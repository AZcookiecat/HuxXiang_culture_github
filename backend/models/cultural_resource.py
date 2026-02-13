from app import db
from datetime import datetime


class CulturalResource(db.Model):
    __tablename__ = 'cultural_resources'

    id = db.Column(db.Integer, primary_key=True)
    title = db.Column(db.String(200), nullable=False)  # 资源标题
    description = db.Column(db.Text)  # 描述
    content = db.Column(db.Text)  # 详细内容
    type = db.Column(db.String(50), nullable=False)  # 类型：历史遗迹、传统艺术、民间文学等
    category = db.Column(db.String(50), nullable=False)  # 分类
    tags = db.Column(db.String(200))  # 标签
    author = db.Column(db.String(100))  # 作者
    source = db.Column(db.String(200))  # 来源
    cover_image = db.Column(db.String(255))  # 封面图片
    media_url = db.Column(db.String(255))  # 媒体文件链接（视频、音频、3D模型等）
    status = db.Column(db.String(20), default='published')  # 状态：draft, published, archived
    priority = db.Column(db.Integer, default=0)  # 排序优先级
    view_count = db.Column(db.Integer, default=0)  # 浏览次数
    like_count = db.Column(db.Integer, default=0)  # 点赞数
    created_at = db.Column(db.DateTime, default=datetime.utcnow)
    updated_at = db.Column(db.DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)

    def __repr__(self):
        return f'<CulturalResource {self.title}>'