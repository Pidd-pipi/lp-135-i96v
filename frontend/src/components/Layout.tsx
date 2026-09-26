import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Layout = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-8">
              <NavLink to="/" className="text-2xl font-bold text-primary-600">
                GiveTrack 公益捐赠
              </NavLink>
              <nav className="hidden md:flex space-x-6">
                <NavLink to="/" className={({ isActive }) => isActive ? 'text-primary-600 font-medium' : 'text-gray-600 hover:text-primary-600'}>
                  首页
                </NavLink>
                <NavLink to="/projects" className={({ isActive }) => isActive ? 'text-primary-600 font-medium' : 'text-gray-600 hover:text-primary-600'}>
                  公益项目
                </NavLink>
                <NavLink to="/ranking" className={({ isActive }) => isActive ? 'text-primary-600 font-medium' : 'text-gray-600 hover:text-primary-600'}>
                  排行榜
                </NavLink>
              </nav>
            </div>
            <div className="flex items-center space-x-4">
              {user ? (
                <>
                  {user.role === 'org' && (
                    <>
                      <NavLink to="/create-project" className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                        发布项目
                      </NavLink>
                      <NavLink to="/create-update" className="text-gray-600 hover:text-primary-600">
                        发布动态
                      </NavLink>
                    </>
                  )}
                  {user.role === 'admin' && (
                    <NavLink to="/admin" className="text-gray-600 hover:text-primary-600">
                      管理后台
                    </NavLink>
                  )}
                  <NavLink to="/profile" className="text-gray-600 hover:text-primary-600">
                    {user.realName || user.username}
                  </NavLink>
                  <button onClick={handleLogout} className="text-gray-600 hover:text-primary-600">
                    退出
                  </button>
                </>
              ) : (
                <>
                  <NavLink to="/login" className="text-gray-600 hover:text-primary-600">
                    登录
                  </NavLink>
                  <NavLink to="/register" className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700">
                    注册
                  </NavLink>
                </>
              )}
            </div>
          </div>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>
      <footer className="bg-white border-t mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center text-gray-500 text-sm">
            © 2026 GiveTrack 公益捐赠 版权所有
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Layout;
