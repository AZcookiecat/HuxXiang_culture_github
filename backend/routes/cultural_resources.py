from flask import Blueprint, request, jsonify
from models.cultural_resource import CulturalResource
from app import db
from datetime import datetime


cultural_resources_bp = Blueprint('cultural_resources', __name__, url_prefix='/api/resources')


@cultural_resources_bp.route('/', methods=['GET'])
def get_resources():
    """获取文化资源列表"""
    page = request.args.get('page', 1, type=int)
    per_page = request.args.get('per_page', 10, type=int)
    resource_type = request.args.get('type')
    category = request.args.get('category')
    search = request.args.get('search')

    query = CulturalResource.query

    if resource_type:
        query = query.filter(CulturalResource.type == resource_type)
    
    if category:
        query = query.filter(CulturalResource.category == category)
        
    if search:
        query = query.filter(CulturalResource.title.contains(search) | 
                             CulturalResource.description.contains(search))

    resources = query.paginate(page=page, per_page=per_page, error_out=False)

    return jsonify({
        'resources': [
            {
                'id': r.id,
                'title': r.title,
                'description': r.description,
                'type': r.type,
                'category': r.category,
                'cover_image': r.cover_image,
                'view_count': r.view_count,
                'created_at': r.created_at.isoformat(),
                'updated_at': r.updated_at.isoformat()
            } for r in resources.items
        ],
        'pagination': {
            'page': page,
            'pages': resources.pages,
            'per_page': per_page,
            'total': resources.total
        }
    })


@cultural_resources_bp.route('/<int:id>', methods=['GET'])
def get_resource(id):
    """获取特定文化资源详情"""
    resource = CulturalResource.query.get_or_404(id)
    
    # 增加浏览次数
    resource.view_count += 1
    db.session.commit()
    
    return jsonify({
        'id': resource.id,
        'title': resource.title,
        'content': resource.content,
        'description': resource.description,
        'type': resource.type,
        'category': resource.category,
        'tags': resource.tags.split(',') if resource.tags else [],
        'author': resource.author,
        'source': resource.source,
        'cover_image': resource.cover_image,
        'media_url': resource.media_url,
        'priority': resource.priority,
        'view_count': resource.view_count,
        'like_count': resource.like_count,
        'created_at': resource.created_at.isoformat(),
        'updated_at': resource.updated_at.isoformat()
    })


@cultural_resources_bp.route('/', methods=['POST'])
def create_resource():
    """创建新的文化资源（需要管理员权限）"""
    data = request.get_json()
    
    resource = CulturalResource(
        title=data.get('title'),
        description=data.get('description'),
        content=data.get('content'),
        type=data.get('type'),
        category=data.get('category'),
        tags=','.join(data.get('tags', [])),
        author=data.get('author'),
        source=data.get('source'),
        cover_image=data.get('cover_image'),
        media_url=data.get('media_url'),
        priority=data.get('priority', 0),
        status=data.get('status', 'published')
    )
    
    db.session.add(resource)
    db.session.commit()
    
    return jsonify({
        'message': '文化资源创建成功',
        'id': resource.id
    }), 201