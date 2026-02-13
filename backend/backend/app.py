from flask import Flask, request, jsonify
from flask_cors import CORS
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate
import os
import re
import bcrypt  
from config import config

app = Flask(__name__)

# 1. 环境变量与配置加载（增加容错）
env = os.getenv('FLASK_ENV', 'development')  # 简化写法
if env not in config:  # 防止配置不存在报错
    raise ValueError(f"FLASK_ENV={env} 无对应的配置,请检查config.py")
app.config.from_object(config[env])

# 2. 跨域配置（更精细，避免全量跨域）
CORS(app, resources={r"/api/*": {"origins": app.config.get('CORS_ORIGINS', '*')}})

# 3. 数据库初始化（规范写法）
db = SQLAlchemy(app)
migrate = Migrate(app, db)

class SerializerMixin:
    """模型序列化混入类"""
    def to_dict(self, exclude=None):
        if exclude is None:
            exclude = []
        # 默认排除密码字段
        # 替换原来的 exclude.append('password')
        if 'password' not in exclude:
            exclude.append('password')
        #下面是字典推导式,用于将模型实例的属性转换为字典
        result = {}
        for c in self.__table__.columns:
            if c.name not in exclude:
                result[c.name] = getattr(self,c.name)
        return result

# 4. 辅助函数
def hash_password(password: str) -> str: #这儿的 -> str是函数返回值类型注解,这个注解不会强制约束函数的返回值类型,只是为了代码可读性
    """
    安全的密码加密(替换hashlib,使用bcrypt不可逆加密+盐值)
    :param password: 明文密码
    :return: 加密后的密码字符串
    """
    salt = bcrypt.gensalt()  # 随机盐值，每个用户不同
    hashed = bcrypt.hashpw(password.encode('utf-8'), salt)
    return hashed.decode('utf-8')

def verify_password(plain_password: str, hashed_password: str) -> bool:
    """
    验证密码（后续登录会用到）
    :param plain_password: 明文密码
    :param hashed_password: 数据库中的加密密码
    :return: 是否匹配
    """
    return bcrypt.checkpw(plain_password.encode('utf-8'), hashed_password.encode('utf-8'))

def is_valid_email(email: str) -> bool:
    """校验邮箱格式(符合RFC 5322标准)"""
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None

def is_valid_username(username: str) -> bool:
    """校验用户名格式(字母/数字/下划线,2-20位)"""
    pattern = r'^[a-zA-Z0-9_]{2,20}$'
    return re.match(pattern, username) is not None

# 5. 统一错误处理
@app.errorhandler(400)
def bad_request(e):
    return jsonify({'code': 400, 'message': str(e.description) if hasattr(e, 'description') else '请求参数错误'}), 400

@app.errorhandler(500)
def server_error(e):
    return jsonify({'code': 500, 'message': '服务器内部错误'}), 500

# 6. 数据库模型
class User(db.Model, SerializerMixin):
    __tablename__ = 'users'  # 显式指定表名，避免ORM自动生成不规则名称
    id = db.Column(db.Integer, primary_key=True)
    username = db.Column(db.String(50), unique=True, nullable=False, comment='用户名')
    email = db.Column(db.String(100), unique=True, nullable=False, comment='邮箱')
    password = db.Column(db.String(255), nullable=False, comment='加密密码')
    created_at = db.Column(db.DateTime, default=db.func.current_timestamp(), comment='创建时间')
    
    def __repr__(self):
        return f'<User {self.username}>'

class Product(db.Model, SerializerMixin):
    __tablename__ = 'products'
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(100), nullable=False, comment='商品名称')
    description = db.Column(db.Text, nullable=False, comment='商品描述')
    price = db.Column(db.Float, nullable=False, comment='商品价格')
    image = db.Column(db.String(255), nullable=True, comment='商品图片URL')
    category = db.Column(db.String(50), nullable=False, comment='商品分类')
    created_at = db.Column(db.DateTime, default=db.func.current_timestamp(), comment='创建时间')
    
    def __repr__(self):
        return f'<Product {self.name}>'

@app.route('/api/register', methods=['POST'])
def register():
    # 1. 获取 JSON 数据
    data = request.get_json()
    if not data or not isinstance(data, dict):
        return jsonify({'code': 400, 'message': '无效的 JSON 数据'}), 400

    # 2. 提取并清洗参数
    username = data.get('username', '').strip()
    email = data.get('email', '').strip()
    password = data.get('password', '').strip()

    # 3. 参数验证
    validation_errors = []
    if not username: validation_errors.append('用户名不能为空')
    elif not is_valid_username(username): validation_errors.append('用户名格式不正确')
    
    if not email: validation_errors.append('邮箱不能为空')
    elif not is_valid_email(email): validation_errors.append('邮箱格式不正确')
    
    if not password: validation_errors.append('密码不能为空')
    elif len(password) < 6: validation_errors.append('密码长度不能少于 6 位')

    if validation_errors:
        return jsonify({'code': 400, 'message': validation_errors[0]}), 400

    # 4. 业务逻辑处理
    try:
        # 检查唯一性
        if User.query.filter(User.username == username).first():
            return jsonify({'code': 400, 'message': '用户名已存在'}), 400
        if User.query.filter(User.email == email).first():
            return jsonify({'code': 400, 'message': '邮箱已被注册'}), 400

        # 创建并保存
        new_user = User(
            username=username,
            email=email,
            password=hash_password(password)
        )
        db.session.add(new_user)
        db.session.commit()
        
        return jsonify({
            'code': 201,
            'message': '注册成功',
            'data': new_user.to_dict()
        }), 201

    except Exception as e:
        db.session.rollback()
        app.logger.error(f"注册失败: {str(e)}")
        return jsonify({'code': 500, 'message': '服务器繁忙，请稍后再试'}), 500

if __name__ == '__main__':
    # 生产环境建议用Gunicorn，此处仅为开发调试
    app.run(
        host=app.config.get('HOST', '0.0.0.0'),
        port=app.config.get('PORT', 5000),
        debug=app.config.get('DEBUG', False)  # 生产环境关闭debug
    )