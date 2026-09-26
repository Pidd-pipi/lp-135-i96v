import React, { useState } from 'react';
import { Project } from '../types';
import { donationAPI } from '../api';
import { useAuth } from '../context/AuthContext';

interface DonationModalProps {
  project: Project;
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

const DonationModal: React.FC<DonationModalProps> = ({ project, isOpen, onClose, onSuccess }) => {
  const { user } = useAuth();
  const [amount, setAmount] = useState('');
  const [paymentMethod, setPaymentMethod] = useState<'wechat' | 'alipay' | 'bank'>('alipay');
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || parseFloat(amount) <= 0) {
      alert('请输入有效金额');
      return;
    }

    setLoading(true);
    try {
      await donationAPI.createDonation({
        projectId: project.id,
        amount: parseFloat(amount),
        paymentMethod,
        isAnonymous,
        message,
      });
      alert('捐赠成功！感谢您的爱心');
      onSuccess();
      onClose();
    } catch (error: any) {
      alert(error.response?.data?.message || '捐赠失败，请重试');
    } finally {
      setLoading(false);
    }
  };

  const quickAmounts = [10, 50, 100, 500, 1000];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl p-8 max-w-md w-full mx-4">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">爱心捐赠</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="mb-6 p-4 bg-gray-50 rounded-lg">
          <p className="text-sm text-gray-600">正在为项目捐赠</p>
          <p className="font-semibold text-gray-900">{project.title}</p>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-3">捐赠金额</label>
            <div className="grid grid-cols-5 gap-2 mb-3">
              {quickAmounts.map((val) => (
                <button
                  key={val}
                  type="button"
                  onClick={() => setAmount(String(val))}
                  className={`py-2 rounded-lg border ${amount === String(val) ? 'border-primary-600 bg-primary-50 text-primary-600' : 'border-gray-200 hover:border-primary-400'}`}
                >
                  ¥{val}
                </button>
              ))}
            </div>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">¥</span>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="自定义金额"
                className="w-full pl-8 pr-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                min="1"
              />
            </div>
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-3">支付方式</label>
            <div className="grid grid-cols-3 gap-3">
              <button
                type="button"
                onClick={() => setPaymentMethod('wechat')}
                className={`py-3 rounded-lg border ${paymentMethod === 'wechat' ? 'border-green-500 bg-green-50 text-green-600' : 'border-gray-200'}`}
              >
                微信支付
              </button>
              <button
                type="button"
                onClick={() => setPaymentMethod('alipay')}
                className={`py-3 rounded-lg border ${paymentMethod === 'alipay' ? 'border-blue-500 bg-blue-50 text-blue-600' : 'border-gray-200'}`}
              >
                支付宝
              </button>
              <button
                type="button"
                onClick={() => setPaymentMethod('bank')}
                className={`py-3 rounded-lg border ${paymentMethod === 'bank' ? 'border-gray-700 bg-gray-50 text-gray-700' : 'border-gray-200'}`}
              >
                银行卡
              </button>
            </div>
          </div>

          <div className="mb-6">
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={isAnonymous}
                onChange={(e) => setIsAnonymous(e.target.checked)}
                className="w-4 h-4 text-primary-600 rounded"
              />
              <span className="ml-2 text-sm text-gray-600">匿名捐赠</span>
            </label>
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">留言（选填）</label>
            <textarea
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="写下您想对受助者说的话..."
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
              rows={3}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-primary-600 text-white py-4 rounded-lg font-semibold hover:bg-primary-700 disabled:opacity-50"
          >
            {loading ? '处理中...' : `确认捐赠 ¥${amount || '0'}`}
          </button>
        </form>
      </div>
    </div>
  );
};

export default DonationModal;
