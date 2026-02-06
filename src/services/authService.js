import { request } from './api';

/**
 * 用户认证服务
 */

// 登录功能
export const login = async (usernameOrEmail, password) => {
  try {
    const response = await request('/auth/login', 'POST', {
      username: usernameOrEmail,
      password
    });
    
    if (response.success) {
      // 保存用户信息和token到本地存储
      localStorage.setItem('access_token', response.access_token);
      localStorage.setItem('user', JSON.stringify(response.user));
    }
    
    return response;
  } catch (error) {
    console.error('登录请求错误:', error);
    throw error;
  }
};

// 注册功能
export const register = async (username, email, password) => {
  try {
    const response = await request('/auth/register', 'POST', {
      username,
      email,
      password
    });
    
    return response;
  } catch (error) {
    console.error('注册请求错误:', error);
    throw error;
  }
};

// 登出功能
export const logout = () => {
  localStorage.removeItem('access_token');
  localStorage.removeItem('user');
  return { success: true, message: '登出成功' };
};

// 获取当前用户信息
export const getCurrentUser = () => {
  const userStr = localStorage.getItem('user');
  if (userStr) {
    try {
      return JSON.parse(userStr);
    } catch (error) {
      console.error('解析用户信息失败:', error);
      return null;
    }
  }
  return null;
};

// 检查用户是否已登录
export const isAuthenticated = () => {
  const token = localStorage.getItem('access_token');
  return token !== null;
};

// 检查用户是否为管理员
export const isAdmin = () => {
  const user = getCurrentUser();
  return user && user.role === 'admin';
};

// 更新用户资料
export const updateProfile = async (userData) => {
  try {
    const response = await request('/auth/profile', 'PUT', userData);
    if (response.success) {
      // 更新本地存储中的用户信息
      const currentUser = getCurrentUser();
      const updatedUser = { ...currentUser, ...userData };
      localStorage.setItem('user', JSON.stringify(updatedUser));
    }
    return response;
  } catch (error) {
    console.error('更新资料请求错误:', error);
    throw error;
  }
};

export default {
  login,
  register,
  logout,
  getCurrentUser,
  isAuthenticated,
  isAdmin,
  updateProfile
};