import os
from app import create_app, db
from models.user import User
from models.cultural_resource import CulturalResource
from models.community_post import CommunityPost, user_likes_table


def init_database():
    """初始化数据库并创建初始数据"""
    app = create_app()
    
    with app.app_context():
        # 创建所有表
        print("正在创建数据库表...")
        db.create_all()
        
        # 创建点赞关联表（如果不存在）
        print("正在创建点赞关联表...")
        user_likes_table.create(db.engine, checkfirst=True)
        print("点赞关联表创建完成")
        
        print("数据库表创建完成")
        
        # 检查是否已有管理员账户
        admin = User.query.filter_by(username='admin').first()
        if not admin:
            print("正在创建管理员账户...")
            admin = User(
                username='admin',
                email='admin@example.com',
                role='admin'
            )
            admin.set_password('admin123')  # 使用默认管理员密码
            
            db.session.add(admin)
            db.session.commit()
            print("管理员账户创建完成")
        else:
            print("管理员账户已存在")
            
        # 检查是否有示例文化资源
        sample_resource = CulturalResource.query.filter_by(title='湖湘文化简介').first()
        if not sample_resource:
            print("正在创建示例文化资源...")
            sample_resource = CulturalResource(
                title='湖湘文化简介',
                description='湖湘文化是湖南地区特有的地域文化，具有深厚的历史底蕴。',
                content='湖湘文化是指湖南地区特有的地域文化，源远流长，博大精深...',
                type='history',
                category='introduction',
                tags='湖湘,湖南,文化,历史',
                author='系统管理员',
                status='published',
                priority=1
            )
            
            db.session.add(sample_resource)
            db.session.commit()
            print("示例文化资源创建完成")
        else:
            print("示例文化资源已存在")
        
        # 检查是否有示例社区帖子
        sample_post = CommunityPost.query.filter_by(title='欢迎来到湖湘文化社区').first()
        if not sample_post:
            print("正在创建示例社区帖子...")
            # 需要先获取一个用户作为作者
            user = User.query.first()
            if not user:
                # 如果没有用户，创建一个普通用户
                user = User(
                    username='demo_user',
                    email='demo@example.com',
                    role='user'
                )
                user.set_password('demo123')
                db.session.add(user)
                
            sample_post = CommunityPost(
                title='欢迎来到湖湘文化社区',
                content='欢迎大家在这里分享关于湖湘文化的见解、故事和资源。让我们一起探讨湖湘文化的魅力！',
                author_id=user.id,
                category='文化讨论',
                status='published'
            )
            
            db.session.add(sample_post)
            db.session.commit()
            print("示例社区帖子创建完成")
        else:
            print("示例社区帖子已存在")
        
        print("数据库初始化完成！")


if __name__ == '__main__':
    init_database()