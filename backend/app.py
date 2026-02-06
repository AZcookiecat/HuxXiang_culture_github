from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from models.user import User
from models.cultural_resource import CulturalResource
from models.community_post import CommunityPost
import os


app = Flask(__name__)

# 从app模块导入db和create_app
from app import db, create_app


app = create_app()


@app.shell_context_processor
def make_shell_context():
    return {'db': db, 'User': User, 'CulturalResource': CulturalResource, 'CommunityPost': CommunityPost}


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)