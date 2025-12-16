<template>
  <div class="knowledge-graph-container">
    <div class="knowledge-graph-wrapper">
      <!-- 图谱头部 -->
      <div class="graph-header">
        <h1>湖湘文化知识图谱</h1>
        <p>探索湖湘文化中的人物、事件、地点和文化元素之间的关联关系</p>
      </div>

      <!-- 控制按钮 -->
      <div class="graph-controls">
        <button class="btn btn-primary" @click="resetGraph">重置视图</button>
        <button class="btn btn-outline" @click="zoomIn">放大</button>
        <button class="btn btn-outline" @click="zoomOut">缩小</button>
        <button class="btn btn-outline" @click="toggleAnimation">{{ isAnimating ? '暂停动画' : '播放动画' }}</button>
      </div>

      <!-- 图谱画布 -->
      <div class="graph-canvas" ref="graphCanvas" @wheel.prevent="handleWheel">
        <svg id="graph-svg" ref="graphSvg" @mousedown="startDrag" @mousemove="handleDrag" @mouseup="endDrag" @mouseleave="endDrag"></svg>
        <div v-if="isLoading" class="graph-loading">
          <div class="loading-spinner"></div>
        </div>
        <div v-if="hasError" class="graph-error">
          <div class="graph-error-content">
            <div class="graph-error-icon">
              <i class="fas fa-exclamation-circle"></i>
            </div>
            <h3>加载失败</h3>
            <p>无法加载知识图谱数据，请稍后重试。</p>
            <button class="btn btn-primary" @click="loadGraphData">重试</button>
          </div>
        </div>
        <div class="zoom-hint">
          <i class="fas fa-search-plus"></i> 鼠标滚轮缩放 | <i class="fas fa-hand-pointer"></i> 拖拽移动
        </div>
        <!-- 图例 -->
        <div class="graph-legend">
          <h4>图例</h4>
          <div class="legend-item">
            <div class="legend-color" style="background-color: #2ecc71;"></div>
            <span>历史人物</span>
          </div>
          <div class="legend-item">
            <div class="legend-color" style="background-color: #e74c3c;"></div>
            <span>历史遗迹</span>
          </div>
          <div class="legend-item">
            <div class="legend-color" style="background-color: #9b59b6;"></div>
            <span>遗产文化</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed, inject } from 'vue'
import '../assets/css/knowledge-graph.css'

export default {
  name: 'KnowledgeGraphPage',

  setup() {
    // 获取全局事件总线
    const eventBus = inject('$eventBus')
    
    // 显示提示消息
    const showAlert = (message, type = 'info') => {
      if (eventBus) {
        eventBus.emit('show-alert', { message, type })
      } else {
        // 如果没有eventBus，不显示任何弹窗，只在控制台打印消息
        console.log(`[${type}] ${message}`)
      }
    }
    
    // 引用
    const graphCanvas = ref(null)
    const graphSvg = ref(null)

    // 状态
    const nodes = ref([])
    const links = ref([])
    const isLoading = ref(true)
    const hasError = ref(false)
    const isAnimating = ref(true)
    const isDragging = ref(false)
    const dragStart = ref({ x: 0, y: 0 })
    const transform = ref({ x: 0, y: 0, scale: 1.5 })

    // 模拟数据 - 按照层次结构组织：湖湘文化 -> 分类节点 -> 具体节点
    const mockNodes = [
      // 顶层节点
      { id: 1, name: '湖湘文化', level: 1, x: 0, y: 0 },
      
      // 分类节点 (第二层)
      { id: 2, name: '历史人物', level: 2, x: -200, y: 0 },
      { id: 3, name: '历史遗迹', level: 2, x: 200, y: 0 },
      { id: 4, name: '遗产文化', level: 2, x: 0, y: 200 },
      
      // 历史人物节点 (第三层)
      { id: 5, name: '毛泽东', level: 3, x: -300, y: -80 },
      { id: 6, name: '曾国藩', level: 3, x: -300, y: 80 },
      { id: 7, name: '刘少奇', level: 3, x: -200, y: -150 },
      { id: 8, name: '彭德怀', level: 3, x: -200, y: 150 },
      { id: 9, name: '左宗棠', level: 3, x: -100, y: 150 },
      { id: 10, name: '谭嗣同', level: 3, x: -100, y: 80 },
      { id: 11, name: '魏源', level: 3, x: -100, y: -80 },
      { id: 12, name: '贾谊', level: 3, x: -100, y: -150 },
      { id: 13, name: '屈原', level: 3, x: -200, y: -80 },
      { id: 14, name: '范仲淹', level: 3, x: -300, y: 0 },
      
      // 历史遗迹节点 (第三层)
      { id: 15, name: '岳阳楼', level: 3, x: 300, y: -150 },
      { id: 16, name: '岳麓书院', level: 3, x: 300, y: -80 },
      { id: 17, name: '岳麓山', level: 3, x: 300, y: 0 },
      { id: 18, name: '橘子洲', level: 3, x: 300, y: 80 },
      { id: 19, name: '韶山', level: 3, x: 300, y: 150 },
      { id: 20, name: '湘潭', level: 3, x: 200, y: 150 },
      { id: 21, name: '桃花源', level: 3, x: 200, y: 80 },
      { id: 22, name: '张家界', level: 3, x: 200, y: 0 },
      { id: 23, name: '衡山', level: 3, x: 200, y: -80 },
      { id: 24, name: '船山学社', level: 3, x: 200, y: -150 },
      
      // 遗产文化节点 (第三层)
      { id: 25, name: '湘绣', level: 3, x: -100, y: 250 },
      { id: 26, name: '花鼓戏', level: 3, x: 0, y: 250 },
      { id: 27, name: '湘菜', level: 3, x: 100, y: 250 },
      { id: 28, name: '齐白石', level: 3, x: -150, y: 320 },
      { id: 29, name: '《岳阳楼记》', level: 3, x: 0, y: 320 },
      { id: 30, name: '《楚辞》', level: 3, x: 150, y: 320 },
      { id: 31, name: '《海国图志》', level: 3, x: -75, y: 380 },
      { id: 32, name: '《曾国藩家书》', level: 3, x: 75, y: 380 },
      { id: 33, name: '湘军文化', level: 3, x: 0, y: 450 }
    ]

    const mockLinks = [
      // 湖湘文化连接到分类节点
      { source: 1, target: 2, level: 1 },
      { source: 1, target: 3, level: 1 },
      { source: 1, target: 4, level: 1 },
      
      // 历史人物分类连接具体人物节点
      { source: 2, target: 5, level: 2 },
      { source: 2, target: 6, level: 2 },
      { source: 2, target: 7, level: 2 },
      { source: 2, target: 8, level: 2 },
      { source: 2, target: 9, level: 2 },
      { source: 2, target: 10, level: 2 },
      { source: 2, target: 11, level: 2 },
      { source: 2, target: 12, level: 2 },
      { source: 2, target: 13, level: 2 },
      { source: 2, target: 14, level: 2 },
      
      // 历史遗迹分类连接具体遗迹节点
      { source: 3, target: 15, level: 2 },
      { source: 3, target: 16, level: 2 },
      { source: 3, target: 17, level: 2 },
      { source: 3, target: 18, level: 2 },
      { source: 3, target: 19, level: 2 },
      { source: 3, target: 20, level: 2 },
      { source: 3, target: 21, level: 2 },
      { source: 3, target: 22, level: 2 },
      { source: 3, target: 23, level: 2 },
      { source: 3, target: 24, level: 2 },
      
      // 遗产文化分类连接具体遗产节点
      { source: 4, target: 25, level: 2 },
      { source: 4, target: 26, level: 2 },
      { source: 4, target: 27, level: 2 },
      { source: 4, target: 28, level: 2 },
      { source: 4, target: 29, level: 2 },
      { source: 4, target: 30, level: 2 },
      { source: 4, target: 31, level: 2 },
      { source: 4, target: 32, level: 2 },
      { source: 4, target: 33, level: 2 },
      
      // 一些相关联的节点连接
      { source: 6, target: 32, level: 3 },
      { source: 6, target: 33, level: 3 },
      { source: 11, target: 31, level: 3 },
      { source: 13, target: 30, level: 3 },
      { source: 14, target: 15, level: 3 },
      { source: 14, target: 29, level: 3 },
      { source: 15, target: 29, level: 3 },
      { source: 16, target: 17, level: 3 }
    ]

    // 加载图谱数据
    const loadGraphData = async () => {
      isLoading.value = true
      hasError.value = false
      
      try {
        // 模拟API请求延迟
        await new Promise(resolve => setTimeout(resolve, 1500))
        
        // 使用模拟数据
        nodes.value = [...mockNodes]
        links.value = [...mockLinks]
        
        // 计算节点位置
        calculateNodePositions()
        
        // 渲染图谱
        renderGraph()
        
        showAlert('知识图谱加载成功！', 'success')
      } catch (error) {
        console.error('Failed to load graph data:', error)
        hasError.value = true
        showAlert('知识图谱加载失败，请稍后重试。', 'error')
      } finally {
        isLoading.value = false
      }
    }

    // 重新实现节点排列算法，让子节点围绕父节点分布并减少连线重叠
    const calculateNodePositions = () => {
      // 先清空现有位置信息
      const nodeMap = new Map();
      nodes.value.forEach(node => {
        nodeMap.set(node.id, node);
      });
      
      // 创建父节点到子节点的映射
      const parentToChildren = new Map();
      links.value.forEach(link => {
        if (!parentToChildren.has(link.source)) {
          parentToChildren.set(link.source, []);
        }
        parentToChildren.get(link.source).push(link.target);
      });
      
      // 首先设置顶层节点位置（湖湘文化）
      const rootNode = nodeMap.get(1);
      if (rootNode) {
        rootNode.x = 0;
        rootNode.y = 0;
      }
      
      // 处理第二层节点（分类节点）
      // 围绕根节点以圆形分布
      const categoryNodes = nodes.value.filter(node => node.level === 2);
      const categoryRadius = 180;
      categoryNodes.forEach((node, index) => {
        const angle = (index / categoryNodes.length) * Math.PI * 2;
        node.x = Math.cos(angle) * categoryRadius;
        node.y = Math.sin(angle) * categoryRadius;
        
        // 确保遗产文化节点在下方
        if (node.id === 4) { // 遗产文化
          node.x = 0;
          node.y = 200;
        } else if (node.id === 2) { // 历史人物
          node.x = -200;
          node.y = 0;
        } else if (node.id === 3) { // 历史遗迹
          node.x = 200;
          node.y = 0;
        }
      });
      
      // 处理第三层节点（具体节点）
      // 让每个具体节点围绕其对应的父节点分布
      nodes.value.forEach(node => {
        if (node.level === 3) {
          // 找到当前节点的父节点
          const parentLink = links.value.find(link => link.target === node.id && link.level === 2);
          if (parentLink && nodeMap.has(parentLink.source)) {
            const parentNode = nodeMap.get(parentLink.source);
            const children = parentToChildren.get(parentLink.source) || [];
            const childIndex = children.indexOf(node.id);
            
            if (childIndex !== -1) {
              // 根据子节点索引计算角度位置，让子节点围绕父节点分布
              const radius = 120; // 子节点距离父节点的距离
              // 根据父节点位置调整角度分布，避免连线重叠
              let angleOffset = 0;
              if (parentNode.id === 2) { // 历史人物分类 - 左侧
                angleOffset = Math.PI * 0.75;
              } else if (parentNode.id === 3) { // 历史遗迹分类 - 右侧
                angleOffset = Math.PI * 0.25;
              } else if (parentNode.id === 4) { // 遗产文化分类 - 下方
                angleOffset = 0;
              }
              
              // 调整角度分布，避免子节点之间过于拥挤
              const angleStep = Math.PI * 2 / Math.max(children.length, 6);
              const angle = angleOffset + childIndex * angleStep;
              
              // 计算子节点位置
              node.x = parentNode.x + Math.cos(angle) * radius;
              node.y = parentNode.y + Math.sin(angle) * radius;
            }
          }
        }
      });
    }

    // 渲染图谱
    const renderGraph = () => {
      if (!graphSvg.value) return

      const svg = graphSvg.value
      svg.innerHTML = ''

      // 设置SVG视口
      const { width, height } = graphCanvas.value.getBoundingClientRect()
      svg.setAttribute('width', width)
      svg.setAttribute('height', height)
      
      // 创建变换组
      const g = document.createElementNS('http://www.w3.org/2000/svg', 'g')
      // 调整y坐标使图谱向上移动，距离画布上边框20px
      const translateY = 20 + (height * 0.4); // 20px边距 + 画布高度的40%
      g.setAttribute('transform', `translate(${width / 2}, ${translateY}) scale(${transform.value.scale})`)
      svg.appendChild(g)

      // 创建箭头标记
      createArrowMarkers(svg)

      // 绘制连接线
      links.value.forEach(link => {
        const source = nodes.value.find(n => n.id === link.source)
        const target = nodes.value.find(n => n.id === link.target)
        
        if (source && target) {
          const line = document.createElementNS('http://www.w3.org/2000/svg', 'path')
          
          // 计算贝塞尔曲线控制点
          const dx = target.x - source.x
          const dy = target.y - source.y
          const length = Math.sqrt(dx * dx + dy * dy)
          const tx = source.x + dx * 0.3
          const ty = source.y + dy * 0.3
          const tx2 = source.x + dx * 0.7
          const ty2 = source.y + dy * 0.7
          
          // 设置路径
          const d = `M ${source.x},${source.y} C ${tx},${ty} ${tx2},${ty2} ${target.x},${target.y}`
          line.setAttribute('d', d)
          line.setAttribute('class', `graph-link`)
          // 统一使用浅灰色连接线
          line.setAttribute('stroke', '#95a5a6')
          line.setAttribute('stroke-width', '2px')
          line.setAttribute('fill', 'none')
          line.setAttribute('opacity', '0.7')
          
          // 添加点击事件
          line.addEventListener('click', () => {
            showAlert(`查看关系: ${source.name} -> ${target.name}`, 'info')
          })
          
          g.appendChild(line)
        }
      })

      // 绘制节点
      nodes.value.forEach(node => {
        // 创建节点圆圈
        const circle = document.createElementNS('http://www.w3.org/2000/svg', 'circle')
        circle.setAttribute('cx', node.x)
        circle.setAttribute('cy', node.y)
        
        // 根据文本长度动态设置节点大小
        const textLength = node.name.length;
        let radius = 16; // 默认半径
        
        if (textLength <= 2) {
          radius = 16;
        } else if (textLength <= 4) {
          radius = 20;
        } else if (textLength <= 6) {
          radius = 24;
        } else {
          radius = 28; // 最长文本的半径
        }
        
        // 应用半径和基础样式类
        circle.setAttribute('r', radius);
        circle.setAttribute('class', `graph-node`);
        
        // 根据分类设置颜色
        if (node.level === 1) {
          circle.setAttribute('fill', '#3498db'); // 湖湘文化主节点 - 蓝色
        } else if (node.id === 2) { // 历史人物分类
          circle.setAttribute('fill', '#2ecc71'); // 绿色
        } else if (node.id === 3) { // 历史遗迹分类
          circle.setAttribute('fill', '#e74c3c'); // 红色
        } else if (node.id === 4) { // 遗产文化分类
          circle.setAttribute('fill', '#9b59b6'); // 紫色
        } else {
          // 第三层节点：根据所属分类设置颜色
          const parentLink = links.value.find(link => link.target === node.id && link.level === 2);
          if (parentLink) {
            if (parentLink.source === 2) { // 属于历史人物
              circle.setAttribute('fill', '#2ecc71'); // 绿色
            } else if (parentLink.source === 3) { // 属于历史遗迹
              circle.setAttribute('fill', '#e74c3c'); // 红色
            } else if (parentLink.source === 4) { // 属于遗产文化
              circle.setAttribute('fill', '#9b59b6'); // 紫色
            }
          }
        }
        
        circle.setAttribute('stroke', 'white');
        circle.setAttribute('stroke-width', '2px');
        
        // 添加点击事件
        circle.addEventListener('click', () => {
          showAlert(`查看节点: ${node.name}`, 'info')
        })
        
        // 添加拖拽功能
        let isNodeDragging = false
        circle.addEventListener('mousedown', (e) => {
          e.stopPropagation()
          isNodeDragging = true
        })
        
        document.addEventListener('mousemove', (e) => {
          if (isNodeDragging && graphSvg.value) {
            const rect = graphSvg.value.getBoundingClientRect()
            const scale = transform.value.scale
            const centerX = rect.width / 2
            const centerY = rect.height / 2
            
            node.x = (e.clientX - rect.left - centerX) / scale
            node.y = (e.clientY - rect.top - centerY) / scale
            
            // 重新渲染连接线
            renderGraph()
          }
        })
        
        document.addEventListener('mouseup', () => {
          isNodeDragging = false
        })
        
        g.appendChild(circle)

        // 创建节点文本
        const text = document.createElementNS('http://www.w3.org/2000/svg', 'text')
        text.setAttribute('x', node.x)
        text.setAttribute('y', node.y)
        text.setAttribute('class', 'graph-node-text')
        text.textContent = node.name
        // 设置字体样式
        text.setAttribute('font-size', '14px')
        text.setAttribute('font-family', '"Microsoft YaHei", "SimHei", sans-serif')
        text.setAttribute('fill', 'white')
        text.setAttribute('font-weight', '700')
        text.setAttribute('text-anchor', 'middle')
        text.setAttribute('dominant-baseline', 'middle')
        g.appendChild(text)
      })
    }

    // 创建箭头标记
    const createArrowMarkers = (svg) => {
      for (let i = 1; i <= 5; i++) {
        const marker = document.createElementNS('http://www.w3.org/2000/svg', 'marker')
        marker.setAttribute('id', `arrowhead-${i}`)
        marker.setAttribute('viewBox', '0 0 10 10')
        marker.setAttribute('refX', '8')
        marker.setAttribute('refY', '5')
        marker.setAttribute('markerWidth', '6')
        marker.setAttribute('markerHeight', '6')
        marker.setAttribute('orient', 'auto-start-reverse')
        
        const path = document.createElementNS('http://www.w3.org/2000/svg', 'path')
        path.setAttribute('d', 'M 0 0 L 10 5 L 0 10 z')
        path.setAttribute('class', `arrowhead arrowhead-level-${i}`)
        
        marker.appendChild(path)
        svg.appendChild(marker)
      }
    }

    // 重置图谱视图
    const resetGraph = () => {
      transform.value = { x: 0, y: 0, scale: 1 }
      calculateNodePositions()
      renderGraph()
      showAlert('图谱视图已重置！', 'info')
    }

    // 放大
    const zoomIn = () => {
      transform.value.scale = Math.min(transform.value.scale * 1.2, 3)
      renderGraph()
    }

    // 缩小
    const zoomOut = () => {
      transform.value.scale = Math.max(transform.value.scale / 1.2, 0.5)
      renderGraph()
    }

    // 处理鼠标滚轮缩放
    const handleWheel = (e) => {
      e.preventDefault()
      const delta = e.deltaY > 0 ? -0.1 : 0.1
      transform.value.scale = Math.max(0.5, Math.min(3, transform.value.scale + delta))
      renderGraph()
    }

    // 开始拖拽
    const startDrag = (e) => {
      if (e.target === graphSvg.value) {
        isDragging.value = true
        dragStart.value = { x: e.clientX - transform.value.x, y: e.clientY - transform.value.y }
      }
    }

    // 处理拖拽
    const handleDrag = (e) => {
      if (isDragging.value) {
        transform.value.x = e.clientX - dragStart.value.x
        transform.value.y = e.clientY - dragStart.value.y
        renderGraph()
      }
    }

    // 结束拖拽
    const endDrag = () => {
      isDragging.value = false
    }

    // 切换动画
    const toggleAnimation = () => {
      isAnimating.value = !isAnimating.value
      showAlert(isAnimating.value ? '动画已开始' : '动画已暂停', 'info')
    }

    // 窗口大小变化时重新渲染
    const handleResize = () => {
      if (nodes.value.length > 0) {
        renderGraph()
      }
    }

    // 组件挂载
    onMounted(() => {
      loadGraphData()
      window.addEventListener('resize', handleResize)
    })

    // 组件卸载
    onUnmounted(() => {
      window.removeEventListener('resize', handleResize)
    })

    return {
      graphCanvas,
      graphSvg,
      nodes,
      links,
      isLoading,
      hasError,
      isAnimating,
      loadGraphData,
      resetGraph,
      zoomIn,
      zoomOut,
      handleWheel,
      startDrag,
      handleDrag,
      endDrag,
      toggleAnimation
    }
  }
}
</script>

<style scoped>
/* 知识图谱容器样式补充 */
.knowledge-graph-container {
  padding-top: 20px;
  padding-bottom: 50px;
  background-color: white;
  min-height: 100vh;
}

/* 加载状态和错误状态样式补充 */
.graph-loading,
.graph-error {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(255, 255, 255, 0.9);
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 图谱头部样式补充 */
.graph-header {
  margin-bottom: 30px;
  text-align: center;
  max-width: 1400px;
  margin-left: auto;
  margin-right: auto;
}

/* 控制按钮样式补充 */
.graph-controls {
  margin-bottom: 20px;
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
  max-width: 1400px;
  margin-left: auto;
  margin-right: auto;
}

/* 图谱画布样式 */
.graph-canvas {
    width: 1400px;
    height: 900px;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    overflow: hidden;
    background-color: #f8f9fa;
    position: relative;
  }
  
  /* 知识图谱包装器样式 */
  .knowledge-graph-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    max-width: 1600px;
    margin: 0 auto;
    padding: 20px;
  }

/* 响应式调整 */
@media (max-width: 1450px) {
  .graph-canvas {
    width: 95%;
  }
  
  .knowledge-graph-wrapper {
    padding: 10px;
  }
}

@media (max-width: 768px) {
  .graph-canvas {
    height: 700px !important;
  }
  
  .graph-controls {
    flex-direction: column;
    align-items: center;
  }
  
  .graph-controls .btn {
    width: 200px;
    margin: 5px 0;
  }
}
</style>