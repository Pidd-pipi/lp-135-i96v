import { useState, useEffect } from 'react';
import { Navigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { projectAPI } from '../api';
import { Project } from '../types';

const categories = [
  { value: 'education', label: '助学' },
  { value: 'elderly', label: '助老' },
  { value: 'medical', label: '医疗' },
  { value: 'disaster', label: '救灾' },
  { value: 'environment', label: '环保' },
  { value: 'other', label: '其他' },
];

const categoryLabel = (v: string) => categories.find((c) => c.value === v)?.label || v;

// statusMeta 后台各状态的中文文案与徽标样式。
const statusMeta: Record<string, { label: string; cls: string }> = {
  pending: { label: '审核中', cls: 'bg-amber-100 text-amber-700' },
  approved: { label: '筹款中', cls: 'bg-green-100 text-green-700' },
  paused: { label: '暂停筹款', cls: 'bg-orange-100 text-orange-700' },
  completed: { label: '已完成', cls: 'bg-gray-200 text-gray-700' },
  rejected: { label: '已驳回', cls: 'bg-red-100 text-red-700' },
};

// toDateInput 把后端时间格式化为 <input type="date"> 需要的 YYYY-MM-DD。
const toDateInput = (v?: string) => (v ? v.slice(0, 10) : '');

interface EditForm {
  description: string;
  category: string;
  targetAmount: string;
  executionPlan: string;
  startDate: string;
  endDate: string;
}

const MyProjects = () => {
  const { user } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState<Project | null>(null);
  const [form, setForm] = useState<EditForm>({
    description: '', category: 'education', targetAmount: '',
    executionPlan: '', startDate: '', endDate: '',
  });
  const [saving, setSaving] = useState(false);
  const [busyId, setBusyId] = useState<string>('');

  useEffect(() => {
    if (user?.role === 'org') loadProjects();
  }, [user]);

  const loadProjects = async () => {
    try {
      const res = await projectAPI.getMyProjects();
      setProjects(res.data.projects);
    } catch (error: any) {
      alert(error.response?.data?.message || '加载项目失败');
    } finally {
      setLoading(false);
    }
  };

  if (!user || user.role !== 'org') {
    return <Navigate to="/" />;
  }

  const startEdit = (p: Project) => {
    setEditing(p);
    setForm({
      description: p.description || '',
      category: p.category,
      targetAmount: String(p.targetAmount),
      executionPlan: p.executionPlan || '',
      startDate: toDateInput(p.startDate),
      endDate: toDateInput(p.endDate),
    });
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editing) return;
    const target = parseFloat(form.targetAmount);
    if (!target || target <= 0) {
      alert('请输入有效的目标金额');
      return;
    }
    if (target < editing.currentAmount) {
      alert(`目标金额不能低于已筹金额（¥${editing.currentAmount.toLocaleString()}）`);
      return;
    }
    if (form.startDate && form.endDate && form.endDate < form.startDate) {
      alert('结束日期不能早于开始日期');
      return;
    }
    setSaving(true);
    try {
      const res = await projectAPI.updateProject(editing.id, {
        description: form.description,
        category: form.category,
        targetAmount: target,
        executionPlan: form.executionPlan,
        startDate: form.startDate || '',
        endDate: form.endDate || '',
      });
      setProjects((prev) => prev.map((p) => (p.id === editing.id ? res.data.project : p)));
      setEditing(null);
      alert('项目信息已保存');
    } catch (error: any) {
      alert(error.response?.data?.message || '保存失败，请重试');
    } finally {
      setSaving(false);
    }
  };

  const handleTogglePause = async (p: Project) => {
    const pause = p.status === 'approved';
    if (!window.confirm(pause ? `确定暂停「${p.title}」的筹款吗？暂停后公开列表将不再展示该项目。` : `确定重新开放「${p.title}」的筹款吗？`)) {
      return;
    }
    setBusyId(p.id);
    try {
      const res = await (pause
        ? projectAPI.pauseProject(p.id)
        : projectAPI.reopenProject(p.id));
      setProjects((prev) => prev.map((item) => (item.id === p.id ? res.data.project : item)));
      alert(pause ? '已暂停筹款' : '已重新开放筹款');
    } catch (error: any) {
      alert(error.response?.data?.message || '操作失败，请重试');
    } finally {
      setBusyId('');
    }
  };

  const handleField = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  if (loading) return <div className="text-center py-20">加载中...</div>;

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">我的项目</h1>
        <Link to="/create-project" className="bg-primary-600 text-white px-5 py-2.5 rounded-lg hover:bg-primary-700">
          发布新项目
        </Link>
      </div>

      <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
        {projects.length === 0 ? (
          <div className="text-center py-16 text-gray-500">
            暂无项目，点击右上角「发布新项目」开始吧
          </div>
        ) : (
          <div className="divide-y divide-gray-100">
            {projects.map((p) => {
              const meta = statusMeta[p.status] || { label: p.status, cls: 'bg-gray-100 text-gray-700' };
              const canEdit = p.status === 'pending' || p.status === 'approved' || p.status === 'paused';
              const canToggle = p.status === 'approved' || p.status === 'paused';
              return (
                <div key={p.id} className="p-6">
                  <div className="flex justify-between items-start gap-6">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-3 mb-2 flex-wrap">
                        <Link to={`/projects/${p.id}`} className="text-lg font-semibold text-gray-900 hover:text-primary-600">
                          {p.title}
                        </Link>
                        <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium ${meta.cls}`}>
                          {meta.label}
                        </span>
                        <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
                          {categoryLabel(p.category)}
                        </span>
                      </div>
                      <div className="grid grid-cols-3 gap-4 text-sm text-gray-600">
                        <div>已筹：<span className="font-semibold text-primary-600">¥{p.currentAmount.toLocaleString()}</span></div>
                        <div>目标：¥{p.targetAmount.toLocaleString()}</div>
                        <div>进度：{p.progress ?? 0}%</div>
                        <div>开始：{toDateInput(p.startDate) || '-'}</div>
                        <div>结束：{toDateInput(p.endDate) || '-'}</div>
                        <div>发布：{new Date(p.createdAt).toLocaleDateString()}</div>
                      </div>
                    </div>
                    <div className="flex flex-col gap-2 shrink-0">
                      {canToggle && (
                        <button
                          onClick={() => handleTogglePause(p)}
                          disabled={busyId === p.id}
                          className={`px-4 py-2 rounded-lg text-sm font-medium disabled:opacity-50 ${
                            p.status === 'approved'
                              ? 'bg-orange-50 text-orange-600 hover:bg-orange-100'
                              : 'bg-green-50 text-green-600 hover:bg-green-100'
                          }`}
                        >
                          {busyId === p.id ? '处理中...' : p.status === 'approved' ? '暂停筹款' : '重新开放'}
                        </button>
                      )}
                      {canEdit ? (
                        <button
                          onClick={() => startEdit(p)}
                          className="px-4 py-2 rounded-lg text-sm font-medium bg-primary-50 text-primary-600 hover:bg-primary-100"
                        >
                          编辑项目
                        </button>
                      ) : (
                        <button
                          disabled
                          title={p.status === 'completed' ? '项目已完成，不可修改' : '当前状态不可修改'}
                          className="px-4 py-2 rounded-lg text-sm font-medium bg-gray-50 text-gray-400 cursor-not-allowed"
                        >
                          不可编辑
                        </button>
                      )}
                      <Link
                        to={`/create-update/${p.id}`}
                        className="px-4 py-2 rounded-lg text-sm font-medium text-center bg-gray-50 text-gray-600 hover:bg-gray-100"
                      >
                        发布动态
                      </Link>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* 编辑弹窗 */}
      {editing && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl p-8 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-bold text-gray-900">编辑项目</h2>
              <button onClick={() => setEditing(null)} className="text-gray-400 hover:text-gray-600 text-2xl leading-none">×</button>
            </div>
            <p className="text-sm text-gray-500 mb-6">
              正在编辑「{editing.title}」，当前已筹 ¥{editing.currentAmount.toLocaleString()}，目标金额不可低于该数额。
            </p>
            <form onSubmit={handleSave} className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">项目分类 *</label>
                <select
                  name="category"
                  value={form.category}
                  onChange={handleField}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500"
                  required
                >
                  {categories.map((c) => <option key={c.value} value={c.value}>{c.label}</option>)}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">目标金额（元）*</label>
                <input
                  type="number"
                  name="targetAmount"
                  value={form.targetAmount}
                  onChange={handleField}
                  min={editing.currentAmount || 1}
                  step="0.01"
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500"
                  required
                />
                <p className="text-xs text-gray-400 mt-1">不得低于已筹金额 ¥{editing.currentAmount.toLocaleString()}</p>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">开始日期</label>
                  <input type="date" name="startDate" value={form.startDate} onChange={handleField}
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500" />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">结束日期</label>
                  <input type="date" name="endDate" value={form.endDate} onChange={handleField}
                    className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500" />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">项目介绍</label>
                <textarea
                  name="description"
                  value={form.description}
                  onChange={handleField}
                  rows={4}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 resize-none"
                  placeholder="请详细描述项目背景、目标和意义"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">执行计划</label>
                <textarea
                  name="executionPlan"
                  value={form.executionPlan}
                  onChange={handleField}
                  rows={5}
                  className="w-full px-4 py-2.5 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 resize-none"
                  placeholder="请描述执行计划、时间节点和预期成果"
                />
              </div>
              <div className="flex gap-4 pt-2">
                <button type="button" onClick={() => setEditing(null)}
                  className="flex-1 px-6 py-3 border border-gray-200 rounded-lg font-medium text-gray-700 hover:bg-gray-50">
                  取消
                </button>
                <button type="submit" disabled={saving}
                  className="flex-1 bg-primary-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-primary-700 disabled:opacity-50">
                  {saving ? '保存中...' : '保存修改'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default MyProjects;
