<template>
  <div class="post-detail-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i> 返回社区
      </button>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <!-- 帖子详情 -->
      <div class="post-detail">
        <!-- 帖子头部 -->
        <div class="post-header">
          <div class="post-author">
            <img :src="currentPost?.authorAvatar" alt="Author avatar" class="author-avatar" loading="lazy" />
            <div class="author-info">
              <span class="author-name">{{ currentPost?.authorName }}</span>
              <span class="post-date">{{ formatDate(currentPost?.createdAt) }}</span>
            </div>
          </div>
          <div class="post-stats">
            <span class="stat-item"><i class="far fa-eye"></i> {{ currentPost?.views }} 浏览</span>
            <span class="stat-item"><i class="far fa-comment"></i> {{ currentPost?.comments }} 评论</span>
            <span class="stat-item"><i class="far fa-thumbs-up"></i> {{ currentPost?.likes }} 点赞</span>
          </div>
        </div>
        
        <!-- 帖子内容 -->
        <div class="post-content">
          <h1 class="post-title">{{ currentPost?.title }}</h1>
          <div class="post-body" v-html="formatPostContent(currentPost?.content)"></div>
          <div class="post-tags">
            <span class="post-category">{{ currentPost?.category }}</span>
          </div>
        </div>
        
        <!-- 帖子操作 -->
        <div class="post-actions">
          <button class="action-btn" @click="toggleLike">
            <i :class="['far', liked ? 'fa-thumbs-up' : 'fa-thumbs-up', liked ? 'liked' : '']"></i>
            <span>{{ liked ? '已点赞' : '点赞' }}</span>
          </button>
          <button class="action-btn">
            <i class="far fa-share-square"></i>
            <span>分享</span>
          </button>
          <button class="action-btn" v-if="isUserAuthor">
            <i class="far fa-edit"></i>
            <span>编辑</span>
          </button>
        </div>
      </div>
      
      <!-- 相关推荐 -->
      <div class="related-posts">
        <h3>相关推荐</h3>
        <div class="related-posts-list">
          <div v-for="post in relatedPosts" :key="post.id" class="related-post-item" @click="goToPostDetail(post.id)">
            <h4 class="related-post-title">{{ post.title }}</h4>
            <div class="related-post-meta">
              <span>{{ post.authorName }}</span>
              <span>{{ formatDate(post.createdAt) }}</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 评论区 -->
      <div class="comments-section">
        <h3>评论区</h3>
        
        <!-- 发表评论 -->
        <div class="comment-form" v-if="isLoggedIn">
          <textarea v-model="newComment" placeholder="写下您的评论..." rows="3"></textarea>
          <button class="submit-comment-btn" @click="submitComment">发表评论</button>
        </div>
        <div v-else class="login-prompt">
          请先 <router-link to="/login">登录</router-link> 后发表评论
        </div>
        
        <!-- 评论列表 -->
        <div class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <img :src="comment.authorAvatar" alt="Comment author" class="comment-avatar" loading="lazy" />
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-author">{{ comment.authorName }}</span>
                <span class="comment-date">{{ formatDate(comment.createdAt) }}</span>
              </div>
              <p class="comment-text">{{ comment.content }}</p>
            </div>
          </div>
        </div>
        
        <!-- 加载更多评论 -->
        <div class="load-more">
          <button class="load-more-btn" @click="loadMoreComments" v-if="hasMoreComments">
            加载更多评论
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

// 模拟帖子数据服务
const mockPostService = {
  getPostById(id) {
    // 这里模拟从API获取帖子详情
    // 在实际应用中，应该替换为真实的API调用
    const posts = [
      {
        id: '1',
        title: '探讨湖湘文化的现代传承与发展',
        content: `
          <p>湖湘文化作为中国传统文化的重要组成部分，在现代社会如何更好地传承和发展是一个值得思考的问题。</p>
          
          <h3>一、湖湘文化的历史渊源</h3>
          <p>湖湘文化是指湖南省境内的地域文化，具有悠久的历史和丰富的内涵。它起源于先秦时期，经过秦汉、唐宋、明清等各个历史时期的发展，形成了独特的文化特色。</p>
          
          <h3>二、湖湘文化的现代价值</h3>
          <p>在现代社会，湖湘文化仍然具有重要的价值。它不仅是湖南人民的精神财富，也是中华民族文化宝库中的重要组成部分。湖湘文化中的"经世致用"、"敢为人先"等精神，对于现代社会的发展仍然具有重要的启示意义。</p>
          
          <h3>三、湖湘文化的传承与创新</h3>
          <p>传承和创新是湖湘文化发展的两大主题。我们需要在传承传统文化精髓的同时，不断进行创新，使湖湘文化适应现代社会的发展需求。</p>
          
          <p>以上是我对湖湘文化现代传承与发展的一些思考，欢迎大家一起探讨。</p>
        `,
        excerpt: '湖湘文化作为中国传统文化的重要组成部分，在现代社会如何更好地传承和发展是一个值得思考的问题...',
        authorName: '文化探索者',
        authorAvatar: 'https://picsum.photos/seed/user1/100/100',
        createdAt: '2023-06-15',
        views: 528,
        comments: 42,
        likes: 89,
        category: '文化讨论'
      },
      {
        id: '2',
        title: '岳麓书院的历史与文化价值',
        content: `
          <p>岳麓书院是中国古代四大书院之一，位于湖南省长沙市岳麓山下，是湖湘文化的重要载体。</p>
          
          <h3>一、岳麓书院的历史沿革</h3>
          <p>岳麓书院始建于北宋开宝九年（公元976年），历经宋、元、明、清各代，至清末光绪二十九年（公元1903年）改为湖南高等学堂，1926年正式定名为湖南大学。千余年来，岳麓书院弦歌不绝、办学不已，故有"千年学府"的美称。</p>
          
          <h3>二、岳麓书院的文化价值</h3>
          <p>岳麓书院不仅是中国古代教育的重要场所，也是中国传统文化的重要载体。它培养了大批杰出人才，如王夫之、魏源、曾国藩、左宗棠、杨昌济等，对中国历史和文化的发展产生了深远的影响。</p>
          
          <h3>三、岳麓书院的建筑特色</h3>
          <p>岳麓书院的建筑风格独特，体现了中国传统文化的精神内涵。它坐落在岳麓山下，依山而建，布局严谨，环境优美，是中国古代书院建筑的典范。</p>
        `,
        excerpt: '岳麓书院是中国古代四大书院之一，位于湖南省长沙市岳麓山下，是湖湘文化的重要载体...',
        authorName: '历史学者',
        authorAvatar: 'https://picsum.photos/seed/user2/100/100',
        createdAt: '2023-06-10',
        views: 389,
        comments: 27,
        likes: 65,
        category: '历史研究'
      },
      {
        id: '3',
        title: '湘绣艺术的魅力与传承',
        content: `
          <p>湘绣是中国四大名绣之一，起源于湖南省长沙、湘潭一带，具有悠久的历史和独特的艺术风格。</p>
          
          <h3>一、湘绣的历史与发展</h3>
          <p>湘绣的历史可以追溯到两千多年前的春秋战国时期。经过汉、唐、宋、元、明、清等各个历史时期的发展，湘绣逐渐形成了自己独特的艺术风格。</p>
          
          <h3>二、湘绣的艺术特色</h3>
          <p>湘绣以其精湛的技艺、丰富的色彩和生动的形象而著称于世。它的主要特点是：针法多样、色彩鲜艳、形象生动、题材广泛。湘绣的代表作品有《百鸟朝凤》、《芙蓉鲤鱼》等。</p>
          
          <h3>三、湘绣的传承与发展</h3>
          <p>在现代社会，湘绣面临着传承和发展的挑战。一方面，我们需要保护和传承传统的湘绣技艺；另一方面，我们也需要创新湘绣的表现形式和题材内容，使湘绣适应现代社会的审美需求。</p>
        `,
        excerpt: '湘绣是中国四大名绣之一，起源于湖南省长沙、湘潭一带，具有悠久的历史和独特的艺术风格...',
        authorName: '艺术爱好者',
        authorAvatar: 'https://picsum.photos/seed/user3/100/100',
        createdAt: '2023-06-05',
        views: 412,
        comments: 33,
        likes: 76,
        category: '传统艺术'
      },
      {
        id: '4',
        title: '桃花源记：探寻陶渊明笔下的湘楚秘境',
        content: `
          <p>《桃花源记》的常德原型，为何能成为湖湘文化里的“精神原乡”？</p>
          
          <h3>一、桃花源的历史溯源与现实印证</h3>
          <p>陶渊明笔下的桃花源，是千百年中国人心中最向往的隐逸秘境。湖南常德的桃花源景区，正是这篇千古名篇的现实溯源地。这里林壑幽深、落英缤纷，复刻着文中“阡陌交通，鸡犬相闻”的悠然景致。自东晋以来，历代文人墨客在此留下无数题咏，让这片土地成为隐逸文化的物质载体。</p>
          
          <h3>二、隐逸文化与湖湘精神的交融</h3>
          <p>湖湘文化向来以“经世致用、敢为人先”的入世精神著称，而桃花源所代表的隐逸文化，却与这种入世精神相融相生。它成为湖湘人在奋进之外，对安宁生活的美好期许。这种“进可经世济民，退可归隐山林”的双重选择，正是湖湘文化包容多元的最佳体现。</p>
          
          <h3>三、当代桃花源的文化价值</h3>
          <p>如今的桃花源，早已不只是一处景区，更是国人心中的精神净土。它所承载的“和谐共生”理念，在今天依然具有深刻的现实意义。我们可以探讨，在快节奏的现代生活中，如何传承这种古老的生活智慧，让“桃花源”的理想照进现实。</p>
        `,
        excerpt: '《桃花源记》的常德原型，为何能成为湖湘文化里的“精神原乡”？',
        authorName: '文化探索者',
        authorAvatar: 'https://picsum.photos/seed/user4/100/100',
        createdAt: '2023-05-30',
        views: 678,
        comments: 54,
        likes: 123,
        category: '文化讨论'
      },
      {
        id: '5',
        title: '湘菜：辣里寻味的湖湘烟火气',
        content: `
          <p>不止于辣！湘菜的百味江湖，藏着怎样的湖湘性格？</p>
          
          <h3>一、湘菜的风味密码：不止于辣</h3>
          <p>作为中国八大菜系之一，湘菜以“辣”出圈，却从不止于辣。鲜、香、酸、辣交融，才是湘菜的精髓。从桌桌必点的辣椒炒肉，到红亮醇厚的剁椒鱼头，从香糯不腻的毛氏红烧肉，到火宫殿里外焦里嫩的臭豆腐，每一道菜都裹着浓浓的湖湘烟火气。</p>
          
          <h3>二、湘菜里的湖湘性格</h3>
          <p>湘菜的辣，是直爽的辣；湘菜的鲜，是食材本真的鲜。这恰如湖湘人民豪爽热情、质朴实在的性格。一方水土养一方人，一方美食铸一方文化。湘菜中对食材本味的尊重，对火候的极致追求，都映射着湖湘人“吃得苦、霸得蛮、耐得烦”的精神特质。</p>
          
          <h3>三、湘菜的传承与创新</h3>
          <p>在全球化的今天，湘菜正以开放的姿态走向世界。我们既要守护传统湘菜的经典味道，也要勇于创新，让更多人爱上湖湘味道。不妨来聊聊，你心中的“湘菜灵魂菜”是哪一道？你认为湘菜该如何在传承中创新？</p>
        `,
        excerpt: '不止于辣！湘菜的百味江湖，藏着怎样的湖湘性格？',
        authorName: '美食爱好者',
        authorAvatar: 'https://picsum.photos/seed/user5/100/100',
        createdAt: '2023-05-25',
        views: 892,
        comments: 76,
        likes: 156,
        category: '饮食文化'
      },
      {
        id: '6',
        title: '衡山：寿岳之山的千年香火与文化传承',
        content: `
          <p>南岳衡山，为何能成为湖湘大地的信仰高地与自然秘境？</p>
          
          <h3>一、衡山的自然与人文之美</h3>
          <p>衡山，五岳中的“寿岳”，坐落于湖南衡阳，既是国家级风景名胜区，也是佛道共存的千年圣地。这里群峰叠翠，祝融峰的日出云海震撼人心，藏经殿的古木参天静谧清幽，水帘洞的飞瀑流泉灵动多姿，自然之美与人文之韵在此完美交融。</p>
          
          <h3>二、佛道共存的文化奇观</h3>
          <p>山上的南岳大庙，佛殿、道宫、儒祠齐聚一堂，千年来香火不断，是中国南方规模最大、保存最完整的庙宇之一。这种“三教共存”的格局，不仅体现了湖湘文化的包容精神，也让衡山成为了独特的文化研究样本。</p>
          
          <h3>三、衡山的当代意义</h3>
          <p>衡山的美，不仅在山水，更在千百年沉淀的文化与信仰。它是湖湘大地的精神地标，也是无数人心中的祈福圣地。我们可以探讨，在当代社会，如何更好地保护和传承这份珍贵的文化遗产，让衡山的香火永续，让寿岳的精神滋养更多人。</p>
        `,
        excerpt: '南岳衡山，为何能成为湖湘大地的信仰高地与自然秘境？',
        authorName: '山水行者',
        authorAvatar: 'https://picsum.photos/seed/user6/100/100',
        createdAt: '2023-05-20',
        views: 567,
        comments: 43,
        likes: 98,
        category: '自然文化'
      },
      {
        id: '7',
        title: '花鼓戏：乡音里的湖湘百态',
        content: `
          <p>一口乡音入戏来！湖南花鼓戏，如何守住湖湘的文化根脉？</p>
          
          <h3>一、花鼓戏的乡土基因</h3>
          <p>湖南花鼓戏，是湖湘大地最具代表性的地方戏曲，也是国家级非物质文化遗产。它以活泼明快的唱腔、生动诙谐的表演、贴近生活的剧情，演绎着寻常百姓的喜怒哀乐，深深扎根于湖湘的乡土之间。《刘海砍樵》的经典对唱家喻户晓，成为刻在湖湘人DNA里的乡音。</p>
          
          <h3>二、传统艺术的传承困境</h3>
          <p>在流行文化的冲击下，这门传统艺术也曾面临传承的挑战。年轻观众的流失、专业人才的断层，都让花鼓戏的生存空间受到挤压。但值得欣慰的是，越来越多的年轻演员加入，为花鼓戏融入现代元素，让老戏焕发出新的生机。</p>
          
          <h3>三、让乡音走向更广阔的舞台</h3>
          <p>我们可以思考，如何借助新媒体的力量，让更多年轻人爱上这门湖湘传统艺术？比如，开发花鼓戏主题的短视频、动漫，或者与流行音乐跨界合作，让古老的唱腔与现代审美产生共鸣。</p>
        `,
        excerpt: '一口乡音入戏来！湖南花鼓戏，如何守住湖湘的文化根脉？',
        authorName: '戏曲爱好者',
        authorAvatar: 'https://picsum.photos/seed/user7/100/100',
        createdAt: '2023-05-15',
        views: 432,
        comments: 31,
        likes: 76,
        category: '传统艺术'
      },
      {
        id: '8',
        title: '马王堆汉墓：穿越千年的湘楚文明密码',
        content: `
          <p>马王堆汉墓的惊世发掘，如何改写我们对汉代湘楚文明的认知？</p>
          
          <h3>一、震惊世界的考古发现</h3>
          <p>1972年，长沙马王堆汉墓的发掘震惊世界。这处西汉长沙国丞相利苍及其家属的墓葬，出土了3000余件珍贵文物，成为20世纪中国最重大的考古发现之一。千年不腐的辛追夫人遗体、薄如蝉翼的素纱襌衣（仅重49克）、包罗万象的帛书竹简，每一件都堪称奇迹。</p>
          
          <h3>二、湘楚文明的璀璨见证</h3>
          <p>马王堆汉墓的文物，带着鲜明的湘楚特色，既融入了中原文化的精髓，又保留着南方楚文化的浪漫与灵动。其中，帛书《五十二病方》是中国现存最早的医方著作，素纱襌衣展现了汉代精湛的纺织技术，这些都让我们看到了两千多年前湖湘大地的文明高度。</p>
          
          <h3>三、文物活化的当代探索</h3>
          <p>如今，马王堆汉墓的文物不仅在博物馆里静静陈列，更通过数字技术、文创产品等方式“活”了起来。我们可以探讨，如何让这些沉睡千年的文物更好地与当代对话，让更多人了解并爱上湖湘的古老文明？</p>
        `,
        excerpt: '马王堆汉墓的惊世发掘，如何改写我们对汉代湘楚文明的认知？',
        authorName: '考古爱好者',
        authorAvatar: 'https://picsum.photos/seed/user8/100/100',
        createdAt: '2023-05-10',
        views: 789,
        comments: 65,
        likes: 143,
        category: '历史研究'
      },
      {
        id: '9',
        title: '铜官窑：黑石号上的大唐湘瓷传奇',
        content: `
          <p>长沙铜官窑，为何能成为大唐海上丝绸之路的“陶瓷名片”？</p>
          
          <h3>一、世界釉下多彩陶瓷的发源地</h3>
          <p>长沙铜官窑，坐落于长沙望城区，是世界釉下多彩陶瓷的发源地，也是唐代三大瓷窑之一。它突破了当时青瓷一统天下的局面，首创了釉下多彩工艺，为中国陶瓷史写下了浓墨重彩的一笔。</p>
          
          <h3>二、黑石号沉船的历史见证</h3>
          <p>1998年，印尼海域发现的“黑石号”沉船，让长沙铜官窑再次惊艳世界。船上装载的数万件铜官窑瓷器，见证了大唐王朝通过海上丝绸之路与世界的深度交流。这些瓷器上的诗文、绘画与异域风格的装饰，反映了大唐盛世的开放胸襟。</p>
          
          <h3>三、千年窑火的当代重燃</h3>
          <p>今天的铜官窑古镇，不仅复原了古窑遗址，更成为了集文化体验、休闲旅游于一体的新地标。我们可以思考，如何让铜官窑的“天下第一窑”美誉，在今天焕发新的光彩，成为湖湘文化“走出去”的重要名片？</p>
        `,
        excerpt: '长沙铜官窑，为何能成为大唐海上丝绸之路的“陶瓷名片”？',
        authorName: '陶瓷爱好者',
        authorAvatar: 'https://picsum.photos/seed/user9/100/100',
        createdAt: '2023-05-05',
        views: 654,
        comments: 48,
        likes: 112,
        category: '文化遗产'
      }
    ]
    
    return posts.find(post => post.id === id)
  },
  
  getRelatedPosts(id) {
    // 模拟获取相关帖子
    const posts = [
      {
        id: '2',
        title: '岳麓书院的历史与文化价值',
        authorName: '历史学者',
        createdAt: '2023-06-10'
      },
      {
        id: '3',
        title: '湘绣艺术的魅力与传承',
        authorName: '艺术爱好者',
        createdAt: '2023-06-05'
      }
    ]
    
    return posts.filter(post => post.id !== id)
  },
  
  getCommentsByPostId(postId) {
    // 模拟获取评论数据
    const comments = [
      {
        id: '1',
        content: '非常赞同您的观点！湖湘文化确实需要在传承中创新，才能更好地适应现代社会的发展需求。',
        authorName: '文化爱好者',
        authorAvatar: 'https://picsum.photos/seed/comment1/100/100',
        createdAt: '2023-06-16'
      },
      {
        id: '2',
        content: '我认为湖湘文化中的"经世致用"精神对现代社会仍然具有重要的指导意义。我们应该继承和发扬这种精神。',
        authorName: '传统文化研究者',
        authorAvatar: 'https://picsum.photos/seed/comment2/100/100',
        createdAt: '2023-06-17'
      },
      {
        id: '3',
        content: '期待看到更多关于湖湘文化的讨论和研究！',
        authorName: '学生',
        authorAvatar: 'https://picsum.photos/seed/comment3/100/100',
        createdAt: '2023-06-18'
      }
    ]
    
    return comments
  }
}

export default {
  name: 'PostDetailPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    }
  },
  setup(props) {
    const router = useRouter()
    const route = useRoute()
    
    const postId = route.params.id
    const currentPost = ref(null)
    const relatedPosts = ref([])
    const comments = ref([])
    const newComment = ref('')
    const isLoggedIn = ref(false)
    const user = ref(null)
    const liked = ref(false)
    const isUserAuthor = ref(false)
    const hasMoreComments = ref(false)
    
    // 获取帖子详情
    const fetchPostDetail = () => {
      const post = mockPostService.getPostById(postId)
      if (post) {
        currentPost.value = post
        
        // 更新浏览量
        currentPost.value.views += 1
        
        // 获取相关帖子
        relatedPosts.value = mockPostService.getRelatedPosts(postId)
        
        // 获取评论
        comments.value = mockPostService.getCommentsByPostId(postId)
        
        // 检查是否已登录
        const userInfo = localStorage.getItem('user')
        if (userInfo) {
          isLoggedIn.value = true
          user.value = JSON.parse(userInfo)
          
          // 检查是否是作者
          isUserAuthor.value = user.value.username === currentPost.value.authorName
          
          // 检查是否已点赞（模拟）
          const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
          liked.value = likedPosts.includes(postId)
        }
      } else {
        // 如果帖子不存在，显示提示
        if (props.showAlert) {
          props.showAlert('帖子不存在或已被删除', 'error')
        }
        // 重定向回社区页面
        setTimeout(() => {
          router.push('/community')
        }, 1500)
      }
    }
    
    // 格式化帖子内容
    const formatPostContent = (content) => {
      if (!content) return ''
      return content.trim()
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }
    
    // 发表评论
    const submitComment = () => {
      if (!newComment.value.trim()) {
        if (props.showAlert) {
          props.showAlert('评论内容不能为空', 'warning')
        }
        return
      }
      
      // 创建新评论
      const newCommentObj = {
        id: Date.now().toString(),
        content: newComment.value.trim(),
        authorName: user.value.username,
        authorAvatar: `https://picsum.photos/seed/${user.value.username}/100/100`,
        createdAt: new Date().toISOString().split('T')[0]
      }
      
      // 添加到评论列表
      comments.value.unshift(newCommentObj)
      
      // 更新评论数
      if (currentPost.value) {
        currentPost.value.comments += 1
      }
      
      // 清空评论框
      newComment.value = ''
      
      // 显示成功提示
      if (props.showAlert) {
        props.showAlert('评论发表成功', 'success')
      }
    }
    
    // 切换点赞状态
    const toggleLike = () => {
      if (!isLoggedIn.value) {
        if (props.showAlert) {
          props.showAlert('请先登录后点赞', 'info')
        }
        return
      }
      
      liked.value = !liked.value
      
      // 更新点赞数
      if (currentPost.value) {
        currentPost.value.likes += liked.value ? 1 : -1
      }
      
      // 保存点赞状态（模拟）
      const likedPosts = JSON.parse(localStorage.getItem('likedPosts') || '[]')
      if (liked.value) {
        likedPosts.push(postId)
      } else {
        const index = likedPosts.indexOf(postId)
        if (index > -1) {
          likedPosts.splice(index, 1)
        }
      }
      localStorage.setItem('likedPosts', JSON.stringify(likedPosts))
    }
    
    // 跳转到其他帖子详情
    const goToPostDetail = (id) => {
      router.push(`/post-detail/${id}`)
    }
    
    // 返回社区页面
    const goBack = () => {
      router.push('/community')
    }
    
    // 加载更多评论
    const loadMoreComments = () => {
      // 模拟加载更多评论
      if (props.showAlert) {
        props.showAlert('已加载全部评论', 'info')
      }
      hasMoreComments.value = false
    }
    
    // 页面加载时获取数据
    onMounted(() => {
      fetchPostDetail()
      // 滚动到页面顶部
      window.scrollTo({ top: 0, behavior: 'smooth' })
    })
    
    return {
      currentPost,
      relatedPosts,
      comments,
      newComment,
      isLoggedIn,
      liked,
      isUserAuthor,
      hasMoreComments,
      formatPostContent,
      formatDate,
      submitComment,
      toggleLike,
      goToPostDetail,
      goBack,
      loadMoreComments
    }
  }
}
</script>

<style scoped>
.post-detail-page {
  padding: 2rem 0;
}

.back-button-container {
  position: fixed;
  top: 2rem;
  left: 2rem;
  z-index: 1000;
}

.back-button {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 25px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  font-weight: 500;
  font-size: 1rem;
}

.back-button:hover {
  background-color: var(--primary-dark);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  transform: translateX(-5px) scale(1.05);
  animation: backButtonFloat 1.5s ease-in-out infinite;
}

@keyframes backButtonFloat {
  0%, 100% {
    transform: translateX(0) scale(1.05);
  }
  50% {
    transform: translateX(5px) scale(1.05);
  }
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.post-detail {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  margin-bottom: 2rem;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #eee;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.author-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: bold;
  color: #333;
  font-size: 1.1rem;
}

.post-date {
  font-size: 0.9rem;
  color: #999;
}

.post-stats {
  display: flex;
  gap: 1.5rem;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.9rem;
  color: #666;
}

.post-title {
  font-size: 2rem;
  color: #333;
  margin: 0 0 1.5rem 0;
  line-height: 1.3;
}

.post-body {
  color: #333;
  line-height: 1.8;
  margin-bottom: 2rem;
}

.post-body p {
  margin-bottom: 1rem;
}

.post-body h3 {
  font-size: 1.5rem;
  color: #333;
  margin: 1.5rem 0 1rem 0;
}

.post-tags {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.post-category {
  background-color: #f0f0f0;
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.9rem;
  color: #666;
}

.post-actions {
  display: flex;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.action-btn {
  background-color: transparent;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  cursor: pointer;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s;
}

.action-btn:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

.action-btn i.liked {
  color: var(--primary-color);
}

.related-posts {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.related-posts h3 {
  margin-top: 0;
  margin-bottom: 1rem;
  color: #333;
  font-size: 1.3rem;
}

.related-posts-list {
  display: grid;
  gap: 1rem;
}

.related-post-item {
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.related-post-item:hover {
  background-color: #f8f9fa;
  border-color: #ddd;
}

.related-post-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: #333;
}

.related-post-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.85rem;
  color: #999;
}

.comments-section {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.comments-section h3 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #333;
  font-size: 1.3rem;
}

.comment-form {
  margin-bottom: 2rem;
}

.comment-form textarea {
  width: 100%;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: vertical;
  font-size: 1rem;
  margin-bottom: 1rem;
}

.submit-comment-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  font-size: 1rem;
  transition: background-color 0.3s;
}

.submit-comment-btn:hover {
  background-color: var(--primary-dark);
}

.login-prompt {
  padding: 1.5rem;
  background-color: #f8f9fa;
  border-radius: 4px;
  text-align: center;
  margin-bottom: 2rem;
  color: #666;
}

.login-prompt a {
  color: var(--primary-color);
  text-decoration: none;
}

.login-prompt a:hover {
  text-decoration: underline;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.comment-item {
  display: flex;
  gap: 1rem;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.comment-author {
  font-weight: bold;
  color: #333;
}

.comment-date {
  font-size: 0.85rem;
  color: #999;
}

.comment-text {
  color: #333;
  line-height: 1.6;
  margin: 0;
}

.load-more {
  text-align: center;
  margin-top: 2rem;
}

.load-more-btn {
  background-color: transparent;
  border: 1px solid #ddd;
  padding: 0.75rem 1.5rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.load-more-btn:hover {
  background-color: #f8f9fa;
  border-color: #ccc;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .post-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .post-title {
    font-size: 1.5rem;
  }
  
  .post-actions {
    flex-wrap: wrap;
  }
  
  .comment-item {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .comment-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }
}
</style>