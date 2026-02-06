from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_cors import CORS
from config import Config


db = SQLAlchemy()


def create_app():
    app = Flask(__name__)
    app.config.from_object(Config)
    
    # 初始化扩展
    db.init_app(app)
    CORS(app)  # 允许跨域请求
    
    # 注册蓝图
    from routes.main import main_bp
    from routes.cultural_resources import cultural_resources_bp
    from routes.community import community_bp
    from routes.auth import auth_bp
    
    app.register_blueprint(main_bp)
    app.register_blueprint(cultural_resources_bp)
    app.register_blueprint(community_bp)
    app.register_blueprint(auth_bp)
    
    return app