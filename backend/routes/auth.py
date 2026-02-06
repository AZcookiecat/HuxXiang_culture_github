from flask import Blueprint, request, jsonify
from flask_jwt_extended import create_access_token, jwt_required, get_jwt_identity
from models.user import User
from app import db
import re


auth_bp = Blueprint('auth', __name__, url_prefix='/api/auth')


@auth_bp.route('/register', methods=['POST'])
def register():
    """用户注册"""
    data = request.get_json()
    
    # 验证输入
    if not data or not data.get('username') or len(data['username']) < 3:
        return jsonify({'message': '用户名至少需要3个字符'}), 400
        
    if not data or not data.get('email') or not re.match(r'^[^@]+@[^@]+\.[^@]+$', data['email']):
        return jsonify({'message': '邮箱格式不正确'}), 400
        
    if not data or not data.get('password') or len(data['password']) < 6:
        return jsonify({'message': '密码至少需要6个字符'}), 400
    
    # 检查用户名或邮箱是否已存在
    if User.query.filter_by(username=data['username']).first():
        return jsonify({'message': '用户名已被占用'}), 400
        
    if User.query.filter_by(email=data['email']).first():
        return jsonify({'message': '邮箱已被注册'}), 400
    
    # 创建新用户
    user = User(
        username=data['username'],
        email=data['email']
    )
    user.set_password(data['password'])
    
    db.session.add(user)
    db.session.commit()
    
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


@auth_bp.route('/login', methods=['POST'])
def login():
    """用户登录，支持用户名或邮箱登录"""
    data = request.get_json()
    
    if not data:
        return jsonify({'message': '请求数据不能为空'}), 400
    
    # 获取用户名/邮箱和密码
    username_or_email = data.get('username') or data.get('usernameOrEmail')
    password = data.get('password')
    
    if not username_or_email or not password:
        return jsonify({'message': '用户名/邮箱和密码不能为空'}), 400
    
    # 尝试通过用户名或邮箱查找用户
    user = User.query.filter((User.username == username_or_email) | (User.email == username_or_email)).first()
    
    if user and user.check_password(password):
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


@auth_bp.route('/profile', methods=['GET'])
@jwt_required()
def profile():
    """获取用户个人信息"""
    user_id = get_jwt_identity()
    user = User.query.get(user_id)
    
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


@auth_bp.route('/profile', methods=['PUT'])
@jwt_required()
def update_profile():
    """更新用户个人信息"""
    user_id = get_jwt_identity()
    user = User.query.get(user_id)
    
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
        existing_user = User.query.filter(User.username == data['username'], User.id != user_id).first()
        if existing_user:
            return jsonify({'message': '用户名已被占用'}), 400
        user.username = data['username']
    
    db.session.commit()
    
    return jsonify({'success': True, 'message': '资料更新成功'})


@auth_bp.route('/logout', methods=['POST'])
@jwt_required()
def logout():
    """用户登出（预留接口）"""
    # JWT是无状态的，客户端只需删除token即可
    return jsonify({'success': True, 'message': '登出成功'})