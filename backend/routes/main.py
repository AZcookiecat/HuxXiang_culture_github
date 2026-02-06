from flask import Blueprint, jsonify
from app import db


main_bp = Blueprint('main', __name__)


@main_bp.route('/')
def index():
    return jsonify({
        'message': '湖湘文化数字化平台 API 接口',
        'version': '1.0.0',
        'endpoints': {
            'auth': '/api/auth',
            'resources': '/api/resources',
            'community': '/api/community'
        }
    })


@main_bp.route('/health')
def health_check():
    """健康检查接口"""
    try:
        # 尝试连接数据库
        db.session.execute('SELECT 1')
        return jsonify({'status': 'healthy', 'database': 'connected'}), 200
    except Exception as e:
        return jsonify({'status': 'unhealthy', 'database': 'error', 'error': str(e)}), 500