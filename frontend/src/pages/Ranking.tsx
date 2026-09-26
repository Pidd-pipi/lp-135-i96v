import { useState, useEffect } from 'react';
import { rankingAPI } from '../api';
import { RankingItem } from '../types';

const Ranking = () => {
  const [activeTab, setActiveTab] = useState<'donation' | 'service'>('donation');
  const [donationRankings, setDonationRankings] = useState<RankingItem[]>([]);
  const [serviceRankings, setServiceRankings] = useState<RankingItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadRankings();
  }, []);

  const loadRankings = async () => {
    try {
      const [donationRes, serviceRes] = await Promise.all([
        rankingAPI.getDonationRanking(20),
        rankingAPI.getServiceRanking(20),
      ]);
      setDonationRankings(donationRes.data.rankings);
      setServiceRankings(serviceRes.data.rankings);
    } catch (error) {
      console.error('加载排行榜失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const currentRankings = activeTab === 'donation' ? donationRankings : serviceRankings;

  const getRankStyle = (rank: number) => {
    switch (rank) {
      case 1:
        return 'bg-yellow-400 text-white';
      case 2:
        return 'bg-gray-300 text-white';
      case 3:
        return 'bg-orange-400 text-white';
      default:
        return 'bg-gray-100 text-gray-600';
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center mb-12">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">公益排行榜</h1>
        <p className="text-gray-500">感谢每一位爱心人士的付出</p>
      </div>

      <div className="flex justify-center gap-4 mb-8">
        <button
          onClick={() => setActiveTab('donation')}
          className={`px-8 py-3 rounded-xl font-medium transition-colors ${
            activeTab === 'donation'
              ? 'bg-primary-600 text-white'
              : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
          }`}
        >
          捐赠金额榜
        </button>
        <button
          onClick={() => setActiveTab('service')}
          className={`px-8 py-3 rounded-xl font-medium transition-colors ${
            activeTab === 'service'
              ? 'bg-primary-600 text-white'
              : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
          }`}
        >
          服务时长榜
        </button>
      </div>

      {loading ? (
        <div className="text-center py-20">加载中...</div>
      ) : (
        <div className="bg-white rounded-2xl shadow-sm overflow-hidden">
          <div className="grid grid-cols-12 gap-4 px-6 py-4 bg-gray-50 border-b border-gray-100 text-sm font-medium text-gray-500">
            <div className="col-span-2">排名</div>
            <div className="col-span-6">用户</div>
            <div className="col-span-4 text-right">
              {activeTab === 'donation' ? '累计捐赠' : '服务时长'}
            </div>
          </div>

          {currentRankings.length === 0 ? (
            <div className="text-center py-12 text-gray-500">暂无排行数据</div>
          ) : (
            <div className="divide-y divide-gray-100">
              {currentRankings.map((item) => (
                <div key={item.userId} className="grid grid-cols-12 gap-4 px-6 py-4 items-center hover:bg-gray-50">
                  <div className="col-span-2">
                    <div className={`w-10 h-10 rounded-full flex items-center justify-center font-bold ${getRankStyle(item.rank)}`}>
                      {item.rank}
                    </div>
                  </div>
                  <div className="col-span-6 flex items-center gap-3">
                    <div className="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center">
                      <span className="text-primary-600 font-medium">
                        {(item.realName || item.username).charAt(0)}
                      </span>
                    </div>
                    <div>
                      <div className="font-medium text-gray-900">{item.realName || item.username}</div>
                    </div>
                  </div>
                  <div className="col-span-4 text-right">
                    <span className={`text-xl font-bold ${activeTab === 'donation' ? 'text-primary-600' : 'text-green-600'}`}>
                      {activeTab === 'donation'
                        ? `¥${item.totalDonation?.toLocaleString()}`
                        : `${item.serviceHours} 小时`}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default Ranking;
