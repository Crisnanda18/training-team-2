import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { getCurrentUser, updateProfile } from '../services/api';
import { VITE_DEBUG_MODE } from '../config';

function Profile() {
  const [user, setUser] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    bio: ''
  });
  const [message, setMessage] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const userData = await getCurrentUser();
        setUser(userData);
        setFormData({
          name: userData.data.name || '',
          bio: userData.data.bio || ''
        });
      } catch (error) {
        if (VITE_DEBUG_MODE) {
          console.error('Failed to fetch user data:', error);
        }
      }
    };
    
    fetchUserData();
  }, []);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      // VULNERABILITY #2: No authorization check - can update any user's profile
      const currentUser = await getCurrentUser();
      if (currentUser.data.id !== user.data.id) {
        setMessage('You are not authorized to update this profile.');
        return;
      }
      const response = await updateProfile(user.data.id, formData);
      
      // VULNERABILITY #5: Updating localStorage with potentially sensitive data
      setMessage('Profile updated successfully!');
    } catch (error) {
      setMessage('Failed to update profile');
      if (VITE_DEBUG_MODE) {
        console.error('Update error:', error);
      }
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <h1 className="text-xl font-bold text-gray-900">SecureTask - Profile</h1>
            <button
              onClick={() => navigate('/dashboard')}
              className="text-blue-500 hover:underline"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-6">User Profile</h2>

          {message && (
            <div className={`mb-4 px-4 py-3 rounded ${
              message.includes('success') 
                ? 'bg-green-100 text-green-700' 
                : 'bg-red-100 text-red-700'
            }`}>
              {message}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-gray-700 font-semibold mb-2">
                Email (Read-only)
              </label>
              <input
                type="email"
                value={user?.data.email || ''}
                className="w-full px-4 py-2 border border-gray-300 rounded bg-gray-100"
                disabled
              />
            </div>

            <div>
              <label className="block text-gray-700 font-semibold mb-2">
                Name
              </label>
              <input
                type="text"
                name="name"
                value={formData.name}
                onChange={handleChange}
                className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
                required
              />
            </div>

            <div>
              <label className="block text-gray-700 font-semibold mb-2">
                Bio
              </label>
              {/* VULNERABILITY #3: No sanitization - XSS possible in bio field */}
              <textarea
                name="bio"
                value={formData.bio}
                onChange={handleChange}
                placeholder="Tell us about yourself... (Try: <img src=x onerror=alert('XSS')>)"
                className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
                rows="4"
              />
            </div>

            <button
              type="submit"
              className="w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
            >
              Update Profile
            </button>
          </form>

          {/* Display current bio with XSS vulnerability */}
          {user?.data.bio && (
            <div className="mt-6 p-4 bg-gray-50 rounded">
              <h3 className="font-semibold mb-2">Current Bio:</h3>
              {/* VULNERABILITY #3: Rendering unsanitized HTML 
              fixed: delete dangerouslySetInnerHTML */}
              <div>{user.data.bio}</div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default Profile;
