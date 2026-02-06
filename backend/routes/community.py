from flask import Blueprint, request, jsonify
from flask_jwt_extended import jwt_required, get_jwt_identity
from models.user import User
from models.community_post import CommunityPost, Comment
from app import db


community_bp = Blueprint('community', __name__, url_prefix='/api/community')


@community_bp.route('/posts', methods=['GET'])
def get_posts():
    """获取社区帖子列表"""
    page = request.args.get('page', 1, type=int)
    per_page = request.args.get('per_page', 10, type=int)
    category = request.args.get('category')
    search = request.args.get('search')

    query = CommunityPost.query

    if category:
        query = query.filter(CommunityPost.category == category)
        
    if search:
        query = query.filter(CommunityPost.title.contains(search) | 
                             CommunityPost.content.contains(search))

    posts = query.order_by(CommunityPost.created_at.desc()).paginate(
        page=page, per_page=per_page, error_out=False
    )

    return jsonify({
        'posts': [
            {
                'id': p.id,
                'title': p.title,
                'summary': p.content[:100] + '...' if len(p.content) > 100 else p.content,
                'author': {
                    'id': p.author.id,
                    'username': p.author.username,
                    'avatar': p.author.avatar
                },
                'category': p.category,
                'view_count': p.view_count,
                'like_count': p.like_count,
                'comment_count': p.comment_count,
                'created_at': p.created_at.isoformat()
            } for p in posts.items
        ],
        'pagination': {
            'page': page,
            'pages': posts.pages,
            'per_page': per_page,
            'total': posts.total
        }
    })


@community_bp.route('/posts/<int:id>', methods=['GET'])
def get_post(id):
    """获取帖子详情"""
    post = CommunityPost.query.get_or_404(id)
    
    # 增加浏览次数
    post.view_count += 1
    db.session.commit()
    
    return jsonify({
        'id': post.id,
        'title': post.title,
        'content': post.content,
        'author': {
            'id': post.author.id,
            'username': post.author.username,
            'avatar': post.author.avatar,
            'bio': post.author.bio
        },
        'category': post.category,
        'view_count': post.view_count,
        'like_count': post.like_count,
        'comment_count': post.comment_count,
        'created_at': post.created_at.isoformat(),
        'updated_at': post.updated_at.isoformat(),
        'comments': [
            {
                'id': c.id,
                'content': c.content,
                'author': {
                    'id': c.author.id,
                    'username': c.author.username,
                    'avatar': c.author.avatar
                },
                'created_at': c.created_at.isoformat(),
                'replies': [
                    {
                        'id': r.id,
                        'content': r.content,
                        'author': {
                            'id': r.author.id,
                            'username': r.author.username,
                            'avatar': r.author.avatar
                        },
                        'created_at': r.created_at.isoformat()
                    } for r in c.replies
                ]
            } for c in post.comments if c.parent_id is None
        ]
    })


@community_bp.route('/posts', methods=['POST'])
@jwt_required()
def create_post():
    """创建帖子"""
    user_id = get_jwt_identity()
    user = User.query.get(user_id)
    
    if not user:
        return jsonify({'message': '用户不存在'}), 404
    
    data = request.get_json()
    
    post = CommunityPost(
        title=data.get('title'),
        content=data.get('content'),
        author_id=user.id,
        category=data.get('category', 'discussion'),
        status=data.get('status', 'published')
    )
    
    db.session.add(post)
    db.session.commit()
    
    return jsonify({
        'message': '帖子发布成功',
        'post_id': post.id
    }), 201


@community_bp.route('/posts/<int:post_id>/comments', methods=['POST'])
@jwt_required()
def add_comment(post_id):
    """添加评论"""
    user_id = get_jwt_identity()
    user = User.query.get(user_id)
    
    if not user:
        return jsonify({'message': '用户不存在'}), 404
    
    post = CommunityPost.query.get_or_404(post_id)
    
    data = request.get_json()
    
    comment = Comment(
        content=data.get('content'),
        author_id=user.id,
        post_id=post_id
    )
    
    if data.get('reply_to'):  # 如果是回复某条评论
        parent_comment = Comment.query.get(data.get('reply_to'))
        if parent_comment and parent_comment.post_id == post_id:
            comment.parent = parent_comment
    
    db.session.add(comment)
    db.session.commit()
    
    # 更新评论数量
    post.comment_count += 1
    db.session.commit()
    
    return jsonify({
        'message': '评论成功',
        'comment_id': comment.id
    }), 201