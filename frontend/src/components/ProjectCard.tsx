import { Link } from 'react-router-dom';
import { Project } from '../types';

interface ProjectCardProps {
  project: Project;
}

const categoryMap: Record<string, { label: string; color: string }> = {
  education: { label: '助学', color: 'bg-blue-100 text-blue-800' },
  elderly: { label: '助老', color: 'bg-orange-100 text-orange-800' },
  medical: { label: '医疗', color: 'bg-green-100 text-green-800' },
  disaster: { label: '救灾', color: 'bg-red-100 text-red-800' },
  environment: { label: '环保', color: 'bg-emerald-100 text-emerald-800' },
  other: { label: '其他', color: 'bg-gray-100 text-gray-800' },
};

const ProjectCard: React.FC<ProjectCardProps> = ({ project }) => {
  const category = categoryMap[project.category] || categoryMap.other;

  return (
    <Link to={`/projects/${project.id}`} className="block">
      <div className="bg-white rounded-xl shadow-sm hover:shadow-lg transition-shadow overflow-hidden">
        <div className="h-48 bg-gradient-to-br from-primary-400 to-primary-600 flex items-center justify-center">
          <svg className="w-16 h-16 text-white opacity-80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        </div>
        <div className="p-5">
          <div className="flex items-center justify-between mb-3">
            <span className={`px-3 py-1 rounded-full text-xs font-medium ${category.color}`}>
              {category.label}
            </span>
            {project.status === 'completed' && (
              <span className="px-3 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                已完成
              </span>
            )}
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2 line-clamp-2">
            {project.title}
          </h3>
          <p className="text-gray-500 text-sm mb-4 line-clamp-2">
            {project.description}
          </p>
          <div className="mb-3">
            <div className="flex justify-between text-sm mb-1">
              <span className="text-gray-500">已筹 ¥{project.currentAmount?.toLocaleString()}</span>
              <span className="text-primary-600 font-medium">{project.progress}%</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-primary-600 h-2 rounded-full transition-all"
                style={{ width: `${project.progress}%` }}
              />
            </div>
          </div>
          <div className="flex justify-between text-sm text-gray-500">
            <span>目标 ¥{project.targetAmount?.toLocaleString()}</span>
            <span>{project.organization?.name}</span>
          </div>
        </div>
      </div>
    </Link>
  );
};

export default ProjectCard;
