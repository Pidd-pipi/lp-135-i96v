import { useState, useEffect } from 'react';
import { projectAPI } from '../api';
import { Project, ProjectCategory } from '../types';
import ProjectCard from '../components/ProjectCard';

const categories = [
  { value: 'all', label: '全部' },
  { value: 'education', label: '助学' },
  { value: 'elderly', label: '助老' },
  { value: 'medical', label: '医疗' },
  { value: 'disaster', label: '救灾' },
  { value: 'environment', label: '环保' },
  { value: 'other', label: '其他' },
];

const Projects = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [category, setCategory] = useState<string>('all');
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({ page: 1, limit: 12, total: 0, totalPages: 0 });

  useEffect(() => {
    loadProjects();
  }, [category, pagination.page]);

  const loadProjects = async () => {
    setLoading(true);
    try {
      const response = await projectAPI.getProjects({
        category: category === 'all' ? undefined : category,
        page: pagination.page,
        limit: pagination.limit,
      });
      setProjects(response.data.projects);
      setPagination(prev => ({
        ...prev,
        total: response.data.total,
        totalPages: response.data.totalPages,
      }));
    } catch (error) {
      console.error('加载项目失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-gray-900">公益项目</h1>
      </div>

      <div className="flex flex-wrap gap-3 mb-8">
        {categories.map((cat) => (
          <button
            key={cat.value}
            onClick={() => {
              setCategory(cat.value);
              setPagination(prev => ({ ...prev, page: 1 }));
            }}
            className={`px-6 py-2 rounded-full text-sm font-medium transition-colors ${
              category === cat.value
                ? 'bg-primary-600 text-white'
                : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
            }`}
          >
            {cat.label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="text-center py-20">加载中...</div>
      ) : projects.length === 0 ? (
        <div className="text-center py-20 text-gray-500">暂无相关项目</div>
      ) : (
        <>
          <div className="grid grid-cols-3 gap-6 mb-8">
            {projects.map((project) => (
              <ProjectCard key={project.id} project={project} />
            ))}
          </div>

          {pagination.totalPages > 1 && (
            <div className="flex justify-center gap-2">
              <button
                onClick={() => setPagination(prev => ({ ...prev, page: Math.max(1, prev.page - 1) }))}
                disabled={pagination.page === 1}
                className="px-4 py-2 rounded-lg border border-gray-200 disabled:opacity-50"
              >
                上一页
              </button>
              {Array.from({ length: pagination.totalPages }, (_, i) => i + 1).map((page) => (
                <button
                  key={page}
                  onClick={() => setPagination(prev => ({ ...prev, page }))}
                  className={`px-4 py-2 rounded-lg ${
                    page === pagination.page
                      ? 'bg-primary-600 text-white'
                      : 'border border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  {page}
                </button>
              ))}
              <button
                onClick={() => setPagination(prev => ({ ...prev, page: Math.min(prev.totalPages, prev.page + 1) }))}
                disabled={pagination.page === pagination.totalPages}
                className="px-4 py-2 rounded-lg border border-gray-200 disabled:opacity-50"
              >
                下一页
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default Projects;
