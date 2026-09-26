import { useState, useEffect } from 'react';
import { useNavigate, Navigate, useParams } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { projectAPI } from '../api';
import { Project } from '../types';

const CreateUpdate = () => {
  const { user } = useAuth();
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState(projectId || '');
  const [formData, setFormData] = useState({
    title: '',
    content: '',
  });
  const [loading, setLoading] = useState(false);
  const [loadingProjects, setLoadingProjects] = useState(true);

  useEffect(() => {
    if (user?.role === 'org') {
      loadMyProjects();
    }
  }, [user]);

  const loadMyProjects = async () => {
    try {
      const response = await projectAPI.getMyProjects();
      setProjects(response.data.projects);
    } catch (error) {
      console.error('加载项目列表失败:', error);
    } finally {
      setLoadingProjects(false);
    }
  };

  if (!user || user.role !== 'org') {
    return <Navigate to="/" />;
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProject) {
      alert('请选择项目');
      return;
    }
    setLoading(true);
    try {
      await projectAPI.createProjectUpdate(selectedProject, {
        ...formData,
        images: [],
      });
      alert('动态发布成功');
      navigate(`/projects/${selectedProject}`);
    } catch (error: any) {
      alert(error.response?.data?.message || '发布失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  if (loadingProjects) {
    return <div className="text-center py-20">加载中...</div>;
  }

  return (
    <div className="max-w-3xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">发布项目动态</h1>

      <form onSubmit={handleSubmit} className="bg-white rounded-2xl shadow-sm p-8 space-y-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">选择项目 *</label>
          <select
            name="projectId"
            value={selectedProject}
            onChange={(e) => setSelectedProject(e.target.value)}
            className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            required
          >
            <option value="">请选择要发布动态的项目</option>
            {projects.map((project) => (
              <option key={project.id} value={project.id}>
                {project.title} ({project.status === 'approved' ? '已通过' : project.status === 'pending' ? '审核中' : '已完成'})
              </option>
            ))}
          </select>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">动态标题 *</label>
          <input
            type="text"
            name="title"
            value={formData.title}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            placeholder="请输入动态标题"
            required
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">动态内容 *</label>
          <textarea
            name="content"
            value={formData.content}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
            placeholder="请详细描述项目进展情况、取得的成果等"
            rows={8}
            required
          />
        </div>

        <div className="flex gap-4 pt-4">
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="flex-1 px-6 py-3 border border-gray-200 rounded-lg font-medium text-gray-700 hover:bg-gray-50"
          >
            取消
          </button>
          <button
            type="submit"
            disabled={loading}
            className="flex-1 bg-primary-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-primary-700 disabled:opacity-50"
          >
            {loading ? '发布中...' : '发布动态'}
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateUpdate;
