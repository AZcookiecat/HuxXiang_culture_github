from flask import Blueprint, request, jsonify, current_app
from flask_jwt_extended import jwt_required
from models.cultural_resource import CulturalResource
from sqlalchemy import text
import math


cultural_resources_bp = Blueprint('cultural_resources', __name__, url_prefix='/api/resources')


@cultural_resources_bp.route('/', methods=['GET'])
def get_resources():
    """获取文化资源列表"""
    try:
        page = request.args.get('page', 1, type=int)
        limit = request.args.get('limit', 10, type=int)
        category = request.args.get('category')
        search = request.args.get('search')
        
        offset = (page - 1) * limit
        
        # 构建查询
        query = current_app.db.session.query(CulturalResource)
        
        if category:
            query = query.filter(CulturalResource.category == category)
            
        if search:
            query = query.filter(
                CulturalResource.title.contains(search) |
                CulturalResource.description.contains(search)
            )
        
        total = query.count()
        resources = query.offset(offset).limit(limit).all()
        
        result = [{
            'id': r.id,
            'title': r.title,
            'description': r.description,
            'type': r.type,
            'category': r.category,
            'tags': r.tags.split(',') if r.tags else [],
            'author': r.author,
            'cover_image': r.cover_image,
            'view_count': r.view_count,
            'like_count': r.like_count,
            'created_at': r.created_at.isoformat(),
        } for r in resources]
        
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
        return jsonify({'message': '获取文化资源列表失败: ' + str(e)}), 500


@cultural_resources_bp.route('/<int:id>', methods=['GET'])
def get_resource(id):
    """获取单个文化资源"""
    try:
        resource = current_app.db.session.get(CulturalResource, id)
        
        if not resource:
            return jsonify({'message': '文化资源不存在'}), 404
            
        # 增加浏览量（排除作者自己）
        # 这里暂时不考虑身份验证，后续可以根据需要添加
        resource.view_count += 1
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'data': {
                'id': resource.id,
                'title': resource.title,
                'description': resource.description,
                'content': resource.content,
                'type': resource.type,
                'category': resource.category,
                'tags': resource.tags.split(',') if resource.tags else [],
                'author': resource.author,
                'source': resource.source,
                'cover_image': resource.cover_image,
                'media_url': resource.media_url,
                'view_count': resource.view_count,
                'like_count': resource.like_count,
                'created_at': resource.created_at.isoformat(),
                'updated_at': resource.updated_at.isoformat()
            }
        })
    except Exception as e:
        return jsonify({'message': '获取文化资源失败: ' + str(e)}), 500


@cultural_resources_bp.route('/<int:id>/like', methods=['POST'])
@jwt_required()
def like_resource(id):
    """点赞文化资源"""
    try:
        resource = current_app.db.session.get(CulturalResource, id)
        
        if not resource:
            return jsonify({'message': '文化资源不存在'}), 404
            
        resource.like_count += 1
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '点赞成功',
            'like_count': resource.like_count
        })
    except Exception as e:
        return jsonify({'message': '点赞失败: ' + str(e)}), 500


@cultural_resources_bp.route('/', methods=['POST'])
@jwt_required()
def create_resource():
    """创建文化资源"""
    try:
        data = request.get_json()
        
        # 验证必需字段
        required_fields = ['title', 'content', 'type', 'category']
        for field in required_fields:
            if not data.get(field):
                return jsonify({'message': f'{field} 是必需的'}), 400
        
        resource = CulturalResource(
            title=data['title'],
            description=data.get('description'),
            content=data['content'],
            type=data['type'],
            category=data['category'],
            tags=','.join(data.get('tags', [])),
            author=data.get('author'),
            source=data.get('source'),
            cover_image=data.get('cover_image'),
            media_url=data.get('media_url')
        )
        
        current_app.db.session.add(resource)
        current_app.db.session.commit()
        
        return jsonify({
            'success': True,
            'message': '文化资源创建成功',
            'data': {
                'id': resource.id,
                'title': resource.title
            }
        }), 201
    except Exception as e:
        current_app.db.session.rollback()
        return jsonify({'message': '创建文化资源失败: ' + str(e)}), 500