<template>
  <div class="resource-detail-page">
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
          'https://picsum.photos/seed/yuelu1/1200/600',
          'https://picsum.photos/seed/yuelu2/1200/600',
          'https://picsum.photos/seed/yuelu3/1200/600',
          'https://picsum.photos/seed/yuelu4/1200/600'
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
          'https://picsum.photos/seed/xiangxiu1/1200/600',
          'https://picsum.photos/seed/xiangxiu2/1200/600',
          'https://picsum.photos/seed/xiangxiu3/1200/600'
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
          'https://picsum.photos/seed/yueyanglou1/1200/600',
          'https://picsum.photos/seed/yueyanglou2/1200/600',
          'https://picsum.photos/seed/yueyanglou3/1200/600'
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
      '13': {
        id: '13',
        title: '曾国藩',
        category: '历史伟人',
        introduction: '曾国藩（1811年11月26日－1872年3月12日），字伯涵，号涤生，湖南长沙府湘乡县人，中国清朝后期政治家、战略家、理学家、文学家，湘军的创立者和统帅。',
        images: [
          'https://picsum.photos/seed/zengguofan2/1200/600',
          'https://picsum.photos/seed/zengguofan3/1200/600',
          'https://picsum.photos/seed/zengguofan4/1200/600'
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
          image: 'https://picsum.photos/seed/history/800/400'
        },
        {
          title: '建筑特色',
          texts: [
            '岳麓书院古建筑群分为教学、藏书、祭祀、园林、纪念五大建筑格局。岳麓书院主体建筑面积有31000多平方米，分为书院主体、附属文庙及新建的中国书院博物馆。',
            '岳麓书院占地面积25000平方米，现存建筑大部分为明清遗物，其古建筑在布局上采用中轴对称、纵深多进的院落形式。主体建筑如头门、大门、二门、讲堂、御书楼集中于中轴线上，讲堂布置在中轴线的中央。斋舍、祭祀专祠等排列于两旁。',
            '中轴对称、层层递进的院落，除了营造一种庄严、神妙、幽远的纵深感和视觉效应之外，还体现了儒家文化尊卑有序、等级有别、主次鲜明的社会伦理关系。'
          ],
          image: 'https://picsum.photos/seed/architecture/800/400'
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
          image: 'https://picsum.photos/seed/scholars/800/400'
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
          image: 'https://picsum.photos/seed/embroidery1/800/400'
        },
        {
          title: '艺术特色',
          texts: [
            '湘绣以其精湛的针法、丰富的色彩和逼真的表现力著称。常用的针法有掺针、毛针、游针、齐针等70多种。',
            '湘绣的色彩运用大胆而细腻，讲究"色相分明，色阶渐变"，通过不同颜色的丝线相互搭配，形成丰富的层次感和立体感。',
            '湘绣的题材广泛，包括花鸟、山水、人物、走兽等，其中以狮、虎等动物题材最为著名，有"湘绣狮虎冠天下"的美誉。'
          ],
          image: 'https://picsum.photos/seed/embroidery2/800/400'
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
          image: 'https://picsum.photos/seed/yueyanghistory/800/400'
        },
        {
          title: '建筑艺术',
          texts: [
            '岳阳楼建筑风格独特，主楼高19.42米，为三层、四柱、飞檐、盔顶、纯木结构。',
            '其盔顶造型在中国古代建筑中独一无二，充分体现了古代劳动人民的智慧和创造力。',
            '楼内藏有大量历代名人的诗词、书画和碑刻，其中最著名的当属范仲淹的《岳阳楼记》木雕屏。'
          ],
          image: 'https://picsum.photos/seed/yueyangarchitecture/800/400'
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
      ]
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