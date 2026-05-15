import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getAllUsers } from '../services/api';
import { getUserData } from '../utils/storage';

function AdminPanel() {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState('');
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const userData = getUserData();
    setUser(userData);
    
    // VULNERABILITY #2: Client-side authorization check only
    // User can bypass this by modifying localStorage
    if (userData?.role !== 'admin') {
      // Should redirect, but let's allow it for training purposes
      console.warn('Non-admin user accessing admin panel!');
    }
    
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      // VULNERABILITY #2: Admin endpoint has no server-side authorization
      const response = await getAllUsers();
      setUsers(response.data.users);
    } catch (err) {
      setError('Failed to load users');
      console.error('Error loading users:', err);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <h1 className="text-xl font-bold text-gray-900">SecureTask - Admin Panel</h1>
            <button
              onClick={() => navigate('/dashboard')}
              className="text-blue-500 hover:underline"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-6">All Users</h2>
          
          {/* VULNERABILITY #2: Showing this only works if role is admin in localStorage */}
          {user?.role !== 'admin' && (
            <div className="mb-4 p-4 bg-yellow-100 border border-yellow-300 text-yellow-800 rounded">
              ⚠️ You are not an admin, but you can still access this page due to missing server-side authorization!
            </div>
          )}

          {error && (
            <div className="mb-4 p-4 bg-red-100 text-red-700 rounded">
              {error}
            </div>
          )}

          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Email
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Role
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Password
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {users.map((u) => (
                  <tr key={u.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {u.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {u.name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {u.email}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`px-2 py-1 rounded ${
                        u.role === 'admin' 
                          ? 'bg-purple-100 text-purple-800' 
                          : 'bg-gray-100 text-gray-800'
                      }`}>
                        {u.role}
                      </span>
                    </td>
                    {/* VULNERABILITY #2 & #5: Displaying plain text passwords from API */}
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-red-600 font-mono">
                      {u.password || 'N/A'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded">
            <h3 className="font-semibold text-red-800 mb-2">🚨 Security Issues on This Page:</h3>
            <ul className="list-disc list-inside text-sm text-red-700 space-y-1">
              <li>No server-side authorization check - anyone can access this endpoint</li>
              <li>Passwords are visible in plain text</li>
              <li>Client-side role check can be bypassed</li>
              <li>Sensitive user data exposed without proper access control</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminPanel;
