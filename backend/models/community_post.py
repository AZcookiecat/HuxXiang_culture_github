from app import db
from datetime import datetime


# 用户点赞关联表
user_likes_table = db.Table('user_post_likes',
    db.Column('user_id', db.Integer, db.ForeignKey('users.id'), primary_key=True),
    db.Column('post_id', db.Integer, db.ForeignKey('community_posts.id'), primary_key=True)
)


class CommunityPost(db.Model):
    __tablename__ = 'community_posts'

    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(200), nullable=False)
    content = db.Column(db.Text, nullable=False)
    author_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    category = db.Column(db.String(50), nullable=False)  # 分类：讨论、提问、分享等
    status = db.Column(db.String(20), default='published')  # 状态：draft, published, hidden
    view_count = db.Column(db.Integer, default=0)  # 浏览次数
    like_count = db.Column(db.Integer, default=0)  # 点赞数
    comment_count = db.Column(db.Integer, default=0)  # 评论数
    created_at = db.Column(db.DateTime, default=datetime.utcnow)
    updated_at = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # 关联评论
    comments = db.relationship('Comment', backref='post', lazy=True, foreign_keys='Comment.post_id')

    # 关联点赞用户
    liked_users = db.relationship('User', secondary=user_likes_table, lazy='subquery',
                                  backref=db.backref('liked_posts', lazy=True))

    def __repr__(self):
        return f'<CommunityPost {self.title}>'


class Comment(db.Model):
    __tablename__ = 'comments'

    id = db.Column(db.Integer, primary_key=True)
    content = db.Column(db.Text, nullable=False)
    author_id = db.Column(db.Integer, db.ForeignKey('users.id'), nullable=False)
    post_id = db.Column(db.Integer, db.ForeignKey('community_posts.id'), nullable=False)
    parent_id = db.Column(db.Integer, db.ForeignKey('comments.id'), nullable=True)  # 回复评论的功能
    created_at = db.Column(db.DateTime, default=datetime.utcnow)
    updated_at = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    # 子评论
    replies = db.relationship('Comment', backref=db.backref('parent', remote_side=[id]), lazy=True, foreign_keys='Comment.parent_id')

    def __repr__(self):
        return f'<Comment {self.id}>'