from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
from flask_jwt_extended import JWTManager
from config import Config


db = SQLAlchemy()
jwt = JWTManager()


def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)
    
    # 初始化扩展
    db.init_app(app)
    CORS(app)  # 允许跨域请求
    jwt.init_app(app)  # 初始化JWT
    
    # 注册蓝图
    from routes.main import main_bp
    from routes.cultural_resources import cultural_resources_bp
    from routes.community import community_bp
    from routes.auth import auth_bp
    
    app.register_blueprint(main_bp)
    app.register_blueprint(cultural_resources_bp)
    app.register_blueprint(community_bp)
    app.register_blueprint(auth_bp)
    
    # 注册全局异常处理器
    @app.errorhandler(Exception)
    def handle_exception(e):
        # 记录错误日志
        app.logger.error(f"Unhandled exception: {str(e)}", exc_info=True)
        
        # 返回JSON格式的错误响应
        import traceback
        return {"error": "服务器内部错误", "details": str(e)}, 500
    
    # 注册shell上下文处理器
    @app.shell_context_processor
    def make_shell_context():
        from models.user import User
        from models.cultural_resource import CulturalResource
        from models.community_post import CommunityPost
        return {'db': db, 'User': User, 'CulturalResource': CulturalResource, 'CommunityPost': CommunityPost}
    
    return app