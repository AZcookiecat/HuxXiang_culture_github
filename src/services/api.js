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

    if (data) {
      options.body = JSON.stringify(data);
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
    const result = await response.json();

    if (!response.ok) {
      throw new Error(result.error || `请求失败: ${response.status}`);
    }

    return result;
  } catch (error) {
    console.error('API请求错误:', error);
    throw error;
  }
};