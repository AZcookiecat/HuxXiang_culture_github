from flask import Blueprint, request, jsonify, current_app
from flask_jwt_extended import create_access_token, jwt_required, get_jwt_identity
from models.user import User
from sqlalchemy import text
import re


auth_bp = Blueprint('auth', __name__, url_prefix='/api/auth')


@auth_bp.route('/register', methods=['POST'])
def register():
    """用户注册"""
    try:
        data = request.get_json()
        
        # 验证输入
        if not data or not data.get('username') or len(data['username']) < 3:
            return jsonify({'message': '用户名至少需要3个字符'}), 400
            
        if not data or not data.get('email') or not re.match(r'^[^@]+@[^@]+\.[^@]+$', data['email']):
            return jsonify({'message': '邮箱格式不正确'}), 400
            
        if not data or not data.get('password') or len(data['password']) < 6:
            return jsonify({'message': '密码至少需要6个字符'}), 400
        
        # 检查用户名或邮箱是否已存在
        existing_user = current_app.db.session.query(User).filter(
            (User.username == data['username']) | (User.email == data['email'])
        ).first()
        
        if existing_user:
            return jsonify({'message': '用户名或邮箱已被注册'}), 400
        
        # 创建新用户
        user = User(
            username=data['username'],
            email=data['email']
        )
        user.set_password(data['password'])
        
        current_app.db.session.add(user)
        current_app.db.session.flush()  # 获取分配的ID，但暂不提交事务
        
        # 设置默认头像，使用用户名首字母和用户ID生成
        initial = data['username'][0].upper()
        user.avatar = f"https://picsum.photos/seed/{initial}{user.id}/100"
        
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '注册成功',
            'user': {
                'id': user.id,
                'username': user.username,
                'email': user.email,
                'role': user.role,
                'avatar': user.avatar
            }
        }), 201
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '注册失败: ' + str(e)}), 500


@auth_bp.route('/login', methods=['POST'])
def login():
    """用户登录，支持用户名或邮箱登录"""
    try:
        data = request.get_json()
        
        if not data:
            return jsonify({'message': '请求数据不能为空'}), 400
        
        # 获取用户名/邮箱和密码
        username_or_email = data.get('username') or data.get('usernameOrEmail')
        password = data.get('password')
        
        if not username_or_email or not password:
            return jsonify({'message': '用户名/邮箱和密码不能为空'}), 400
        
        # 尝试通过用户名或邮箱查找用户
        user = current_app.db.session.query(User).filter(
            (User.username == username_or_email) | (User.email == username_or_email)
        ).first()
        
        if user and user.check_password(password) and user.is_active:
            access_token = create_access_token(identity=user.id)
            return jsonify({
                'success': True,
                'message': '登录成功',
                'access_token': access_token,
                'user': {
                    'id': user.id,
                    'username': user.username,
                    'email': user.email,
                    'role': user.role,
                    'avatar': user.avatar
                }
            })
        
        return jsonify({'message': '用户名/邮箱或密码错误'}), 401
    except Exception as e:
        return jsonify({'message': '登录失败: ' + str(e)}), 500


@auth_bp.route('/profile', methods=['GET'])
@jwt_required()
def profile():
    """获取用户个人信息"""
    try:
        user_id = get_jwt_identity()
        user = current_app.db.session.get(User, user_id)
        
        if not user:
            return jsonify({'message': '用户不存在'}), 404
        
        return jsonify({
            'success': True,
            'data': {
                'id': user.id,
                'username': user.username,
                'email': user.email,
                'bio': user.bio,
                'avatar': user.avatar,
                'role': user.role,
                'created_at': user.created_at.isoformat()
            }
        })
    except Exception as e:
        return jsonify({'message': '获取用户信息失败: ' + str(e)}), 500


@auth_bp.route('/profile', methods=['PUT'])
@jwt_required()
def update_profile():
    """更新用户个人信息"""
    try:
        user_id = get_jwt_identity()
        user = current_app.db.session.get(User, user_id)
        
        if not user:
            return jsonify({'message': '用户不存在'}), 404
        
        data = request.get_json()
        
        if not data:
            return jsonify({'message': '请求数据不能为空'}), 400

        if 'bio' in data:
            user.bio = data['bio']
        if 'avatar' in data:
            user.avatar = data['avatar']
        if 'username' in data:
            # 检查用户名是否已被其他用户使用
            existing_user = current_app.db.session.query(User).filter(
                User.username == data['username'], 
                User.id != user_id
            ).first()
            if existing_user:
                return jsonify({'message': '用户名已被占用'}), 400
            user.username = data['username']
        
        current_app.db.session.commit()
        
        return jsonify({'success': True, 'message': '资料更新成功'})
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '更新用户资料失败: ' + str(e)}), 500


@auth_bp.route('/upload-avatar', methods=['POST'])
@jwt_required()
def upload_avatar():
    """上传用户头像"""
    try:
        user_id = get_jwt_identity()
        user = current_app.db.session.get(User, user_id)
        
        if not user:
            return jsonify({'message': '用户不存在'}), 404
        
        if 'avatar' not in request.files:
            return jsonify({'message': '没有上传头像文件'}), 400
        
        file = request.files['avatar']
        
        if file.filename == '':
            return jsonify({'message': '未选择文件'}), 400
        
        # 检查文件类型
        allowed_extensions = {'png', 'jpg', 'jpeg', 'gif', 'webp'}
        if not ('.' in file.filename and \
                file.filename.rsplit('.', 1)[1].lower() in allowed_extensions):
            return jsonify({'message': '不支持的文件格式'}), 400
        
        # 生成唯一的文件名
        import uuid
        import os
        extension = file.filename.rsplit('.', 1)[1].lower()
        unique_filename = f"{user.id}_{uuid.uuid4().hex}.{extension}"
        
        # 确定上传目录 - 首先尝试使用环境变量定义的上传目录
        import tempfile
        upload_base = os.environ.get('AVATAR_UPLOAD_PATH', None)
        
        if upload_base:
            # 如果设置了环境变量，则使用环境变量指定的路径
            upload_folder = os.path.join(upload_base, 'avatars')
        else:
            # 否则使用相对于项目根目录的路径
            upload_folder = os.path.join(current_app.root_path, '..', 'public', 'static', 'avatars')
        
        # 使用abspath规范化路径
        upload_folder = os.path.abspath(upload_folder)
        
        os.makedirs(upload_folder, exist_ok=True)
        
        # 保存文件
        file_path = os.path.join(upload_folder, unique_filename)
        file.save(file_path)
        
        # 更新用户头像路径
        # 使用正斜杠构建Web访问路径
        user.avatar = f"/static/avatars/{unique_filename}"
        current_app.db.session.commit()
        
        return jsonify({
            'success': True, 
            'message': '头像上传成功',
            'avatar_url': user.avatar
        })
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '上传头像失败: ' + str(e)}), 500


@auth_bp.route('/logout', methods=['POST'])
@jwt_required()
def logout():
    """用户登出（预留接口）"""
    # JWT是无状态的，客户端只需删除token即可
    return jsonify({'success': True, 'message': '登出成功'})