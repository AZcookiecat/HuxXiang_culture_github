import os
from datetime import timedelta


class Config:
    # 数据库配置 - 根据实际环境调整
    SQLALCHEMY_DATABASE_URI = os.environ.get('DATABASE_URL') or 'mysql+pymysql://root:password@localhost/huxiang_culture'
    SQLALCHEMY_TRACK_MODIFICATIONS = False
    
    # 密钥配置
    SECRET_KEY = os.environ.get('SECRET_KEY') or 'huxiang-culture-secret-key-dev'
    
    # JWT配置
    JWT_SECRET_KEY = os.environ.get('JWT_SECRET_KEY') or 'jwt-huxiang-secret-key-dev'
    JWT_ACCESS_TOKEN_EXPIRES = timedelta(hours=1)
    
    # 文件上传配置
    UPLOAD_FOLDER = 'static/uploads'
    MAX_CONTENT_LENGTH = 16 * 1024 * 1024  # 16MB max file size