import { useState, useEffect } from 'react';
import { Link, Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { projectAPI } from '../api';
import { Project } from '../types';

const statusMeta: Record<string, { label: string; className: string }> = {
  pending: { label: '待审核', className: 'bg-yellow-100 text-yellow-800' },
  approved: { label: '筹款中', className: 'bg-green-100 text-green-800' },
  rejected: { label: '已拒绝', className: 'bg-red-100 text-red-800' },
  paused: { label: '已停募', className: 'bg-orange-100 text-orange-800' },
  completed: { label: '已完成', className: 'bg-gray-200 text-gray-700' },
};

const categoryMap: Record<string, string> = {
  education: '助学',
  elderly: '助老',
  medical: '医疗',
  disaster: '救灾',
  environment: '环保',
  other: '其他',
};

const MyProjects = () => {
  const { user } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [actingId, setActingId] = useState<string | null>(null);

  useEffect(() => {
    if (user?.role === 'org') {
      loadProjects();
    }
  }, [user]);

  const loadProjects = async () => {
    try {
      const response = await projectAPI.getMyProjects();
      setProjects(response.data.projects);
    } catch (error) {
      console.error('加载我的项目失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleStatusChange = async (project: Project, action: 'pause' | 'resume') => {
    const confirmed = action === 'pause'
      ? window.confirm(`确定暂停「${project.title}」的筹款吗？停募后公开列表将不再展示该项目。`)
      : window.confirm(`确定重新开放「${project.title}」的筹款吗？`);
    if (!confirmed) return;
    setActingId(project.id);
    try {
      const response = await projectAPI.changeProjectStatus(project.id, action);
      setProjects((prev) =>
        prev.map((p) => (p.id === project.id ? response.data.project : p))
      );
      alert(action === 'pause' ? '项目已暂停筹款' : '项目已重新开放筹款');
    } catch (error: any) {
      alert(error.response?.data?.message || '状态变更失败');
    } finally {
      setActingId(null);
    }
  };

  if (!user || user.role !== 'org') {
    return <Navigate to="/" />;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">我的项目</h1>
        <Link
          to="/create-project"
          className="bg-primary-600 text-white px-5 py-2.5 rounded-lg font-medium hover:bg-primary-700"
        >
          + 发布新项目
        </Link>
      </div>

      {loading ? (
        <div className="text-center py-20">加载中...</div>
      ) : projects.length === 0 ? (
        <div className="bg-white rounded-2xl shadow-sm text-center py-20 text-gray-500">
          暂未发布项目，点击右上角「发布新项目」开始吧
        </div>
      ) : (
        <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
          <table className="w-full">
            <thead>
              <tr className="bg-gray-50 text-left text-sm text-gray-500">
                <th className="px-6 py-4 font-medium">项目</th>
                <th className="px-6 py-4 font-medium">分类</th>
                <th className="px-6 py-4 font-medium">已筹 / 目标</th>
                <th className="px-6 py-4 font-medium">起止日期</th>
                <th className="px-6 py-4 font-medium">状态</th>
                <th className="px-6 py-4 font-medium text-right">操作</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {projects.map((project) => {
                const status = statusMeta[project.status] || { label: project.status, className: 'bg-gray-100 text-gray-700' };
                const fmtDate = (d?: string) => (d ? new Date(d).toLocaleDateString() : '-');
                const busy = actingId === project.id;
                return (
                  <tr key={project.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4">
                      <Link to={`/projects/${project.id}`} className="font-medium text-gray-900 hover:text-primary-600">
                        {project.title}
                      </Link>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">{categoryMap[project.category] || project.category}</td>
                    <td className="px-6 py-4 text-sm">
                      <span className="text-primary-600 font-semibold">¥{project.currentAmount?.toLocaleString()}</span>
                      <span className="text-gray-400"> / ¥{project.targetAmount?.toLocaleString()}</span>
                      <span className="ml-2 text-gray-400">{project.progress}%</span>
                    </td>
                    <td className="px-6 py-4 text-sm text-gray-600">
                      {fmtDate(project.startDate)} ~ {fmtDate(project.endDate)}
                    </td>
                    <td className="px-6 py-4">
                      <span className={`px-3 py-1 rounded-full text-xs font-medium ${status.className}`}>
                        {status.label}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex justify-end gap-2">
                        {project.status !== 'completed' && (
                          <Link
                            to={`/projects/${project.id}/edit`}
                            className="px-3 py-1.5 text-sm bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100"
                          >
                            编辑
                          </Link>
                        )}
                        {project.status === 'approved' && (
                          <button
                            onClick={() => handleStatusChange(project, 'pause')}
                            disabled={busy}
                            className="px-3 py-1.5 text-sm bg-orange-50 text-orange-600 rounded-lg hover:bg-orange-100 disabled:opacity-50"
                          >
                            {busy ? '处理中...' : '暂停筹款'}
                          </button>
                        )}
                        {project.status === 'paused' && (
                          <button
                            onClick={() => handleStatusChange(project, 'resume')}
                            disabled={busy}
                            className="px-3 py-1.5 text-sm bg-green-50 text-green-600 rounded-lg hover:bg-green-100 disabled:opacity-50"
                          >
                            {busy ? '处理中...' : '重新开放'}
                          </button>
                        )}
                        {project.status === 'completed' && (
                          <span className="px-3 py-1.5 text-sm text-gray-400">不可修改</span>
                        )}
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      <p className="text-sm text-gray-400 mt-4">
        提示：暂停筹款后项目不会出现在公开列表中，但项目详情仍可访问，历史捐赠与电子凭证均保留。
      </p>
    </div>
  );
};

export default MyProjects;
