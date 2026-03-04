<template>
  <div class="poetry-digitalization-page">
    <!-- 返回按钮 -->
    <div class="back-button-container">
      <button class="back-button" @click="goBack">
        <i class="fas fa-arrow-left"></i>
        <span>返回数字化展示</span>
      </button>
    </div>

    <!-- 主标题区域 -->
    <div class="section-title">
      <h2>湖湘诗词数字化</h2>
      <p>数字化呈现湖湘地区经典诗词作品，配合注释和朗诵音频</p>
    </div>

    <!-- 内容区域 -->
    <div class="content-container">
      <!-- 专题介绍 -->
      <div class="introduction-section">
        <h3>专题介绍</h3>
        <p>湖湘诗词数字化项目旨在通过现代信息技术手段，将湖湘地区的经典诗词作品进行数字化处理和呈现。我们收集了从古代到近现代的湖湘地区代表性诗词作品，为每首诗词提供详细的注释、作者介绍、创作背景等信息，并邀请专业朗诵者录制朗诵音频，让用户能够通过多种感官体验湖湘诗词的魅力。</p>
        <p>该项目不仅是对湖湘文化遗产的保护和传承，也是对传统诗词教育方式的创新和尝试。通过数字化平台，用户可以更加便捷地学习和欣赏湖湘诗词，深入了解湖湘文化的深厚底蕴。</p>
      </div>

      <!-- 诗词展示 -->
      <div class="poetry-showcase">
        <h3>经典诗词展示</h3>
        <div class="poetry-grid">
          <!-- 诗词卡片 1 -->
          <div class="poetry-card">
            <div class="poetry-header">
              <h4>《岳阳楼记》</h4>
              <p class="poetry-author">范仲淹</p>
            </div>
            <div class="poetry-content">
              <p>庆历四年春，滕子京谪守巴陵郡。越明年，政通人和，百废具兴。乃重修岳阳楼，增其旧制，刻唐贤今人诗赋于其上。属予作文以记之。</p>
              <p>予观夫巴陵胜状，在洞庭一湖。衔远山，吞长江，浩浩汤汤，横无际涯；朝晖夕阴，气象万千。此则岳阳楼之大观也，前人之述备矣。然则北通巫峡，南极潇湘，迁客骚人，多会于此，览物之情，得无异乎？</p>
            </div>
            <div class="poetry-actions">
              <button class="action-btn" @click="playAudio('yueyanglou')">
                <i class="fas fa-play"></i>
                <span>播放朗诵</span>
              </button>
              <button class="action-btn" @click="toggleAnnotation('yueyanglou')">
                <i class="fas fa-book"></i>
                <span>查看注释</span>
              </button>
            </div>
            <!-- 注释对话框 -->
            <div class="annotation-dialog" v-if="annotationVisible.yueyanglou">
              <div class="annotation-header">
                <h5>{{ annotations.yueyanglou.title }}</h5>
                <button class="close-btn" @click="toggleAnnotation('yueyanglou')">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <div class="annotation-body">
                <div class="annotation-section">
                  <h6>作者简介</h6>
                  <p>{{ annotations.yueyanglou.author }}</p>
                </div>
                <div class="annotation-section">
                  <h6>创作背景</h6>
                  <p>{{ annotations.yueyanglou.background }}</p>
                </div>
                <div class="annotation-section">
                  <h6>词句注释</h6>
                  <div class="annotation-notes">
                    <div class="note-item" v-for="(note, index) in annotations.yueyanglou.notes" :key="index">
                      <p class="note-line">{{ note.line }}</p>
                      <p class="note-content">{{ note.note }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 诗词卡片 2 -->
          <div class="poetry-card">
            <div class="poetry-header">
              <h4>《望洞庭湖赠张丞相》</h4>
              <p class="poetry-author">孟浩然</p>
            </div>
            <div class="poetry-content">
              <p>八月湖水平，涵虚混太清。</p>
              <p>气蒸云梦泽，波撼岳阳城。</p>
              <p>欲济无舟楫，端居耻圣明。</p>
              <p>坐观垂钓者，徒有羡鱼情。</p>
            </div>
            <div class="poetry-actions">
              <button class="action-btn" @click="playAudio('dongtinghu')">
                <i class="fas fa-play"></i>
                <span>播放朗诵</span>
              </button>
              <button class="action-btn" @click="toggleAnnotation('dongtinghu')">
                <i class="fas fa-book"></i>
                <span>查看注释</span>
              </button>
            </div>
            <!-- 注释对话框 -->
            <div class="annotation-dialog" v-if="annotationVisible.dongtinghu">
              <div class="annotation-header">
                <h5>{{ annotations.dongtinghu.title }}</h5>
                <button class="close-btn" @click="toggleAnnotation('dongtinghu')">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <div class="annotation-body">
                <div class="annotation-section">
                  <h6>作者简介</h6>
                  <p>{{ annotations.dongtinghu.author }}</p>
                </div>
                <div class="annotation-section">
                  <h6>创作背景</h6>
                  <p>{{ annotations.dongtinghu.background }}</p>
                </div>
                <div class="annotation-section">
                  <h6>词句注释</h6>
                  <div class="annotation-notes">
                    <div class="note-item" v-for="(note, index) in annotations.dongtinghu.notes" :key="index">
                      <p class="note-line">{{ note.line }}</p>
                      <p class="note-content">{{ note.note }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 诗词卡片 3 -->
          <div class="poetry-card">
            <div class="poetry-header">
              <h4>《长沙过贾谊宅》</h4>
              <p class="poetry-author">刘长卿</p>
            </div>
            <div class="poetry-content">
              <p>三年谪宦此栖迟，万古惟留楚客悲。</p>
              <p>秋草独寻人去后，寒林空见日斜时。</p>
              <p>汉文有道恩犹薄，湘水无情吊岂知？</p>
              <p>寂寂江山摇落处，怜君何事到天涯！</p>
            </div>
            <div class="poetry-actions">
              <button class="action-btn" @click="playAudio('jiayizhai')">
                <i class="fas fa-play"></i>
                <span>播放朗诵</span>
              </button>
              <button class="action-btn" @click="toggleAnnotation('jiayizhai')">
                <i class="fas fa-book"></i>
                <span>查看注释</span>
              </button>
            </div>
            <!-- 注释对话框 -->
            <div class="annotation-dialog" v-if="annotationVisible.jiayizhai">
              <div class="annotation-header">
                <h5>{{ annotations.jiayizhai.title }}</h5>
                <button class="close-btn" @click="toggleAnnotation('jiayizhai')">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <div class="annotation-body">
                <div class="annotation-section">
                  <h6>作者简介</h6>
                  <p>{{ annotations.jiayizhai.author }}</p>
                </div>
                <div class="annotation-section">
                  <h6>创作背景</h6>
                  <p>{{ annotations.jiayizhai.background }}</p>
                </div>
                <div class="annotation-section">
                  <h6>词句注释</h6>
                  <div class="annotation-notes">
                    <div class="note-item" v-for="(note, index) in annotations.jiayizhai.notes" :key="index">
                      <p class="note-line">{{ note.line }}</p>
                      <p class="note-content">{{ note.note }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 互动功能 -->
      <div class="interactive-section">
        <h3>互动功能</h3>
        <div class="interactive-features">
          <div class="feature-card">
            <div class="feature-icon">
              <i class="fas fa-microphone fa-3x"></i>
            </div>
            <h4>语音朗诵</h4>
            <p>点击播放按钮，聆听专业朗诵者对诗词的精彩演绎，感受诗词的韵律美。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="fas fa-book-open fa-3x"></i>
            </div>
            <h4>详细注释</h4>
            <p>查看每首诗词的详细注释，了解诗词的背景、典故和深刻内涵。</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <i class="fas fa-search fa-3x"></i>
            </div>
            <h4>诗词检索</h4>
            <p>通过作者、朝代、关键词等多种方式检索诗词，快速找到您感兴趣的作品。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, ref } from 'vue'
import { useRouter } from 'vue-router'

// 获取路由实例
const router = useRouter()
// 获取全局事件总线
const eventBus = inject('$eventBus')

// 注释对话框状态
const annotationVisible = ref({
  yueyanglou: false,
  dongtinghu: false,
  jiayizhai: false
})

// 详细注释内容
const annotations = {
  yueyanglou: {
    title: '《岳阳楼记》注释',
    author: '范仲淹（989-1052），字希文，北宋著名的政治家、思想家、军事家、文学家，世称"范文正公"。',
    background: '本文作于庆历六年（1046年），范仲淹因倡导改革被贬邓州。滕子京重修岳阳楼，请范仲淹作记。',
    notes: [
      { line: '庆历四年春，滕子京谪守巴陵郡。', note: '庆历四年：公元1044年。滕子京：名宗谅，字子京，范仲淹的朋友。谪：贬官。巴陵郡：今湖南岳阳。' },
      { line: '越明年，政通人和，百废具兴。', note: '越：到了。政通人和：政事顺利，百姓和乐。具：通"俱"，全、都。' },
      { line: '乃重修岳阳楼，增其旧制，刻唐贤今人诗赋于其上。', note: '乃：于是。增其旧制：扩大它原有的规模。唐贤：唐代名家。' },
      { line: '属予作文以记之。', note: '属：通"嘱"，嘱托。予：我。作文：写文章。' },
      { line: '予观夫巴陵胜状，在洞庭一湖。', note: '夫：那。胜状：胜景，好景色。' },
      { line: '衔远山，吞长江，浩浩汤汤，横无际涯；朝晖夕阴，气象万千。', note: '衔：衔接。吞：吞没。浩浩汤汤：水势浩大的样子。横无际涯：宽阔无边。朝晖夕阴：或早或晚（一天里）阴晴多变化。气象万千：气象多种多样。' },
      { line: '此则岳阳楼之大观也，前人之述备矣。', note: '大观：雄伟景象。备：详尽。' },
      { line: '然则北通巫峡，南极潇湘，迁客骚人，多会于此，览物之情，得无异乎？', note: '然则：虽然如此，那么。北通巫峡：北面通向巫峡。南极潇湘：南面直到潇水、湘水。迁客：谪迁的人，指降职远调的人。骚人：诗人。会：聚集。览物之情：观赏自然景物而触发的感情。得无异乎：大概会有所不同吧？' }
    ]
  },
  dongtinghu: {
    title: '《望洞庭湖赠张丞相》注释',
    author: '孟浩然（689-740），字浩然，号孟山人，唐代著名的山水田园派诗人，世称"孟襄阳"。',
    background: '此诗是孟浩然投赠张九龄的干谒诗，希望得到张九龄的引荐。',
    notes: [
      { line: '八月湖水平，涵虚混太清。', note: '湖水平：八月长江水涨，洞庭湖水位升高，与岸齐平。涵虚：包含天空，指天空倒映在水中。混太清：与天混同。太清：天空。' },
      { line: '气蒸云梦泽，波撼岳阳城。', note: '气蒸：水气蒸腾。云梦泽：古代大泽，在今湖北南部、湖南北部一带，后淤积为陆地。波撼：波浪摇动。岳阳城：今湖南岳阳，在洞庭湖东岸。' },
      { line: '欲济无舟楫，端居耻圣明。', note: '济：渡。舟楫：船只。端居：闲居。耻圣明：有愧于圣明之世。圣明：指太平盛世，古时认为皇帝圣明，社会就会安定。' },
      { line: '坐观垂钓者，徒有羡鱼情。', note: '垂钓者：比喻有官位的人。羡鱼情：羡慕钓鱼人的心情，比喻想做官。' }
    ]
  },
  jiayizhai: {
    title: '《长沙过贾谊宅》注释',
    author: '刘长卿（约726-约786），字文房，唐代诗人，擅长五言诗，自称"五言长城"。',
    background: '此诗作于大历年间，刘长卿被贬为睦州司马，路过长沙贾谊故宅时所作。',
    notes: [
      { line: '三年谪宦此栖迟，万古惟留楚客悲。', note: '三年谪宦：贾谊被贬为长沙王太傅三年。栖迟：居留。楚客：指贾谊，也指诗人自己，因为两人都是被贬谪的南方人。' },
      { line: '秋草独寻人去后，寒林空见日斜时。', note: '人去后：指贾谊离开长沙之后。日斜时：指傍晚时分。' },
      { line: '汉文有道恩犹薄，湘水无情吊岂知？', note: '汉文：汉文帝刘恒。有道：指汉文帝是明君。恩犹薄：但对贾谊的恩遇仍然微薄。吊：凭吊。岂知：哪里知道。' },
      { line: '寂寂江山摇落处，怜君何事到天涯！', note: '摇落处：指秋天草木凋零的地方。君：指贾谊，也指诗人自己。天涯：指远离京城的地方，这里指长沙。' }
    ]
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

// 播放音频
const playAudio = (poetryId) => {
  showAlert('音频功能即将上线，敬请期待！')
}

// 显示/隐藏注释
const toggleAnnotation = (poetryId) => {
  annotationVisible.value[poetryId] = !annotationVisible.value[poetryId]
}

// 显示提示消息
const showAlert = (message, type = 'info') => {
  if (eventBus) {
    eventBus.emit('show-alert', { message, type })
  } else {
    console.log(`[${type}] ${message}`)
  }
}
</script>

<style scoped>
:root {
  --primary-color: #C8102E; /* 湘红色 */
  --bg-color: #F9FAFB;
  --text-color: #1F2937;
  --light-text: #6B7280;
  --heading-font: 'Microsoft YaHei', 'STXihei', 'SimHei', sans-serif;
  --body-font: 'Microsoft YaHei', 'STXihei', 'SimSun', serif;
}

.poetry-digitalization-page {
  padding: 2rem 0;
  background-color: #fafafa;
  min-height: 100vh;
}

/* 返回按钮样式 */
.back-button-container {
  position: fixed;
  top: 2rem;
  left: 2rem;
  z-index: 1000;
}

.back-button {
  background: var(--primary-color);
  color: white;
  border: none;
  padding: 0.8rem 1.5rem;
  border-radius: 25px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(200, 16, 46, 0.3);
}

.back-button:hover {
  background: #a60e24;
  transform: translateX(-5px) scale(1.05);
  box-shadow: 0 6px 16px rgba(200, 16, 46, 0.4);
}

.back-button i {
  font-size: 1.1rem;
}

/* 主标题样式 */
.section-title {
  text-align: center;
  margin-bottom: 3rem;
  padding: 0 1rem;
  margin-top: 60px;
}

.section-title h2 {
  font-size: 2.5rem;
  color: var(--text-color);
  margin-bottom: 0.5rem;
  font-family: var(--heading-font);
  font-weight: 700;
  line-height: 1.3;
  position: relative;
  display: inline-block;
}

.section-title h2::after {
  content: '';
  position: absolute;
  bottom: -10px;
  left: 50%;
  transform: translateX(-50%);
  width: 100px;
  height: 3px;
  background-color: var(--primary-color);
}

.section-title p {
  color: var(--light-text);
  max-width: 600px;
  margin: 1.5rem auto 0;
  line-height: 1.7;
  font-size: 16px;
}

/* 内容容器 */
.content-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

/* 介绍部分 */
.introduction-section {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  margin-bottom: 3rem;
}

.introduction-section h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.introduction-section h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.introduction-section p {
  color: var(--text-color);
  line-height: 1.7;
  margin-bottom: 1.2rem;
}

/* 诗词展示 */
.poetry-showcase {
  margin-bottom: 3rem;
}

.poetry-showcase h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.poetry-showcase h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.poetry-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
}

.poetry-card {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
  min-height: 450px;
}

.poetry-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.poetry-header {
  padding: 1.5rem;
  background: var(--primary-color);
  color: white;
}

.poetry-header h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.3rem;
  font-weight: 600;
}

.poetry-author {
  margin: 0;
  font-size: 0.95rem;
  opacity: 0.9;
}

.poetry-content {
  padding: 1.5rem;
  line-height: 1.8;
  color: var(--text-color);
  flex: 1;
  text-align: center;
}

.poetry-content p {
  margin-bottom: 1rem;
  text-align: center;
}

.poetry-content p:last-child {
  margin-bottom: 0;
}

.poetry-actions {
  padding: 1.5rem;
  border-top: 1px solid #e9ecef;
  display: flex;
  gap: 1rem;
  align-self: stretch;
}

.action-btn {
  flex: 1;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  padding: 0.8rem;
  border-radius: 4px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
  transform: translateY(-1px);
}

/* 互动功能 */
.interactive-section {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
}

.interactive-section h3 {
  font-size: 1.8rem;
  color: var(--text-color);
  margin-bottom: 1.5rem;
  font-family: var(--heading-font);
  font-weight: 600;
  position: relative;
  padding-left: 1rem;
}

.interactive-section h3::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 4px;
  background-color: var(--primary-color);
}

.interactive-features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.feature-card {
  text-align: center;
  padding: 2rem;
  border-radius: 8px;
  background: #f8f9fa;
  transition: transform 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-icon {
  color: var(--primary-color);
  margin-bottom: 1.5rem;
  background-color: rgba(200, 16, 46, 0.05);
  width: 100px;
  height: 100px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
}

.feature-card h4 {
  font-size: 1.3rem;
  color: var(--text-color);
  margin-bottom: 1rem;
  font-family: var(--heading-font);
  font-weight: 600;
}

.feature-card p {
  color: var(--light-text);
  line-height: 1.7;
}

/* 注释对话框样式 */
.annotation-dialog {
  background: #f8f9fa;
  border-top: 1px solid #e9ecef;
  padding: 1.5rem;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.annotation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e9ecef;
}

.annotation-header h5 {
  font-size: 1.2rem;
  color: var(--text-color);
  font-family: var(--heading-font);
  font-weight: 600;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: var(--light-text);
  cursor: pointer;
  padding: 0.2rem;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: #e9ecef;
  color: var(--text-color);
}

.annotation-body {
  line-height: 1.7;
}

.annotation-section {
  margin-bottom: 1.5rem;
}

.annotation-section h6 {
  font-size: 1rem;
  color: var(--text-color);
  font-family: var(--heading-font);
  font-weight: 600;
  margin-bottom: 0.8rem;
  position: relative;
  padding-left: 0.8rem;
}

.annotation-section h6::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  width: 3px;
  background-color: var(--primary-color);
}

.annotation-section p {
  color: var(--text-color);
  margin-bottom: 0;
}

.annotation-notes {
  margin-top: 1rem;
}

.note-item {
  margin-bottom: 1rem;
  padding: 0.8rem;
  background: white;
  border-radius: 4px;
  border-left: 3px solid var(--primary-color);
}

.note-item:last-child {
  margin-bottom: 0;
}

.note-line {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--text-color);
  margin-bottom: 0.4rem;
}

.note-content {
  font-size: 0.9rem;
  color: var(--light-text);
  margin-bottom: 0;
  line-height: 1.6;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .section-title h2 {
    font-size: 2rem;
  }
  
  .content-container {
    padding: 0 1rem;
  }
  
  .introduction-section,
  .interactive-section {
    padding: 1.5rem;
  }
  
  .poetry-grid,
  .interactive-features {
    gap: 1.5rem;
  }
  
  .poetry-card {
    margin-bottom: 1rem;
  }
  
  .poetry-actions {
    flex-direction: column;
  }
  
  .back-button {
    padding: 0.6rem 1.2rem;
    font-size: 0.85rem;
  }
  
  .annotation-dialog {
    padding: 1.2rem;
  }
  
  .annotation-header h5 {
    font-size: 1.1rem;
  }
  
  .annotation-section h6 {
    font-size: 0.95rem;
  }
}

@media (max-width: 480px) {
  .section-title h2 {
    font-size: 1.8rem;
  }
  
  .poetry-grid {
    grid-template-columns: 1fr;
  }
  
  .feature-card {
    padding: 1.5rem;
  }
  
  .note-item {
    padding: 0.6rem;
  }
  
  .note-line {
    font-size: 0.9rem;
  }
  
  .note-content {
    font-size: 0.85rem;
  }
}
</style>