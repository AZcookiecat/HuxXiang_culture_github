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
        sort_by = request.args.get('sortBy', 'latest')  # 新增排序参数，默认为最新发布
        
        offset = (page - 1) * limit
        
        # 构建查询
        query = current_app.db.session.query(CommunityPost)
        
        if category:
            query = query.filter(CommunityPost.category == category)
        
        query = query.filter(CommunityPost.status == 'published')
        
        # 根据排序参数构建排序条件
        if sort_by == 'popular':  # 热门：按点赞数+评论数
            query = query.order_by((CommunityPost.like_count + CommunityPost.comment_count).desc())
        elif sort_by == 'comments':  # 评论最多：按评论数
            query = query.order_by(CommunityPost.comment_count.desc())
        else:  # 默认按最新发布：按创建时间
            query = query.order_by(CommunityPost.created_at.desc())
        
        total = query.count()
        posts = query.offset(offset).limit(limit).all()
        
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
        
        # 获取当前用户是否已点赞
        liked_by_current_user = False
        try:
            current_user_id = get_jwt_identity()
            if current_user_id:
                from models.community_post import user_likes_table  # 修复导入路径
                existing_like = current_app.db.session.query(user_likes_table).filter(
                    user_likes_table.c.user_id == current_user_id,
                    user_likes_table.c.post_id == id
                ).first()
                liked_by_current_user = bool(existing_like)
        except:
            # 如果用户未登录，get_jwt_identity()会抛出异常，这里忽略即可
            pass
        
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
                'liked_by_current_user': liked_by_current_user,  # 添加当前用户是否点赞的状态
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


@community_bp.route('/posts/<int:id>', methods=['PUT'])
@jwt_required()
def update_post(id):
    """更新帖子"""
    try:
        current_user_id = get_jwt_identity()
        post = current_app.db.session.get(CommunityPost, id)
        
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
        
        # 检查权限：只有作者或管理员可以编辑
        user = current_app.db.session.get(User, current_user_id)
        if post.author_id != current_user_id and user.role != 'admin':
            return jsonify({'message': '没有权限编辑此帖子'}), 403
        
        data = request.get_json()
        if 'title' in data:
            post.title = data['title']
        if 'content' in data:
            post.content = data['content']
        if 'category' in data:
            post.category = data['category']
        
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '帖子更新成功',
            'data': {
                'id': post.id,
                'title': post.title
            }
        })
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '更新帖子失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:id>', methods=['DELETE'])
@jwt_required()
def delete_post(id):
    """删除帖子"""
    try:
        current_user_id = get_jwt_identity()
        post = current_app.db.session.get(CommunityPost, id)
        
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
        
        # 检查权限：只有作者或管理员可以删除
        user = current_app.db.session.get(User, current_user_id)
        if post.author_id != current_user_id and user.role != 'admin':
            return jsonify({'message': '没有权限删除此帖子'}), 403
        
        # 删除相关评论
        from sqlalchemy import delete
        stmt = delete(Comment).where(Comment.post_id == id)
        current_app.db.session.execute(stmt)
        
        # 删除帖子
        current_app.db.session.delete(post)
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '帖子删除成功'
        })
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '删除帖子失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:id>/like', methods=['POST'])
@jwt_required()
def like_post(id):
    """点赞或取消点赞帖子"""
    try:
        current_user_id = get_jwt_identity()
        post = current_app.db.session.get(CommunityPost, id)
        
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
            
        # 检查用户是否已经点赞
        from models.community_post import user_likes_table  # 修复导入路径
        existing_like = current_app.db.session.query(user_likes_table).filter(
            user_likes_table.c.user_id == current_user_id,
            user_likes_table.c.post_id == id
        ).first()
        
        if existing_like:
            # 如果已点赞，则取消点赞
            current_app.db.session.execute(
                user_likes_table.delete().where(
                    user_likes_table.c.user_id == current_user_id
                ).where(
                    user_likes_table.c.post_id == id
                )
            )
            post.like_count = max(0, post.like_count - 1)  # 确保不会小于0
            message = '取消点赞成功'
            liked = False
        else:
            # 如果未点赞，则添加点赞
            current_app.db.session.execute(
                user_likes_table.insert().values(
                    user_id=current_user_id,
                    post_id=id
                )
            )
            post.like_count += 1
            message = '点赞成功'
            liked = True
        
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': message,
            'like_count': post.like_count,
            'liked': liked
        })
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '操作失败: ' + str(e)}), 500


@community_bp.route('/posts/<int:post_id>/comments', methods=['GET'])
def get_comments(post_id):
    """获取帖子的评论列表"""
    try:
        # 获取帖子
        post = current_app.db.session.get(CommunityPost, post_id)
        if not post:
            return jsonify({'message': '帖子不存在'}), 404
        
        # 查询帖子的评论，并按时间倒序排列
        comments = current_app.db.session.query(Comment).filter(
            Comment.post_id == post_id,
            Comment.parent_id.is_(None)  # 只获取一级评论，不包含回复
        ).order_by(Comment.created_at.desc()).all()
        
        # 获取评论作者信息
        from models.user import User
        comments_data = []
        for comment in comments:
            author = current_app.db.session.get(User, comment.author_id)
            # 确保即使作者不存在也有默认值
            if author and author.avatar:
                # 使用用户上传的头像
                avatar_url = author.avatar
            elif author:
                # 如果用户存在但没有头像，使用基于用户名首字母+ID的头像
                initial = author.username[0].upper() if author.username else 'U'
                avatar_url = f'https://picsum.photos/seed/{initial}{author.id}/100'
            else:
                # 如果作者不存在，使用默认头像
                avatar_url = f'https://picsum.photos/seed/default{comment.author_id}/100'
                
            author_data = {
                'id': author.id if author else None,
                'username': author.username if author else '匿名用户',
                'avatar': avatar_url
            }
            comments_data.append({
                'id': comment.id,
                'content': comment.content,
                'created_at': comment.created_at.isoformat(),
                'updated_at': comment.updated_at.isoformat(),
                'author': author_data
            })
        
        return jsonify({
            'success': True,
            'data': comments_data,
            'count': len(comments_data)
        })
    except Exception as e:
        return jsonify({'message': '获取评论失败: ' + str(e)}), 500


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
        
        # 获取评论作者信息用于返回
        from models.user import User
        author = current_app.db.session.get(User, current_user_id)
        # 根据作者信息生成头像URL
        if author and author.avatar:
            # 使用用户上传的头像
            avatar_url = author.avatar
        elif author:
            # 如果用户存在但没有头像，使用基于用户名首字母+ID的头像
            initial = author.username[0].upper() if author.username else 'U'
            avatar_url = f'https://picsum.photos/seed/{initial}{author.id}/100'
        else:
            # 如果作者不存在，使用默认头像
            avatar_url = f'https://picsum.photos/seed/default{current_user_id}/100'
        
        author_data = {
            'id': author.id if author else None,
            'username': author.username if author else '匿名用户',
            'avatar': avatar_url
        }
        
        return jsonify({
            'success': True,
            'message': '评论发布成功',
            'data': {
                'id': comment.id,
                'content': comment.content,
                'author': author_data
            }
        }), 201
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '发布评论失败: ' + str(e)}), 500


@community_bp.route('/posts/related/<int:post_id>', methods=['GET'])
def get_related_posts(post_id):
    """获取相关帖子"""
    try:
        # 获取当前帖子
        current_post = current_app.db.session.get(CommunityPost, post_id)
        if not current_post:
            return jsonify({'message': '帖子不存在'}), 404
        
        # 获取参数
        category = request.args.get('category', '')
        limit = request.args.get('limit', 2, type=int)
        
        # 查询相关帖子：同类别的其他帖子
        query = current_app.db.session.query(CommunityPost).filter(
            CommunityPost.id != post_id,  # 排除当前帖子
            CommunityPost.status == 'published',  # 只查询已发布的帖子
            CommunityPost.category == current_post.category  # 同一类别
        )
        
        # 按点赞数和评论数排序
        related_posts = query.order_by(
            (CommunityPost.like_count + CommunityPost.comment_count).desc()
        ).limit(limit).all()
        
        result = []
        for post in related_posts:
            # 获取作者信息
            author = current_app.db.session.get(User, post.author_id)
            result.append({
                'id': post.id,
                'title': post.title,
                'content': post.content[:150] + '...' if len(post.content) > 150 else post.content,
                'author': {
                    'id': author.id if author else None,
                    'username': author.username if author else '匿名用户',
                    'avatar': author.avatar if author and hasattr(author, 'avatar') else None
                },
                'category': post.category,
                'view_count': post.view_count,
                'like_count': post.like_count,
                'comment_count': post.comment_count,
                'created_at': post.created_at.isoformat(),
            })
        
        return jsonify({
            'success': True,
            'data': result
        })
    except Exception as e:
        return jsonify({'message': '获取相关帖子失败: ' + str(e)}), 500

# 添加删除评论的路由
@community_bp.route('/comments/<int:comment_id>', methods=['DELETE'])
@jwt_required()
def delete_comment(comment_id):
    """删除评论"""
    try:
        current_user_id = get_jwt_identity()
        comment = current_app.db.session.get(Comment, comment_id)
        
        if not comment:
            return jsonify({'message': '评论不存在'}), 404
        
        # 检查权限：只有评论作者或管理员可以删除
        user = current_app.db.session.get(User, current_user_id)
        if comment.author_id != current_user_id and user.role != 'admin':
            return jsonify({'message': '没有权限删除此评论'}), 403
        
        # 删除评论及其所有回复
        from sqlalchemy import delete
        # 先删除回复
        stmt_replies = delete(Comment).where(Comment.parent_id == comment_id)
        current_app.db.session.execute(stmt_replies)
        
        # 再删除主评论
        stmt_main = delete(Comment).where(Comment.id == comment_id)
        current_app.db.session.execute(stmt_main)
        
        # 更新帖子的评论计数
        post = current_app.db.session.get(CommunityPost, comment.post_id)
        if post:
            # 重新计算评论数（减去被删除的评论和它的回复）
            remaining_comments = current_app.db.session.query(Comment).filter(
                Comment.post_id == post.id,
                Comment.parent_id.is_(None)
            ).count()
            post.comment_count = remaining_comments
            current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '评论删除成功'
        })
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '删除评论失败: ' + str(e)}), 500
