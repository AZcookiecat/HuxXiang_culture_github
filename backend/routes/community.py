from flask import Blueprint, request, jsonify, current_app
from flask_jwt_extended import jwt_required, get_jwt_identity
from models.community_post import CommunityPost, Comment
from models.user import User
from sqlalchemy import text
import math


community_bp = Blueprint('community', __name__, url_prefix='/api/community')


@community_bp.route('/posts', methods=['GET'])
def get_posts():
    """获取帖子列表"""
    try:
        page = request.args.get('page', 1, type=int)
        limit = request.args.get('limit', 10, type=int)
        category = request.args.get('category')
        
        offset = (page - 1) * limit
        
        # 构建查询
        query = current_app.db.session.query(CommunityPost)
        
        if category:
            query = query.filter(CommunityPost.category == category)
        
        query = query.filter(CommunityPost.status == 'published')
        
        total = query.count()
        posts = query.order_by(CommunityPost.created_at.desc()).offset(offset).limit(limit).all()
        
        result = []
        for post in posts:
            # 获取作者信息
            author = current_app.db.session.get(User, post.author_id)
            result.append({
                'id': post.id,
                'title': post.title,
                'summary': post.content[:100] + '...' if len(post.content) > 100 else post.content,
                'author': {
                    'id': author.id,
                    'username': author.username,
                    'avatar': author.avatar
                },
                'category': post.category,
                'view_count': post.view_count,
                'like_count': post.like_count,
                'comment_count': post.comment_count,
                'created_at': post.created_at.isoformat(),
            })
        
        return jsonify({
            'success': True,
            'data': result,
            'pagination': {
                'page': page,
                'limit': limit,
                'total': total,
                'pages': math.ceil(total / limit)
            }
        })
    except Exception as e:
        return jsonify({'message': '获取帖子列表失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:id>', methods=['GET'])
def get_post(id):
    """获取单个帖子"""
    try:
        post = current_app.db.session.get(CommunityPost, id)
        
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
            
        if post.status != 'published':
            # 如果不是发布的状态，需要验证用户权限
            current_user_id = get_jwt_identity()
            if not current_user_id or current_user_id != post.author_id:
                return jsonify({'message': '没有权限查看此帖子'}), 403
        
        # 增加浏览量
        post.view_count += 1
        current_app.db.session.commit()
        
        # 获取作者信息
        author = current_app.db.session.get(User, post.author_id)
        
        # 获取评论
        comments = current_app.db.session.query(Comment).filter(
            Comment.post_id == id,
            Comment.parent_id.is_(None)  # 只获取顶级评论
        ).order_by(Comment.created_at.desc()).all()
        
        comments_data = []
        for comment in comments:
            comment_author = current_app.db.session.get(User, comment.author_id)
            replies = current_app.db.session.query(Comment).filter(
                Comment.parent_id == comment.id
            ).all()
            
            reply_data = []
            for reply in replies:
                reply_author = current_app.db.session.get(User, reply.author_id)
                reply_data.append({
                    'id': reply.id,
                    'content': reply.content,
                    'author': {
                        'id': reply_author.id,
                        'username': reply_author.username,
                        'avatar': reply_author.avatar
                    },
                    'created_at': reply.created_at.isoformat()
                })
            
            comments_data.append({
                'id': comment.id,
                'content': comment.content,
                'author': {
                    'id': comment_author.id,
                    'username': comment_author.username,
                    'avatar': comment_author.avatar
                },
                'replies': reply_data,
                'created_at': comment.created_at.isoformat()
            })
        
        return jsonify({
            'success': True,
            'data': {
                'id': post.id,
                'title': post.title,
                'content': post.content,
                'author': {
                    'id': author.id,
                    'username': author.username,
                    'avatar': author.avatar,
                    'bio': author.bio
                },
                'category': post.category,
                'view_count': post.view_count,
                'like_count': post.like_count,
                'comment_count': post.comment_count,
                'created_at': post.created_at.isoformat(),
                'updated_at': post.updated_at.isoformat(),
                'comments': comments_data
            }
        })
    except Exception as e:
        return jsonify({'message': '获取帖子详情失败: ' + str(e)}), 500


@community_bp.route('/posts', methods=['POST'])
@jwt_required()
def create_post():
    """创建帖子"""
    try:
        current_user_id = get_jwt_identity()
        data = request.get_json()
        
        required_fields = ['title', 'content', 'category']
        for field in required_fields:
            if not data.get(field):
                return jsonify({'message': f'{field} 是必需的'}), 400
        
        post = CommunityPost(
            title=data['title'],
            content=data['content'],
            author_id=current_user_id,
            category=data['category']
        )
        
        current_app.db.session.add(post)
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '帖子发布成功',
            'data': {
                'id': post.id,
                'title': post.title
            }
        }), 201
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '发布帖子失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:id>/like', methods=['POST'])
@jwt_required()
def like_post(id):
    """点赞帖子"""
    try:
        post = current_app.db.session.get(CommunityPost, id)
        
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
            
        post.like_count += 1
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '点赞成功',
            'like_count': post.like_count
        })
    except Exception as e:
        return jsonify({'message': '点赞失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:post_id>/comments', methods=['POST'])
@jwt_required()
def add_comment(post_id):
    """添加评论"""
    try:
        current_user_id = get_jwt_identity()
        data = request.get_json()
        
        if not data.get('content'):
            return jsonify({'message': '评论内容不能为空'}), 400
        
        # 检查帖子是否存在
        post = current_app.db.session.get(CommunityPost, post_id)
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
        
        comment = Comment(
            content=data['content'],
            author_id=current_user_id,
            post_id=post_id
        )
        
        # 如果是回复评论
        if data.get('parent_id'):
            parent_comment = current_app.db.session.get(Comment, data['parent_id'])
            if not parent_comment or parent_comment.post_id != post_id:
                return jsonify({'message': '回复的评论不存在'}), 404
            comment.parent_id = parent_comment.id
        
        current_app.db.session.add(comment)
        current_app.db.session.commit()
        
        # 更新评论数
        post.comment_count += 1
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '评论发布成功',
            'data': {
                'id': comment.id,
                'content': comment.content
            }
        }), 201
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '发布评论失败: ' + str(e)}), 500