import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => {
    // 解包后端统一响应 {code, message, data}，code===0 时返回 data
    if (response.data && typeof response.data === 'object' && 'code' in response.data) {
      if (response.data.code === 0) {
        response.data = response.data.data;
      }
    }
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: (data: any) => api.post('/auth/register', data),
  login: (data: any) => api.post('/auth/login', data),
  getMe: () => api.get('/auth/me'),
  updateProfile: (data: any) => api.put('/auth/profile', data),
};

export const projectAPI = {
  getProjects: (params?: any) => api.get('/projects', { params }),
  getProject: (id: string) => api.get(`/projects/${id}`),
  createProject: (data: any) => api.post('/projects', data),
  getMyProjects: () => api.get('/projects/org/my'),
  getProjectUpdates: (projectId: string) => api.get(`/projects/${projectId}/updates`),
  createProjectUpdate: (projectId: string, data: any) => api.post(`/projects/${projectId}/updates`, data),
};

export const donationAPI = {
  createDonation: (data: any) => api.post('/donations', data),
  getMyDonations: (params?: any) => api.get('/donations/my', { params }),
  getCertificate: (id: string) => api.get(`/donations/${id}/certificate`),
};

export const rankingAPI = {
  getDonationRanking: (limit?: number) => api.get('/ranking/donation', { params: { limit } }),
  getServiceRanking: (limit?: number) => api.get('/ranking/service', { params: { limit } }),
  getStats: () => api.get('/ranking/stats'),
};

export const adminAPI = {
  getPendingProjects: () => api.get('/admin/projects/pending'),
  reviewProject: (id: string, data: any) => api.post(`/admin/projects/${id}/review`, data),
  getPendingOrganizations: () => api.get('/admin/organizations/pending'),
  reviewOrganization: (id: string, data: any) => api.post(`/admin/organizations/${id}/review`, data),
};

export default api;
