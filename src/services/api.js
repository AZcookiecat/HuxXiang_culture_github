// API服务配置
const API_BASE_URL = '/api'; // 假设API基础路径为/api

/**
 * 处理HTTP请求的基础函数
 * @param {string} endpoint - API端点
 * @param {string} method - HTTP方法
 * @param {object} data - 请求数据
 * @returns {Promise} 返回处理后的响应
 */
export const request = async (endpoint, method = 'GET', data = null) => {
  try {
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include' // 包含cookie，用于会话管理
    };

    // 如果有存储的token，则添加到请求头
    const token = localStorage.getItem('access_token');
    if (token) {
      options.headers.Authorization = `Bearer ${token}`;
    }

    if (data) {
      options.body = JSON.stringify(data);
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, options);

    // 检查响应是否为JSON格式
    const contentType = response.headers.get('content-type');
    if (!contentType || !contentType.includes('application/json')) {
      if (response.status >= 200 && response.status < 300) {
        // 如果状态码表示成功但不是JSON响应，尝试获取文本并返回默认格式
        const text = await response.text();
        console.warn(`Non-JSON response received for ${method} ${endpoint}:`, text);
        return { success: true, data: text || {} };
      } else {
        // 如果状态码表示失败且不是JSON响应，抛出错误
        throw new Error(`服务器返回错误响应，状态码: ${response.status}`);
      }
    }

    const result = await response.json();

    if (!response.ok) {
      // 如果响应状态是401（未授权），可能是token过期
      if (response.status === 401) {
        // 清除本地存储的无效token
        localStorage.removeItem('access_token');
        localStorage.removeItem('user');
        // 触发登出事件
        window.dispatchEvent(new Event('logout'));
      }
      throw new Error(result.error || result.message || `请求失败: ${response.status}`);
    }

    return result;
  } catch (error) {
    console.error('API请求错误:', error);
    
    // 如果是JSON解析错误，提供更友好的错误信息
    if (error instanceof SyntaxError || error.message.includes('JSON')) {
      throw new Error('服务器返回的数据格式错误，请联系管理员');
    }
    
    throw error;
  }
};