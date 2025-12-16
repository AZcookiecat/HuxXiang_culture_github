// 模拟用户数据
let mockUsers = [
  {
    id: 1,
    username: 'admin',
    email: 'admin@qq.com',
    password: 'admin',
    isAdmin: true
  },
  {
    id: 2,
    username: '向烁安',
    email: '18821959907@163.com',
    password: 'xsa003X',
    isAdmin: true
  },
  {
    id: 3,
    username: '朱理婧',
    email: '18707320519@qq.com',
    password: 'zlj3065',
    isAdmin: true
  },
  {
    id: 4,
    username: '向婷',
    email: '19313025728@qq.com',
    password: 'xt008X',
    isAdmin: false
  },
  {
    id: 5,
    username: '宋雨昕',
    email: '17369302819@qq.com',
    password: 'syx0024',
    isAdmin: false
  },
  {
    id: 6,
    username: '伍俊辉',
    email: '18373829281@qq.com',
    password: 'wjh0130',
    isAdmin: false
  },
  {
    id: 7,
    username: '张志阳',
    email: '18166232505@qq.com',
    password: 'zzy0315',
    isAdmin: false
  },
  {
    id: 8,
    username: '李神洲',
    email: '19313179395@qq.com',
    password: 'lsz0278',
    isAdmin: false
  },
]

/**
 * 模拟登录功能
 * @param {string} usernameOrEmail - 用户名或邮箱
 * @param {string} password - 密码
 * @returns {Promise} 登录结果的Promise对象
 */
export const mockLogin = (usernameOrEmail, password) => {
  return new Promise((resolve, reject) => {
    // 模拟网络延迟
    setTimeout(() => {
      const user = mockUsers.find(u => 
        (u.username === usernameOrEmail || u.email === usernameOrEmail) && 
        u.password === password
      )

      if (user) {
        // 登录成功，存储用户信息到本地存储
        const userInfo = {
          id: user.id,
          username: user.username,
          email: user.email,
          isAdmin: user.isAdmin
        }
        
        localStorage.setItem('user', JSON.stringify(userInfo))
        
        resolve({
          success: true,
          user: userInfo,
          message: '登录成功'
        })
      } else {
        reject(new Error('用户名或密码错误'))
      }
    }, 500)
  })
}

/**
 * 模拟注册功能
 * @param {string} username - 用户名
 * @param {string} email - 邮箱
 * @param {string} password - 密码
 * @returns {Promise} 注册结果的Promise对象
 */
export const mockRegister = (username, email, password) => {
  return new Promise((resolve, reject) => {
    // 模拟网络延迟
    setTimeout(() => {
      // 检查用户名或邮箱是否已存在
      const existingUser = mockUsers.find(u => 
        u.username === username || u.email === email
      )

      if (existingUser) {
        reject(new Error('用户名或邮箱已被注册'))
      } else {
        // 创建新用户
        const newUser = {
          id: mockUsers.length + 1,
          username,
          email,
          password,
          isAdmin: false
        }
        
        // 添加到模拟用户数据中
        mockUsers.push(newUser)
        
        resolve({
          success: true,
          message: '注册成功'
        })
      }
    }, 500)
  })
}

/**
 * 模拟登出功能
 * @returns {Promise} 登出结果的Promise对象
 */
export const mockLogout = () => {
  return new Promise((resolve) => {
    // 模拟网络延迟
    setTimeout(() => {
      // 从本地存储移除用户信息
      localStorage.removeItem('user')
      resolve({ success: true, message: '登出成功' })
    }, 300)
  })
}

/**
 * 获取当前登录用户信息
 * @returns {Object|null} 用户信息对象或null
 */
export const getCurrentUser = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch (error) {
      console.error('解析用户信息失败:', error)
      return null
    }
  }
  return null
}

/**
 * 检查用户是否已登录
 * @returns {boolean} 是否已登录
 */
export const isAuthenticated = () => {
  return getCurrentUser() !== null
}

/**
 * 检查用户是否为管理员
 * @returns {boolean} 是否为管理员
 */
export const isAdmin = () => {
  const user = getCurrentUser()
  return user && user.isAdmin
}

/**
 * 更新用户资料
 * @param {number} userId - 用户ID
 * @param {Object} updatedData - 更新的用户数据
 * @returns {Promise} 更新结果的Promise对象
 */
export const updateUserProfile = (userId, updatedData) => {
  return new Promise((resolve, reject) => {
    // 模拟网络延迟
    setTimeout(() => {
      const userIndex = mockUsers.findIndex(u => u.id === userId)
      
      if (userIndex === -1) {
        reject(new Error('用户不存在'))
        return
      }
      
      // 更新用户数据（只更新允许修改的字段）
      const user = mockUsers[userIndex]
      
      // 允许更新的字段
      if (updatedData.username) {
        // 检查新用户名是否已被使用
        const usernameExists = mockUsers.some(u => 
          u.id !== userId && u.username === updatedData.username
        )
        
        if (usernameExists) {
          reject(new Error('用户名已被使用'))
          return
        }
        
        user.username = updatedData.username
      }
      
      // 如果有头像数据，这里可以添加头像更新逻辑
      // 注意：在实际应用中，头像通常会存储在服务器上
      if (updatedData.avatar) {
        // 在这个模拟服务中，我们只记录头像更新，但不实际存储
        console.log('头像已更新')
      }
      
      // 更新模拟用户数据
      mockUsers[userIndex] = user
      
      // 如果用户当前已登录，更新localStorage中的用户信息
      const currentUser = getCurrentUser()
      if (currentUser && currentUser.id === userId) {
        const updatedUserInfo = {
          ...currentUser,
          username: user.username
        }
        localStorage.setItem('user', JSON.stringify(updatedUserInfo))
      }
      
      resolve({
        success: true,
        message: '个人资料更新成功',
        user: {
          id: user.id,
          username: user.username,
          email: user.email,
          isAdmin: user.isAdmin
        }
      })
    }, 500)
  })
}

export default {
  mockLogin,
  mockRegister,
  mockLogout,
  getCurrentUser,
  isAuthenticated,
  isAdmin,
  updateUserProfile,
  users: mockUsers
}