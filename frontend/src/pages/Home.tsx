import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { projectAPI, rankingAPI } from '../api';
import { Project, RankingItem } from '../types';
import ProjectCard from '../components/ProjectCard';

const Home = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [donationRankings, setDonationRankings] = useState<RankingItem[]>([]);
  const [stats, setStats] = useState({ totalUsers: 0, totalDonationAmount: 0, totalServiceHours: 0 });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [projectsRes, rankingRes, statsRes] = await Promise.all([
        projectAPI.getProjects({ limit: 6 }),
        rankingAPI.getDonationRanking(5),
        rankingAPI.getStats(),
      ]);
      setProjects(projectsRes.data.projects);
      setDonationRankings(rankingRes.data.rankings);
      setStats(statsRes.data.stats);
    } catch (error) {
      console.error('加载数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-20">加载中...</div>;
  }

  return (
    <div>
      <section className="bg-gradient-to-r from-primary-600 to-primary-700 text-white rounded-3xl p-12 mb-12">
        <div className="max-w-3xl">
          <h1 className="text-4xl font-bold mb-4">用爱心点亮希望</h1>
          <p className="text-xl text-primary-100 mb-8">
            每一份捐赠都能改变一个人的命运，让我们一起传递温暖
          </p>
          <div className="flex gap-4">
            <Link to="/projects" className="bg-white text-primary-600 px-8 py-3 rounded-lg font-semibold hover:bg-gray-100">
              立即捐赠
            </Link>
            <Link to="/register" className="border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white hover:text-primary-600">
              加入我们
            </Link>
          </div>
        </div>
      </section>

      <section className="grid grid-cols-3 gap-8 mb-12">
        <div className="bg-white rounded-2xl p-6 text-center shadow-sm">
          <div className="text-4xl font-bold text-primary-600 mb-2">{stats.totalUsers}</div>
          <div className="text-gray-500">注册用户</div>
        </div>
        <div className="bg-white rounded-2xl p-6 text-center shadow-sm">
          <div className="text-4xl font-bold text-green-600 mb-2">¥{stats.totalDonationAmount.toLocaleString()}</div>
          <div className="text-gray-500">累计捐赠</div>
        </div>
        <div className="bg-white rounded-2xl p-6 text-center shadow-sm">
          <div className="text-4xl font-bold text-blue-600 mb-2">{stats.totalServiceHours}</div>
          <div className="text-gray-500">服务时长</div>
        </div>
      </section>

      <section className="mb-12">
        <div className="flex justify-between items-center mb-8">
          <h2 className="text-2xl font-bold text-gray-900">热门公益项目</h2>
          <Link to="/projects" className="text-primary-600 hover:text-primary-700">
            查看全部 →
          </Link>
        </div>
        <div className="grid grid-cols-3 gap-6">
          {projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      </section>

      <section className="grid grid-cols-2 gap-8">
        <div className="bg-white rounded-2xl p-8 shadow-sm">
          <h3 className="text-xl font-bold text-gray-900 mb-6">捐赠排行榜</h3>
          <div className="space-y-4">
            {donationRankings.map((item) => (
              <div key={item.userId} className="flex items-center gap-4">
                <div className={`w-8 h-8 rounded-full flex items-center justify-center font-bold ${
                  item.rank === 1 ? 'bg-yellow-400 text-white' :
                  item.rank === 2 ? 'bg-gray-300 text-white' :
                  item.rank === 3 ? 'bg-orange-400 text-white' :
                  'bg-gray-100 text-gray-600'
                }`}>
                  {item.rank}
                </div>
                <div className="flex-1">
                  <div className="font-medium text-gray-900">{item.realName || item.username}</div>
                </div>
                <div className="text-primary-600 font-semibold">¥{item.totalDonation?.toLocaleString()}</div>
              </div>
            ))}
          </div>
          <Link to="/ranking" className="block text-center text-primary-600 mt-6 hover:text-primary-700">
            查看完整榜单 →
          </Link>
        </div>

        <div className="bg-gradient-to-br from-primary-500 to-primary-700 rounded-2xl p-8 text-white">
          <h3 className="text-xl font-bold mb-4">为什么选择我们？</h3>
          <ul className="space-y-4">
            <li className="flex items-start gap-3">
              <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              <div>
                <div className="font-semibold">公开透明</div>
                <div className="text-primary-100 text-sm">每笔捐赠可追踪，项目进展实时更新</div>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <div className="font-semibold">资金安全</div>
                <div className="text-primary-100 text-sm">专业机构监管，确保专款专用</div>
              </div>
            </li>
            <li className="flex items-start gap-3">
              <svg className="w-6 h-6 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <div className="font-semibold">电子凭证</div>
                <div className="text-primary-100 text-sm">每笔捐赠生成电子证书，可作为公益证明</div>
              </div>
            </li>
          </ul>
        </div>
      </section>
    </div>
  );
};

export default Home;
