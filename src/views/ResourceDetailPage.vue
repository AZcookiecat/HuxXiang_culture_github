<template>
  <div class="resource-detail-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i> 返回
      </button>
    </div>

    <!-- 资源详情头部 -->
    <div class="resource-header">
      <div class="container">
        <div class="breadcrumb">
          <a href="#" @click.prevent="goBack">首页</a> &gt;
          <a href="#" @click.prevent="goToResources">文化资源</a> &gt;
          <span>{{ resource.category }}</span> &gt;
          <span>{{ resource.title }}</span>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="container">
      <div class="resource-content">
        <!-- 左侧：资源详情 -->
        <div class="resource-main">
          <!-- 资源标题和元信息 -->
          <div class="resource-title-section">
            <h1>{{ resource.title }}</h1>
            <div class="resource-meta">
              <span class="meta-item"><i class="fas fa-tag"></i> {{ resource.category }}</span>
            </div>
          </div>

          <!-- 资源图片轮播 -->
          <div class="resource-gallery">
            <div class="gallery-main">
              <img :src="currentImage" :alt="resource.title" class="main-image" loading="lazy" />
            </div>
            <div class="gallery-thumbs" v-if="resource.images && resource.images.length > 1">
              <div 
                v-for="(img, index) in resource.images" 
                :key="index" 
                class="thumb-item"
                :class="{ active: currentImageIndex === index }"
                @click="setCurrentImage(index)"
              >
                <img :src="img" :alt="`${resource.title} ${index + 1}`" loading="lazy" />
              </div>
            </div>
          </div>

          <!-- 资源简介 -->
          <div class="resource-intro">
            <h2>简介</h2>
            <p>{{ resource.introduction }}</p>
          </div>

          <!-- 资源详细内容 -->
          <div class="resource-detail-content">
            <h2>详细介绍</h2>
            <div class="content-paragraph" v-for="(paragraph, index) in detailParagraphs" :key="index">
              <h3 v-if="paragraph.title">{{ paragraph.title }}</h3>
              <p v-for="(text, textIndex) in paragraph.texts" :key="textIndex">{{ text }}</p>
              <img v-if="paragraph.image" :src="paragraph.image" :alt="paragraph.title" class="content-image" loading="lazy" />
            </div>
          </div>

          <!-- 评论区 -->
          <div class="comments-section">
            <h2>用户评论 ({{ comments.length }})</h2>
            <div class="no-comments">
              <p>暂无评论</p>
            </div>
          </div>
        </div>

        <!-- 右侧：侧边栏 -->
        <div class="resource-sidebar">
          <!-- 资源信息卡片 -->
          <div class="info-card">
            <h3>资源信息</h3>
            <div class="info-item">
              <span class="info-label">分类：</span>
              <span class="info-value">{{ resource.category }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">年代：</span>
              <span class="info-value">{{ resource.era }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">来源：</span>
              <span class="info-value">{{ resource.source }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">收藏地：</span>
              <span class="info-value">{{ resource.location }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">保护级别：</span>
              <span class="info-value">{{ resource.protectionLevel }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <button class="action-btn primary" @click="toggleLikeResource">
              <i :class="resource.isLiked ? 'fas fa-thumbs-up' : 'far fa-thumbs-up' "></i>
              <span>{{ resource.isLiked ? '已点赞' : '点赞' }}</span>
            </button>
            <button class="action-btn" @click="shareResource">
              <i class="fas fa-share-alt"></i>
              <span>分享</span>
            </button>
            <button class="action-btn" @click="reportResource">
              <i class="fas fa-flag"></i>
              <span>举报</span>
            </button>
          </div>

          <!-- 资源贡献者 -->
          <div class="contributor-card">
            <h3>资源贡献者</h3>
            <div class="contributor-info">
              <div class="contributor-avatar">
                <img :src="resource.contributorAvatar" :alt="resource.contributorName" loading="lazy" />
              </div>
              <div class="contributor-details">
                <h4>{{ resource.contributorName }}</h4>
                <p>{{ resource.contributorRole }}</p>
              </div>
            </div>
            <button class="follow-btn" v-if="!resource.isFollowing" @click="followContributor">
              关注
            </button>
            <button class="following-btn" v-else @click="unfollowContributor">
              已关注
            </button>
          </div>

          <!-- 热门标签 -->
          <div class="tags-card">
            <h3>热门标签</h3>
            <div class="tags-list">
              <a href="#" @click.prevent="filterByTag(tag)" v-for="tag in resource.tags" :key="tag" class="tag-item">
                {{ tag }}
              </a>
            </div>
          </div>


        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

export default {
  name: 'ResourceDetailPage',
  props: {
    showAlert: {
      type: Function,
      default: null
    },
    showLoginModal: {
      type: Function,
      default: null
    },
    user: {
      type: Object,
      default: () => ({
        isLoggedIn: false,
        username: '',
        avatar: ''
      })
    }
  },
  setup(props) {
    const route = useRoute()
    const router = useRouter()
    
    // 从路由参数获取资源ID
    const resourceId = route.params.id || '1'
    
    // 资源数据仓库
    const resourcesData = {
      '1': {
        id: '1',
        title: '岳麓书院',
        category: '历史遗迹',
        introduction: '岳麓书院是中国历史上赫赫闻名的四大书院之一，坐落于中国历史文化名城湖南长沙湘江西岸的岳麓山脚下，作为世界上最古老的学府之一，其古代传统的书院建筑至今被完整保存，每一组院落、每一块石碑、每一枚砖瓦、每一支风荷，都闪烁着时光淬炼的人文精神。',
        images: [
          '/src/assets/Imgs/岳麓书院1.jpg',
          '/src/assets/Imgs/岳麓书院2.jpg',
          '/src/assets/Imgs/岳麓书院3.jpg',
        ],
        isLiked: false,
        era: '北宋',
        source: '湖南省博物馆',
        location: '湖南省长沙市岳麓区',
        protectionLevel: '全国重点文物保护单位',
        contributorName: '李文化',
        contributorAvatar: 'https://picsum.photos/seed/contributor/100/100',
        contributorRole: '湖湘文化研究员',
        isFollowing: false,
        tags: ['岳麓书院', '历史建筑', '四大书院', '北宋', '长沙']
      },
      '2': {
        id: '2',
        title: '湘绣',
        category: '传统艺术',
        introduction: '湘绣是中国四大名绣之一，以其精湛的技艺和独特的风格闻名于世。湘绣起源于湖南长沙，历史悠久，工艺精湛，针法细腻，色彩鲜明，被誉为"东方明珠"。',
        images: [
          '/src/assets/Imgs/湘绣1.jpg',
          '/src/assets/Imgs/湘绣2.jpg',
          '/src/assets/Imgs/湘绣3.jpg'
        ],
        isLiked: false,
        era: '清朝',
        source: '湖南省工艺美术协会',
        location: '湖南省长沙市',
        protectionLevel: '国家级非物质文化遗产',
        contributorName: '王绣师',
        contributorAvatar: 'https://picsum.photos/seed/artist1/100/100',
        contributorRole: '湘绣大师',
        isFollowing: false,
        tags: ['湘绣', '传统艺术', '刺绣', '非物质文化遗产', '湖南']
      },
      '3': {
        id: '3',
        title: '岳阳楼',
        category: '建筑风格',
        introduction: '岳阳楼位于湖南省岳阳市古城西门城墙之上，下瞰洞庭，前望君山，自古有"洞庭天下水，岳阳天下楼"之美誉。北宋范仲淹脍炙人口的《岳阳楼记》更使岳阳楼著称于世。',
        images: [
          '/src/assets/Imgs/岳阳楼记1.jpg',
          '/src/assets/Imgs/岳阳楼记2.jpg',
          '/src/assets/Imgs/岳阳楼记3.jpg'
        ],
        isLiked: false,
        era: '唐代',
        source: '岳阳市博物馆',
        location: '湖南省岳阳市',
        protectionLevel: '全国重点文物保护单位',
        contributorName: '张建筑',
        contributorAvatar: 'https://picsum.photos/seed/architect1/100/100',
        contributorRole: '古建筑研究员',
        isFollowing: false,
        tags: ['岳阳楼', '江南三大名楼', '古建筑', '洞庭湖', '范仲淹']
      },
      '4': {
        id: '4',
        title: '桃花源记',
        category: '文学作品',
        introduction: '《桃花源记》是东晋文学家陶渊明的传世散文，以武陵渔人的奇遇为线索，描绘了一个与世隔绝、耕织自足、安宁祥和的理想社会。湖南常德的桃花源景区，正是依托这一文学名篇而建，成为承载着中国人"世外桃源"精神向往的文化地标。',
        images: [
          '/src/assets/Imgs/桃花源记1.jpg',
          '/src/assets/Imgs/桃花源记2.jpg',
          '/src/assets/Imgs/桃花源记3.jpg'
        ],
        isLiked: false,
        era: '东晋',
        source: '常德市博物馆',
        location: '湖南省常德市',
        protectionLevel: '国家5A级景区',
        contributorName: '王文学',
        contributorAvatar: 'https://picsum.photos/seed/writer1/100/100',
        contributorRole: '文学研究员',
        isFollowing: false,
        tags: ['桃花源记', '陶渊明', '文学作品', '常德', '世外桃源']
      },
      '5': {
        id: '5',
        title: '湘菜',
        category: '饮食文化',
        introduction: '湘菜是中国八大菜系之一，以"鲜、香、辣、酸"为核心口味，擅长以辣椒、豆豉、剁椒等调味，代表菜品有剁椒鱼头、小炒黄牛肉、东安子鸡等。它不仅是湖湘饮食文化的核心载体，更以独特的风味和深厚的文化底蕴享誉全国。',
        images: [
          '/src/assets/Imgs/湘菜1.jpg',
          '/src/assets/Imgs/湘菜2.jpeg',
          '/src/assets/Imgs/湘菜3.jpg'
        ],
        isLiked: false,
        era: '战国时期',
        source: '湖南省餐饮协会',
        location: '湖南省',
        protectionLevel: '国家级非物质文化遗产',
        contributorName: '李厨师',
        contributorAvatar: 'https://picsum.photos/seed/chef1/100/100',
        contributorRole: '湘菜大师',
        isFollowing: false,
        tags: ['湘菜', '饮食文化', '八大菜系', '辣椒', '湖南美食']
      },
      '6': {
        id: '6',
        title: '端午节（湖南汨罗）',
        category: '民俗风情',
        introduction: '湖南汨罗是端午节的重要发源地，作为爱国诗人屈原的投江之地，这里保留着祭龙、竞渡、食粽、点雄黄等最具传统特色的端午习俗，是中国端午文化的核心传承地，也是全球华人追思屈原、传承家国情怀的精神家园。',
        images: [
          '/src/assets/Imgs/端午1.jpeg',
          '/src/assets/Imgs/端午2.jpeg',
          '/src/assets/Imgs/端午3.jpeg'
        ],
        isLiked: false,
        era: '战国时期',
        source: '汨罗市文化旅游局',
        location: '湖南省汨罗市',
        protectionLevel: '世界非物质文化遗产',
        contributorName: '张民俗',
        contributorAvatar: 'https://picsum.photos/seed/folk1/100/100',
        contributorRole: '民俗研究员',
        isFollowing: false,
        tags: ['端午节', '汨罗', '屈原', '龙舟竞渡', '非物质文化遗产']
      },
      '7': {
        id: '7',
        title: '衡山（南岳）',
        category: '历史遗迹',
        introduction: '衡山又名南岳，是中国五岳之一，位于湖南衡阳市境内。它既是道教主流全真派的圣地，也是佛教禅宗的重要道场，以"寿岳""宗教圣地"之称闻名于世，是一座兼具自然景观与人文底蕴的文化名山。',
        images: [
          '/src/assets/Imgs/衡山1.jpeg',
          '/src/assets/Imgs/衡山2.jpeg',
          '/src/assets/Imgs/衡山3.jpeg'
        ],
        isLiked: false,
        era: '尧舜时期',
        source: '衡阳市文化旅游局',
        location: '湖南省衡阳市',
        protectionLevel: '国家5A级景区',
        contributorName: '刘宗教',
        contributorAvatar: 'https://picsum.photos/seed/religion1/100/100',
        contributorRole: '宗教文化研究员',
        isFollowing: false,
        tags: ['衡山', '南岳', '五岳', '道教', '佛教', '寿岳']
      },
      '8': {
        id: '8',
        title: '花鼓戏',
        category: '传统艺术',
        introduction: '花鼓戏是湖南最具代表性的地方戏曲剧种，起源于民间歌舞，以活泼明快的曲调、幽默诙谐的表演著称，代表剧目有《刘海砍樵》《打铜锣》《补锅》等。它扎根于湖湘民间，是普通民众表达喜怒哀乐的重要艺术形式。',
        images: [
          '/src/assets/Imgs/花鼓戏1.jpeg',
          '/src/assets/Imgs/花鼓戏2.jpeg',
          '/src/assets/Imgs/花鼓戏3.jpeg'
        ],
        isLiked: false,
        era: '清代中晚期',
        source: '湖南省戏曲家协会',
        location: '湖南省',
        protectionLevel: '国家级非物质文化遗产',
        contributorName: '陈戏曲',
        contributorAvatar: 'https://picsum.photos/seed/actor1/100/100',
        contributorRole: '花鼓戏表演艺术家',
        isFollowing: false,
        tags: ['花鼓戏', '传统艺术', '地方戏曲', '刘海砍樵', '非物质文化遗产']
      },
      '9': {
        id: '9',
        title: '曾国藩故居（富厚堂）',
        category: '历史遗迹',
        introduction: '曾国藩故居富厚堂位于湖南双峰县荷叶镇，是晚清名臣曾国藩的侯府。它集住宅、藏书楼、祠堂于一体，是中国近代规模最大、保存最完整的私家建筑群之一，被誉为"中国乡间第一侯府"。',
        images: [
          '/src/assets/Imgs/曾国藩故居1.jpeg',
          '/src/assets/Imgs/曾国藩故居2.jpeg',
          '/src/assets/Imgs/曾国藩故居3.jpeg'
        ],
        isLiked: false,
        era: '清代',
        source: '双峰县文化旅游局',
        location: '湖南省娄底市双峰县',
        protectionLevel: '全国重点文物保护单位',
        contributorName: '曾后裔',
        contributorAvatar: 'https://picsum.photos/seed/descendant1/100/100',
        contributorRole: '曾国藩研究专家',
        isFollowing: false,
        tags: ['曾国藩', '富厚堂', '历史遗迹', '晚清', '双峰县']
      },
      '10': {
        id: '10',
        title: '马王堆汉墓',
        category: '历史遗迹',
        introduction: '马王堆汉墓位于长沙市芙蓉区，是西汉初期长沙国丞相轪侯利苍及其家属的墓葬。它出土了千年不腐的辛追夫人遗体、素纱襌衣、马王堆汉墓帛书等国宝级文物，被誉为"汉初历史的百科全书"，震惊了世界考古界。',
        images: [
          '/src/assets/Imgs/马王堆1.jpeg',
          '/src/assets/Imgs/马王堆2.jpeg',
          '/src/assets/Imgs/马王堆3.jpeg'
        ],
        isLiked: false,
        era: '西汉初期',
        source: '湖南博物院',
        location: '湖南省长沙市芙蓉区',
        protectionLevel: '全国重点文物保护单位',
        contributorName: '赵考古',
        contributorAvatar: 'https://picsum.photos/seed/archaeologist1/100/100',
        contributorRole: '考古研究员',
        isFollowing: false,
        tags: ['马王堆汉墓', '西汉', '辛追夫人', '素纱襌衣', '考古发现']
      },
      '11': {
        id: '11',
        title: '火宫殿',
        category: '饮食文化',
        introduction: '火宫殿是长沙著名的民俗文化地标，始建于清乾隆十二年（1747年），原为祭祀火神祝融的庙宇。如今，它已成为集湘菜、小吃、戏曲、民俗于一体的"湘文化活化石"，是品尝长沙美食、体验湖湘民俗的必去之地。',
        images: [
          '/src/assets/Imgs/火宫殿1.jpeg',
          '/src/assets/Imgs/火宫殿2.jpeg',
          '/src/assets/Imgs/火宫殿3.jpeg'
        ],
        isLiked: false,
        era: '清代',
        source: '长沙市文化旅游局',
        location: '湖南省长沙市',
        protectionLevel: '市级文物保护单位',
        contributorName: '孙美食',
        contributorAvatar: 'https://picsum.photos/seed/food1/100/100',
        contributorRole: '长沙美食专家',
        isFollowing: false,
        tags: ['火宫殿', '长沙美食', '臭豆腐', '民俗文化', '火神祝融']
      },
      '12': {
        id: '12',
        title: '铜官窑',
        category: '传统艺术',
        introduction: '铜官窑位于长沙市望城区，是唐代南方最大的青瓷窑场，也是世界釉下彩瓷的发源地。它开创了中国陶瓷彩绘的先河，产品曾通过海上丝绸之路远销至东亚、西亚乃至非洲，是唐代湖南对外文化交流的重要见证。',
        images: [
          '/src/assets/Imgs/铜官窑1.jpeg',
          '/src/assets/Imgs/铜官窑2.jpeg',
          '/src/assets/Imgs/铜官窑3.jpeg'
        ],
        isLiked: false,
        era: '唐代',
        source: '长沙市望城区文化旅游局',
        location: '湖南省长沙市望城区',
        protectionLevel: '全国重点文物保护单位',
        contributorName: '周陶瓷',
        contributorAvatar: 'https://picsum.photos/seed/ceramic1/100/100',
        contributorRole: '陶瓷艺术研究员',
        isFollowing: false,
        tags: ['铜官窑', '唐代', '釉下彩瓷', '海上丝绸之路', '陶瓷艺术']
      },
      '13': {
        id: '13',
        title: '曾国藩',
        category: '历史伟人',
        introduction: '曾国藩（1811年11月26日－1872年3月12日），字伯涵，号涤生，湖南长沙府湘乡县人，中国清朝后期政治家、战略家、理学家、文学家，湘军的创立者和统帅。',
        images: [
          '/src/assets/Imgs/曾国藩1.jpeg',
          '/src/assets/Imgs/曾国藩2.jpeg',
          '/src/assets/Imgs/曾国藩3.jpeg'
        ],
        isLiked: false,
        era: '清朝',
        source: '湖南省社会科学院',
        location: '湖南省娄底市双峰县',
        protectionLevel: '历史文化名人',
        contributorName: '李历史',
        contributorAvatar: 'https://picsum.photos/seed/historian1/100/100',
        contributorRole: '清史研究员',
        isFollowing: false,
        tags: ['曾国藩', '湘军', '晚清', '政治家', '文学家', '湖南']
      }
    }
    
    // 默认资源数据
    const defaultResource = {
      id: resourceId,
      title: '资源不存在',
      category: '未知',
      introduction: '抱歉，您访问的资源不存在或已被删除。',
      images: [
        'https://picsum.photos/seed/default/1200/600'
      ],
      isLiked: false,
      era: '',
      source: '',
      location: '',
      protectionLevel: '',
      contributorName: '',
      contributorAvatar: '',
      contributorRole: '',
      isFollowing: false,
      tags: []
    }
    
    // 当前资源数据
    const resource = ref(defaultResource)
    
    // 当前图片索引
    const currentImageIndex = ref(0)
    
    // 评论相关 - 设置为没有评论
    const comments = ref([])
    
    // 保留相关变量以避免报错
    const newComment = ref('')
    const isSubmittingComment = ref(false)
    const replyToId = ref(null)
    const replyContent = ref('')
    const isSubmittingReply = ref(false)
    
    // 资源详细内容段落库
    const detailParagraphsMap = {
      '1': [
        {
          title: '历史沿革',
          texts: [
            '岳麓书院始建于北宋开宝九年（976年），潭州太守朱洞在僧人办学的基础上，由官府捐资兴建，正式创立岳麓书院。',
            '北宋祥符八年（公元1015年），宋真宗召见岳麓山长周式，御笔赐书"岳麓书院"四字门额。',
            '嗣后，历经南宋、元、明、清各代，至清末光绪廿九年（公元1903年），岳麓书院与湖南省城大学堂合并改制为湖南高等学堂，沿用书院旧址。',
            '中华民国15年（公元1926年），湖南高等学堂正式定名湖南大学，仍就书院基址扩建至今。'
          ],
          image: '/src/assets/Imgs/岳麓书院1.jpg'
        },
        {
          title: '建筑特色',
          texts: [
            '岳麓书院古建筑群分为教学、藏书、祭祀、园林、纪念五大建筑格局。岳麓书院主体建筑面积有31000多平方米，分为书院主体、附属文庙及新建的中国书院博物馆。',
            '岳麓书院占地面积25000平方米，现存建筑大部分为明清遗物，其古建筑在布局上采用中轴对称、纵深多进的院落形式。主体建筑如头门、大门、二门、讲堂、御书楼集中于中轴线上，讲堂布置在中轴线的中央。斋舍、祭祀专祠等排列于两旁。',
            '中轴对称、层层递进的院落，除了营造一种庄严、神妙、幽远的纵深感和视觉效应之外，还体现了儒家文化尊卑有序、等级有别、主次鲜明的社会伦理关系。'
          ],
          image: '/src/assets/Imgs/岳麓书院2.jpg'
        },
        {
          title: '文化价值',
          texts: [
            '岳麓书院不仅是湖南大学的文史哲人才培养和研究基地，湖南省旅游胜地，更是是整个长沙市的文化窗口和文化名片。',
            '书院的学规、课程、考试、学田、学舍、藏书等都成为了中国古代书院的典范。它的教育模式和理念，对后世产生了深远的影响。',
            '作为世界上最古老的学府之一，岳麓书院的古代传统的书院建筑至今被完整保存，每一组院落、每一块石碑、每一枚砖瓦、每一支风荷，都闪烁着时光淬炼的人文精神。'
          ]
        },
        {
          title: '著名学者',
          texts: [
            '千年以来，岳麓书院人才辈出，毛泽东、蔡和森、何叔衡、李达等都曾在此求学。',
            '南宋时期，朱熹、张栻曾在此讲学，形成了著名的"朱张会讲"，推动了宋代理学的发展。',
            '明代王阳明心学也在此得到传播和发展，进一步丰富了岳麓书院的学术内涵。'
          ],
          image: '/src/assets/Imgs/岳麓书院3.jpg'
        }
      ],
      '2': [
        {
          title: '历史渊源',
          texts: [
            '湘绣的历史可以追溯到2000多年前的春秋战国时期，湖南地区的刺绣工艺已经相当发达。',
            '到了清朝中后期，湘绣逐渐形成了自己独特的风格，并与苏绣、粤绣、蜀绣并称为中国四大名绣。',
            '20世纪初，湘绣迎来了黄金时期，不仅在国内享有盛誉，还远销海外，成为中国文化的重要使者。'
          ],
          image: '/src/assets/Imgs/湘绣1.jpg'
        },
        {
          title: '艺术特色',
          texts: [
            '湘绣以其精湛的针法、丰富的色彩和逼真的表现力著称。常用的针法有掺针、毛针、游针、齐针等70多种。',
            '湘绣的色彩运用大胆而细腻，讲究"色相分明，色阶渐变"，通过不同颜色的丝线相互搭配，形成丰富的层次感和立体感。',
            '湘绣的题材广泛，包括花鸟、山水、人物、走兽等，其中以狮、虎等动物题材最为著名，有"湘绣狮虎冠天下"的美誉。'
          ],
          image: '/src/assets/Imgs/湘绣2.jpg'
        },
        {
          title: '传承与发展',
          texts: [
            '2006年，湘绣被列入第一批国家级非物质文化遗产名录。近年来，政府和社会各界加大了对湘绣的保护和传承力度。',
            '新一代湘绣艺人在继承传统的基础上，不断创新，将现代艺术元素融入传统湘绣，开发出了许多符合现代人审美需求的新产品。',
            '湘绣不仅是一种传统工艺，更是湖湘文化的重要组成部分，它承载着湖南人民的智慧和创造力，是中华民族优秀传统文化的瑰宝。'
          ]
        }
      ],
      '3': [
        {
          title: '历史变迁',
          texts: [
            '岳阳楼始建于东汉末年，最初是鲁肃为训练水师而建的阅兵台。',
            '唐代开元四年（716年），中书令张说被贬为岳州刺史，在此基础上重建楼阁，正式定名为"岳阳楼"。',
            '北宋庆历四年（1044年），滕子京谪守巴陵郡，重修岳阳楼，并请范仲淹作《岳阳楼记》，使岳阳楼声名远扬。',
            '此后，岳阳楼历经多次损毁和重建，现存建筑为清光绪六年（1880年）重建。'
          ],
          image: '/src/assets/Imgs/岳阳楼记1.jpg'
        },
        {
          title: '建筑艺术',
          texts: [
            '岳阳楼建筑风格独特，主楼高19.42米，为三层、四柱、飞檐、盔顶、纯木结构。',
            '其盔顶造型在中国古代建筑中独一无二，充分体现了古代劳动人民的智慧和创造力。',
            '楼内藏有大量历代名人的诗词、书画和碑刻，其中最著名的当属范仲淹的《岳阳楼记》木雕屏。'
          ],
          image: '/src/assets/Imgs/岳阳楼记2.jpg'
        },
        {
          title: '文化意义',
          texts: [
            '岳阳楼因范仲淹的《岳阳楼记》而名满天下，文中"先天下之忧而忧，后天下之乐而乐"的名句，成为中华民族精神的重要象征。',
            '岳阳楼与湖北武汉黄鹤楼、江西南昌滕王阁并称为"江南三大名楼"，是中国古代建筑的杰出代表。',
            '作为湖湘文化的重要载体，岳阳楼不仅是一座历史建筑，更是中华民族优秀传统文化的重要象征。'
          ]
        }
      ],
      '4': [
        {
          title: '历史沿革',
          texts: [
            '自唐代起，常德桃花源便依托《桃花源记》的记载逐步开发，建有桃川书院、渊明祠等人文景观，吸引历代文人墨客前来凭吊题咏。宋、元、明、清各朝代均对景区进行过修缮与扩建，使其成为江南著名的隐逸文化圣地。1990年以来，当地政府对桃花源进行了大规模保护性开发，恢复了"秦人村""豁然台"等核心景点，使其成为融合自然山水与人文底蕴的国家5A级景区。'
          ],
          image: '/src/assets/Imgs/桃花源记1.jpg'
        },
        {
          title: '文化价值',
          texts: [
            '它不仅是中国隐逸文化与田园思想的集中体现，更承载着中国人对理想社会的永恒追求。文中"黄发垂髫，并怡然自乐"的描写，深刻影响了中国文人的精神世界与文学创作，成为后世描绘乌托邦社会的经典范式。'
          ],
          image: '/src/assets/Imgs/桃花源记2.jpg'
        },
        {
          title: '传承发展',
          texts: [
            '近年来，景区通过举办"中国桃花源文化节""桃花节"等活动，结合实景演出《桃花源记》，让游客沉浸式体验文学中的场景。同时，开发了线上虚拟游览项目与数字文创产品，吸引年轻群体了解这一文化IP，实现了传统文化的创造性转化与创新性发展。'
          ],
          image: '/src/assets/Imgs/桃花源记3.jpg'
        }
      ],
      '5': [
        {
          title: '历史沿革',
          texts: [
            '湘菜的历史可追溯至战国时期，楚地饮食文化中"食不厌精"的传统奠定了其发展基础。历经千年演变，在明清时期逐渐形成"一菜一格，百菜百味"的成熟体系。近代以来，随着湖南人口的迁徙与文化交流，湘菜进一步融合各地风味，尤其是吸收了粤菜的烹饪技法与川菜的调味理念，逐步走向全国乃至世界。'
          ],
          image: '/src/assets/Imgs/湘菜1.jpg'
        },
        {
          title: '艺术特色',
          texts: [
            '湘菜注重刀工与火候，讲究"色、香、味、形、器"的和谐统一。烹饪技法多样，包括煨、炖、腊、蒸、炒等，其中"煨"能使食材软烂入味，"腊"则赋予食材独特的烟熏香气。调味上善用"酸辣"，以辣椒提香、醋酸解腻，既保持了食材的本味，又能激发复合口感。'
          ],
          image: '/src/assets/Imgs/湘菜2.jpeg'
        },
        {
          title: '文化价值',
          texts: [
            '湘菜的"辣"与湖南人"霸蛮""坚韧"的性格特质高度呼应，每一道菜品都蕴含着地域民俗与历史典故。例如"毛氏红烧肉"承载着红色记忆，"祖庵菜"则体现了晚清湖湘士大夫的饮食审美。它不仅是了解湖湘文化的重要窗口，更成为连接全球湖南人的味觉纽带。'
          ],
          image: '/src/assets/Imgs/湘菜3.jpg'
        }
      ],
      '6': [
        {
          title: '历史沿革',
          texts: [
            '自战国时期屈原投汨罗江后，当地百姓便兴起了龙舟竞渡、投粽祭贤的习俗，以寄托对屈原的哀思。这一习俗历经秦、汉、唐、宋等朝代的传承与发展，逐渐形成了一套完整的祭祀仪式与节庆活动。明清时期，汨罗端午龙舟竞渡已成为湘楚地区规模最大的民间盛会，吸引周边数十万群众参与。'
          ],
          image: '/src/assets/Imgs/端午1.jpeg'
        },
        {
          title: '文化意义',
          texts: [
            '汨罗端午不仅承载着对爱国诗人屈原的缅怀，更蕴含着驱邪避灾、祈福安康的文化内涵。龙舟竞渡所展现的集体协作精神，与屈原"上下求索"的爱国精神相互呼应，成为中华民族精神的重要象征。2009年，中国端午节（汨罗）被联合国教科文组织列入《人类非物质文化遗产代表作名录》，进一步提升了其全球影响力。'
          ],
          image: '/src/assets/Imgs/端午2.jpeg'
        },
        {
          title: '传承发展',
          texts: [
            '当地通过举办"中国汨罗江国际龙舟节""屈原文化研讨会"等活动，推动端午文化的活态传承。同时，将端午习俗纳入中小学地方课程，开发"端午香囊""龙舟模型"等文创产品，让传统节日在当代社会中焕发新的生机与活力。此外，汨罗还与韩国、日本等国开展端午文化交流，促进了东亚文化圈的对话与互鉴。'
          ],
          image: '/src/assets/Imgs/端午3.jpeg'
        }
      ],
      '7': [
        {
          title: '历史沿革',
          texts: [
            '衡山的祭祀历史可追溯至尧舜时期，相传舜帝曾南巡至此祭祀山神。汉代以来，道教、佛教相继在此兴盛，形成了"佛道共存""寺观同山"的独特格局。唐代，衡山被封为"南岳真君"，成为官方祭祀的重要场所；宋代，宋徽宗御题"天下南岳"，进一步提升了其地位。历代帝王曾多次前来封禅祭祀，文人墨客如李白、杜甫、朱熹等也在此留下了大量诗词题咏。'
          ],
          image: '/src/assets/Imgs/衡山1.jpeg'
        },
        {
          title: '建筑特色',
          texts: [
            '山上的寺庙道观依山而建，巧妙融合了南方建筑的精巧灵秀与宗教建筑的庄严厚重。南岳大庙是衡山最具代表性的建筑群，采用中轴线对称布局，占地达9.85万平方米，是中国南方最大的古代宫殿式庙宇。庙内的木雕、石雕、泥塑工艺精湛，尤其是正殿的"盘龙石柱"与"铁瓦"，堪称中国古代建筑艺术的杰出代表。此外，祝融峰上的祝融殿、藏经殿等建筑，也以其独特的选址与设计，展现了天人合一的建筑理念。'
          ],
          image: '/src/assets/Imgs/衡山2.jpeg'
        },
        {
          title: '文化价值',
          texts: [
            '衡山是湖湘文化的重要源头之一，孕育了王夫之、魏源等思想家。其宗教文化、祭祀文化与山水文化相互交融，形成了独特的南岳文化体系。道教的"道法自然"与佛教的"慈悲为怀"在此共生，深刻影响了湖湘士人的精神世界，成为中华文化多元一体的生动见证。'
          ],
          image: '/src/assets/Imgs/衡山3.jpeg'
        }
      ],
      '8': [
        {
          title: '历史沿革',
          texts: [
            '花鼓戏起源于清代中晚期的湖南民间歌舞，最初是农民在节庆、婚丧活动中表演的"地花鼓"。后来，它吸收了湘剧、祁剧等剧种的唱腔与表演程式，逐步发展成为成熟的戏曲形式。20世纪以来，经过专业剧团的整理改编，花鼓戏形成了长沙花鼓戏、岳阳花鼓戏、衡州花鼓戏等多个流派，其中长沙花鼓戏影响最广。'
          ],
          image: '/src/assets/Imgs/花鼓戏1.jpeg'
        },
        {
          title: '艺术特色',
          texts: [
            '花鼓戏采用湖南方言演唱，曲调活泼多样，包括"川调""打锣腔""牌子曲"等三大声腔系统。表演贴近生活，以"小丑""小旦""小生"为主要行当，擅长通过夸张的动作与幽默的对白表现民间生活故事。例如《刘海砍樵》中"刘海哥，你是我的夫咯"的经典唱段，至今仍广为流传。'
          ],
          image: '/src/assets/Imgs/花鼓戏2.jpeg'
        },
        {
          title: '传承发展',
          texts: [
            '为保护这一非物质文化遗产，湖南通过建立花鼓戏传承基地、培养非遗传承人、开展校园戏曲普及活动等方式，推动其活态传承。同时，一些剧团结合现代舞台技术进行创新编排，推出了《老表轶事》《桃花烟雨》等新剧目，吸引年轻观众走进传统戏曲。此外，花鼓戏还通过网络直播、短视频等形式传播，进一步扩大了其影响力。'
          ],
          image: '/src/assets/Imgs/花鼓戏3.jpeg'
        }
      ],
      '9': [
        {
          title: '历史沿革',
          texts: [
            '富厚堂始建于清代同治四年（1865年），是曾国藩被封为一等毅勇侯后，由其弟曾国潢主持修建的。建筑历时十年完工，占地达4万余平方米，是曾国藩晚年归隐之地。历经百年风雨，主体建筑保存完好，现为全国重点文物保护单位与国家4A级旅游景区。'
          ],
          image: '/src/assets/Imgs/曾国藩故居1.jpeg'
        },
        {
          title: '建筑特色',
          texts: [
            '富厚堂采用江南民居与官式建筑相结合的风格，布局严谨，雕梁画栋。建筑群分为"公记""朴记""芳记"三大部分，其中"公记"是曾国藩的办公与藏书区域，"朴记""芳记"则为家属居住之所。最具特色的是"求阙斋"藏书楼，藏书达30万卷，涵盖经、史、子、集及西方科技书籍，是中国近代私家藏书的典范。此外，故居的木雕、石雕、彩绘工艺精湛，体现了晚清湖湘士大夫的审美情趣。'
          ],
          image: '/src/assets/Imgs/曾国藩故居2.jpeg'
        },
        {
          title: '文化价值',
          texts: [
            '富厚堂不仅是研究晚清历史与建筑艺术的重要实物资料，更承载着曾国藩的家风家训与湖湘文化"经世致用"的思想。曾国藩所倡导的"耕读传家""修身齐家"理念，通过故居的建筑布局与藏书得以体现，对当代家风建设仍具有重要启示意义。'
          ],
          image: '/src/assets/Imgs/曾国藩故居3.jpeg'
        }
      ],
      '10': [
        {
          title: '历史沿革',
          texts: [
            '1972年至1974年，考古工作者对马王堆汉墓进行了发掘，共清理出三座墓葬，出土文物三千余件。其中，一号墓为利苍之妻辛追的墓葬，二号墓为利苍本人的墓葬，三号墓为他们的儿子的墓葬。墓葬保存完好，尤其是一号墓的辛追夫人遗体，历经两千多年仍保持完整，为研究汉代丧葬习俗与防腐技术提供了珍贵样本。'
          ],
          image: '/src/assets/Imgs/马王堆1.jpeg'
        },
        {
          title: '文化价值',
          texts: [
            '马王堆汉墓出土的文物涵盖帛书、帛画、漆器、丝织品、竹木简等，反映了汉初高度发达的文明水平。例如，《五十二病方》是中国现存最早的医学著作，《五星占》是世界上现存最早的天文著作之一，素纱襌衣则展现了汉代高超的丝织工艺。这些文物不仅改写了中国古代科技史与文化史，更成为展示湖湘文明的重要窗口。'
          ],
          image: '/src/assets/Imgs/马王堆2.jpeg'
        },
        {
          title: '传承发展',
          texts: [
            '出土文物主要收藏于湖南博物院，通过"马王堆汉墓陈列"等专题展览向公众开放。近年来，博物院还运用数字技术对文物进行复原与展示，例如通过3D打印技术复制素纱襌衣，让观众更直观地感受汉代文明的魅力。此外，马王堆汉墓的研究成果还被应用于动画、纪录片等文创产品中，进一步扩大了其文化影响力。'
          ],
          image: '/src/assets/Imgs/马王堆3.jpeg'
        }
      ],
      '11': [
        {
          title: '历史沿革',
          texts: [
            '火宫殿始建于1747年，最初是长沙百姓祭祀火神的场所，每逢农历六月二十三日火神诞辰，都会举行盛大的祭祀活动。民国时期，火宫殿逐渐发展为集餐饮、娱乐于一体的市井场所，成为长沙"老口子"聚集的地方。新中国成立后，火宫殿经过多次扩建与改造，保留了传统的庙宇建筑风貌，同时增加了小吃城、戏台等功能区域。'
          ],
          image: '/src/assets/Imgs/火宫殿1.jpeg'
        },
        {
          title: '文化价值',
          texts: [
            '火宫殿承载着长沙的市井文化与饮食记忆，臭豆腐、糖油粑粑、姊妹团子等特色小吃享誉全国。其中，火宫殿臭豆腐以其外焦里嫩、香辣可口的特点，成为长沙美食的标志性符号。此外，火宫殿还是湖南花鼓戏、湘剧等民间艺术的重要展演场所，戏台每日上演传统戏曲，让游客在品尝美食的同时感受湖湘文化的魅力。'
          ],
          image: 'https://picsum.photos/seed/huogongdian3/800/400'
        },
        {
          title: '传承发展',
          texts: [
            '在保留传统风貌的基础上，火宫殿通过品牌化运营与连锁发展，已在长沙、北京、上海等地开设多家分店。同时，它还积极开展非遗传承与创新活动，例如举办"火宫殿臭豆腐制作技艺"非遗传承培训班，开发火宫殿主题文创产品，让这一老字号品牌在当代焕发新生。'
          ],
          image: 'https://picsum.photos/seed/huogongdian4/800/400'
        }
      ],
      '12': [
        {
          title: '历史沿革',
          texts: [
            '铜官窑始于初唐，盛于中晚唐，衰于五代，烧造历史长达两百多年。1956年，铜官窑遗址被考古发现，出土瓷器数万件，涵盖碗、盘、壶、罐等各类器具。考古研究表明，铜官窑的产品通过湘江进入长江，再经海上丝绸之路远销海外，在印度尼西亚、伊朗、埃及等国均有出土。'
          ],
          image: 'https://picsum.photos/seed/tongguanyao2/800/400'
        },
        {
          title: '艺术特色',
          texts: [
            '铜官窑首创釉下多彩工艺，以氧化铜、氧化铁等为颜料，在瓷胎上绘制纹饰后再罩一层透明釉，经高温烧制而成。纹饰题材丰富，包括人物、动物、花卉、诗文等，其中"诗文瓷器"最具特色，例如"君生我未生，我生君已老"等诗句被直接绘制在瓷器上，反映了唐代的社会风貌与审美情趣。此外，铜官窑的瓷器造型多样，既有实用器具，也有艺术摆件，体现了唐代手工业的高度成就。'
          ],
          image: 'https://picsum.photos/seed/tongguanyao3/800/400'
        },
        {
          title: '文化价值',
          texts: [
            '铜官窑是唐代海上丝绸之路的重要物证，它不仅开创了中国陶瓷彩绘的新纪元，更通过瓷器的外销促进了中外文化交流。其釉下彩工艺对后世的青花、五彩等瓷器产生了深远影响，是中国陶瓷史上的重要里程碑。'
          ],
          image: 'https://picsum.photos/seed/tongguanyao4/800/400'
        }
      ],
      '13': [
        {
          title: '生平简介',
          texts: [
            '曾国藩（1811年11月26日－1872年3月12日），字伯涵，号涤生，湖南长沙府湘乡县人。',
            '道光十八年（1838年）中进士，入翰林院，为军机大臣穆彰阿门生。累迁内阁学士，礼部侍郎，署兵、工、刑、吏部侍郎。',
            '太平天国运动时，曾国藩组建湘军，力挽狂澜，经过多年鏖战后攻灭太平天国。',
            '官至两江总督、直隶总督、武英殿大学士，封一等毅勇侯，谥号"文正"，后世称"曾文正"。'
          ],
          image: 'https://picsum.photos/seed/zengguofanbio/800/400'
        },
        {
          title: '历史功绩',
          texts: [
            '曾国藩是湘军的创立者和统帅，他率领湘军镇压了太平天国运动，维护了清朝的统治。',
            '他倡导洋务运动，主张学习西方先进技术，创办了安庆内军械所等近代工业，是中国近代化建设的开拓者之一。',
            '曾国藩在文化教育方面也有重要贡献，他创办了时务学堂，推动了湖南地区的教育发展。',
            '他的《曾国藩家书》对后世影响深远，被誉为"处世之良法，齐家之要道"。'
          ],
          image: 'https://picsum.photos/seed/zengguofanachieve/800/400'
        },
        {
          title: '思想与影响',
          texts: [
            '曾国藩的思想以儒家思想为基础，融合了程朱理学和经世致用的思想，形成了自己独特的思想体系。',
            '他强调"修身齐家治国平天下"，注重个人道德修养和家庭伦理建设。',
            '曾国藩的用人之道和治军理念对后世产生了深远影响，他提出的"忠义血性"、"扎硬寨，打呆仗"等观点至今仍有借鉴意义。',
            '他与胡林翼并称"曾胡"，与李鸿章、左宗棠、张之洞并称"晚清中兴四大名臣"。'
          ]
        }
      ],
    }
    
    // 当前资源的详细内容段落
    const detailParagraphs = ref([])
    
    // 计算当前图片
    const currentImage = computed(() => {
      if (resource.value.images && resource.value.images.length > 0) {
        return resource.value.images[currentImageIndex.value]
      }
      return ''
    })
    
    // 设置当前图片
    const setCurrentImage = (index) => {
      currentImageIndex.value = index
    }
    
    // 切换点赞资源
    const toggleLikeResource = () => {
      if (!props.user.isLoggedIn) {
        if (props.showLoginModal) {
          props.showLoginModal()
        }
        return
      }
      
      resource.value.isLiked = !resource.value.isLiked
      // 移除点赞计数逻辑
      
      if (props.showAlert) {
        props.showAlert('success', resource.value.isLiked ? '点赞成功！' : '已取消点赞')
      }
    }
    
    // 分享资源
    const shareResource = () => {
      // 模拟分享功能
      if (navigator.share) {
        navigator.share({
          title: resource.value.title,
          text: resource.value.introduction,
          url: window.location.href
        }).then(() => {
          if (props.showAlert) {
            props.showAlert('success', '分享成功！')
          }
        }).catch(() => {
          if (props.showAlert) {
            props.showAlert('error', '分享失败，请稍后再试。')
          }
        })
      } else {
        // 复制链接到剪贴板
        navigator.clipboard.writeText(window.location.href).then(() => {
          if (props.showAlert) {
            props.showAlert('success', '链接已复制到剪贴板！')
          }
        }).catch(() => {
          if (props.showAlert) {
            props.showAlert('error', '复制失败，请手动复制链接。')
          }
        })
      }
    }
    
    // 举报资源
    const reportResource = () => {
      // 模拟举报功能
      if (props.showAlert) {
        props.showAlert('success', '举报信息已提交，我们会尽快处理。')
      }
    }
    
    // 关注贡献者
    const followContributor = () => {
      if (!props.user.isLoggedIn) {
        if (props.showLoginModal) {
          props.showLoginModal()
        }
        return
      }
      
      resource.value.isFollowing = true
      
      if (props.showAlert) {
        props.showAlert('success', '关注成功！')
      }
    }
    
    // 取消关注贡献者
    const unfollowContributor = () => {
      resource.value.isFollowing = false
      
      if (props.showAlert) {
        props.showAlert('success', '已取消关注。')
      }
    }
    
    // 按标签筛选
    const filterByTag = (tag) => {
      // 模拟按标签筛选
      if (props.goToResources) {
        props.goToResources(tag)
      }
    }
    
    // 下载资源
    const downloadResource = (download) => {
      // 模拟下载功能
      if (props.showAlert) {
        props.showAlert('success', `开始下载：${download.name}`)
      }
    }
    
    // 提交评论
    const submitComment = async () => {
      if (!newComment.value.trim()) return
      
      isSubmittingComment.value = true
      
      try {
        // 模拟提交评论
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        const newCommentObj = {
          id: Date.now().toString(),
          userName: props.user.username || '游客',
          userAvatar: props.user.avatar || 'https://picsum.photos/seed/default/100/100',
          content: newComment.value,
          time: new Date().toLocaleString('zh-CN'),
          likes: 0,
          isLiked: false
        }
        
        comments.value.unshift(newCommentObj)
        newComment.value = ''
        
        if (props.showAlert) {
          props.showAlert('success', '评论发布成功！')
        }
      } catch (error) {
        console.error('提交评论失败:', error)
        if (props.showAlert) {
          props.showAlert('error', '评论发布失败，请稍后再试。')
        }
      } finally {
        isSubmittingComment.value = false
      }
    }
    
    // 切换点赞评论
    const toggleLikeComment = (commentId) => {
      if (!props.user.isLoggedIn) {
        if (props.showLoginModal) {
          props.showLoginModal()
        }
        return
      }
      
      const comment = comments.value.find(c => c.id === commentId)
      if (comment) {
        comment.isLiked = !comment.isLiked
        comment.likes += comment.isLiked ? 1 : -1
      }
    }
    
    // 切换回复表单
    const toggleReplyForm = (commentId) => {
      if (!props.user.isLoggedIn) {
        if (props.showLoginModal) {
          props.showLoginModal()
        }
        return
      }
      
      if (replyToId.value === commentId) {
        replyToId.value = null
        replyContent.value = ''
      } else {
        replyToId.value = commentId
        replyContent.value = ''
      }
    }
    
    // 取消回复
    const cancelReply = () => {
      replyToId.value = null
      replyContent.value = ''
    }
    
    // 提交回复
    const submitReply = async (commentId) => {
      if (!replyContent.value.trim()) return
      
      isSubmittingReply.value = true
      
      try {
        // 模拟提交回复
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        const comment = comments.value.find(c => c.id === commentId)
        if (comment) {
          // 这里简化处理，实际项目中应该有回复的嵌套结构
          const newReplyComment = {
            id: Date.now().toString(),
            userName: props.user.username || '游客',
            userAvatar: props.user.avatar || 'https://picsum.photos/seed/default/100/100',
            content: `回复 @${comment.userName}：${replyContent.value}`,
            time: new Date().toLocaleString('zh-CN'),
            likes: 0,
            isLiked: false
          }
          
          comments.value.splice(comments.value.indexOf(comment) + 1, 0, newReplyComment)
          cancelReply()
          
          if (props.showAlert) {
            props.showAlert('success', '回复发布成功！')
          }
        }
      } catch (error) {
        console.error('提交回复失败:', error)
        if (props.showAlert) {
          props.showAlert('error', '回复发布失败，请稍后再试。')
        }
      } finally {
        isSubmittingReply.value = false
      }
    }
    
    // 返回上一页
    const goBack = () => {
      router.back()
    }
    
    // 跳转到资源列表页
    const goToResources = (tag = null) => {
      if (tag) {
        router.push({ name: 'CulturalResourcesPage', query: { tag } })
      } else {
        router.push({ name: 'CulturalResourcesPage' })
      }
    }
    
    // 跳转到资源详情页
    const goToResourceDetail = (id) => {
      router.push({ name: 'ResourceDetailPage', params: { id } })
    }
    
    // 组件挂载时
    onMounted(() => {
      // 根据resourceId获取资源数据
      console.log(`加载资源ID: ${resourceId} 的详细信息`)
      
      // 从资源数据仓库中查找对应ID的资源
      if (resourcesData[resourceId]) {
        // 深拷贝资源数据，避免直接修改数据源
        resource.value = { ...resourcesData[resourceId] }
        
        // 设置对应的详细内容段落
        if (detailParagraphsMap[resourceId]) {
          detailParagraphs.value = [...detailParagraphsMap[resourceId]]
        } else {
          detailParagraphs.value = []
        }
        
        // 模拟增加访问量
        resource.value.views++
      } else {
        console.warn(`未找到资源ID: ${resourceId} 的数据`)
        // 使用默认资源数据
        resource.value = { ...defaultResource }
        detailParagraphs.value = []
      }
      
      // 滚动到页面顶部
      window.scrollTo(0, 0)
    })
    
    return {
      resource,
      currentImageIndex,
      currentImage,
      comments,
      newComment,
      isSubmittingComment,
      replyToId,
      replyContent,
      isSubmittingReply,
      detailParagraphs,
      setCurrentImage,
      toggleLikeResource,
      shareResource,
      reportResource,
      followContributor,
      unfollowContributor,
      filterByTag,
      downloadResource,
      submitComment,
      toggleLikeComment,
      toggleReplyForm,
      cancelReply,
      submitReply,
      goBack,
      goToResources,
      goToResourceDetail
    }
  }
}
</script>

<style scoped>
.resource-detail-page {
  padding: 2rem 0;
}

/* 返回按钮样式 */
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
  0%, 100% { transform: translateX(0) scale(1.05); }
  50% { transform: translateX(5px) scale(1.05); }
}

.resource-header {
  background-color: #f8f9fa;
  padding: 1rem 0;
  border-bottom: 1px solid #eee;
}

.breadcrumb {
  font-size: 0.9rem;
  color: #666;
}

.breadcrumb a {
  color: var(--primary-color);
  text-decoration: none;
}

.breadcrumb a:hover {
  text-decoration: underline;
}

.breadcrumb span {
  color: #666;
}

.resource-content {
  display: flex;
  flex-wrap: wrap;
  gap: 2rem;
  margin: 2rem 0;
}

.resource-main {
  flex: 3;
  min-width: 300px;
}

.resource-sidebar {
  flex: 1;
  min-width: 250px;
}

/* 资源标题和元信息 */
.resource-title-section {
  margin-bottom: 2rem;
}

.resource-title-section h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 1rem;
}

.resource-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
  font-size: 0.9rem;
  color: #666;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* 资源图片轮播 */
.resource-gallery {
  margin-bottom: 2rem;
}

.gallery-main {
  width: 100%;
  height: 400px;
  overflow: hidden;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.main-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gallery-thumbs {
  display: flex;
  gap: 1rem;
  overflow-x: auto;
  padding: 0.5rem 0;
}

.thumb-item {
  width: 100px;
  height: 60px;
  overflow: hidden;
  border-radius: 4px;
  cursor: pointer;
  border: 2px solid transparent;
  flex-shrink: 0;
}

.thumb-item.active {
  border-color: var(--primary-color);
}

.thumb-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.thumb-item:hover img {
  transform: scale(1.1);
}

/* 资源简介和详细内容 */
.resource-intro,
.resource-detail-content {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.resource-intro h2,
.resource-detail-content h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--primary-color);
}

.resource-intro p {
  color: #666;
  line-height: 1.6;
}

.content-paragraph {
  margin-bottom: 2rem;
}

.content-paragraph:last-child {
  margin-bottom: 0;
}

.content-paragraph h3 {
  font-size: 1.2rem;
  color: #333;
  margin: 1rem 0;
}

.content-paragraph p {
  color: #666;
  line-height: 1.6;
  margin-bottom: 1rem;
}

.content-image {
  width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 1rem 0;
}

/* 相关资源 */
.related-resources {
  margin-bottom: 2rem;
}

.related-resources h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.related-item {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
}

.related-item:hover {
  transform: translateY(-5px);
}

.related-image {
  width: 100px;
  height: 100px;
  border-radius: 4px;
  overflow: hidden;
  flex-shrink: 0;
}

.related-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.related-info {
  flex: 1;
}

.related-info h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
}

.related-info h3 a {
  color: #333;
  text-decoration: none;
}

.related-info h3 a:hover {
  color: var(--primary-color);
}

.related-info p {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.4;
  margin-bottom: 0.5rem;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
}

.related-category {
  display: inline-block;
  padding: 0.2rem 0.5rem;
  background-color: #f0f0f0;
  color: #666;
  font-size: 0.8rem;
  border-radius: 10px;
}

/* 评论区 */
.comments-section {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
}

.comments-section h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1.5rem;
}

.comment-form textarea {
  width: 100%;
  padding: 1rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: vertical;
  font-size: 1rem;
}

.comment-form textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-actions {
  text-align: right;
  margin-top: 1rem;
}

.submit-comment {
  padding: 0.5rem 1.5rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

.submit-comment:hover:not(:disabled) {
  background-color: #c62828;
}

.submit-comment:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.login-prompt {
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 4px;
  text-align: center;
}

.login-prompt a {
  color: var(--primary-color);
  text-decoration: none;
}

.login-prompt a:hover {
  text-decoration: underline;
}

.comments-list {
  margin-top: 1.5rem;
}

.comment-item {
  display: flex;
  gap: 1rem;
  padding: 1rem 0;
  border-bottom: 1px solid #eee;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.comment-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.comment-header h4 {
  margin: 0;
  color: #333;
  font-size: 1rem;
}

.comment-header span {
  color: #999;
  font-size: 0.8rem;
}

.comment-text p {
  margin: 0 0 1rem 0;
  color: #666;
  line-height: 1.5;
}

.comment-actions {
  display: flex;
  gap: 1.5rem;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 0.9rem;
}

.action-btn:hover {
  color: var(--primary-color);
}

.reply-form {
  margin-top: 1rem;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.reply-form textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  resize: vertical;
  font-size: 0.9rem;
}

.reply-form textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

.reply-actions {
  text-align: right;
  margin-top: 0.5rem;
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.cancel-reply {
  padding: 0.25rem 1rem;
  background-color: #ddd;
  color: #333;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.submit-reply {
  padding: 0.25rem 1rem;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.submit-reply:hover:not(:disabled) {
  background-color: #c62828;
}

.submit-reply:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.no-comments {
  text-align: center;
  padding: 2rem;
  color: #999;
}

/* 侧边栏样式 */
.info-card,
.contributor-card,
.tags-card,
.download-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-card h3,
.contributor-card h3,
.tags-card h3,
.download-card h3 {
  font-size: 1.2rem;
  color: #333;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--primary-color);
}

.info-item {
  display: flex;
  margin-bottom: 0.75rem;
}

.info-label {
  width: 80px;
  color: #666;
  font-weight: 500;
}

.info-value {
  flex: 1;
  color: #333;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.action-buttons .action-btn {
  justify-content: center;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  font-size: 1rem;
  transition: all 0.3s;
}

.action-buttons .action-btn:hover {
  background-color: #f8f9fa;
  border-color: var(--primary-color);
}

.action-buttons .action-btn.primary {
  background-color: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.action-buttons .action-btn.primary:hover {
  background-color: #c62828;
}

/* 贡献者卡片 */
.contributor-info {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-bottom: 1rem;
}

.contributor-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
}

.contributor-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.contributor-details h4 {
  margin: 0 0 0.25rem 0;
  color: #333;
}

.contributor-details p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.follow-btn,
.following-btn {
  width: 100%;
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9rem;
}

.follow-btn {
  background-color: var(--primary-color);
  color: white;
}

.follow-btn:hover {
  background-color: #c62828;
}

.following-btn {
  background-color: #f8f9fa;
  color: #666;
  border: 1px solid #ddd;
}

.following-btn:hover {
  background-color: #e9ecef;
}

/* 标签卡片 */
.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag-item {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background-color: #f8f9fa;
  color: #666;
  text-decoration: none;
  border-radius: 15px;
  font-size: 0.9rem;
  transition: all 0.3s;
}

.tag-item:hover {
  background-color: var(--primary-color);
  color: white;
}

/* 下载卡片 */
.download-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.download-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background-color: #f8f9fa;
  color: #333;
  text-decoration: none;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.download-item:hover {
  background-color: #e9ecef;
}

.download-item i {
  color: var(--primary-color);
}

.download-size {
  margin-left: auto;
  color: #666;
  font-size: 0.9rem;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .resource-content {
    flex-direction: column;
  }
  
  .gallery-main {
    height: 250px;
  }
  
  .related-grid {
    grid-template-columns: 1fr;
  }
  
  .resource-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .comment-item {
    flex-direction: column;
  }
  
  .comment-avatar {
    width: 50px;
    height: 50px;
  }
}
</style>