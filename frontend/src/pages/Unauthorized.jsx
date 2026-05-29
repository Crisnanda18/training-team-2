import { useNavigate } from 'react-router-dom';

function Unauthorized() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <h1 className="text-xl font-bold text-gray-900">SecureTask</h1>
            <button
              onClick={() => navigate('/dashboard')}
              className="text-blue-500 hover:underline"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <h2 className="text-2xl font-semibold text-gray-900 mb-2">
            Unauthorized Access
          </h2>
          <p className="text-gray-600 mb-6">
            You do not have permission to access this page.
          </p>
          <button
            onClick={() => navigate('/dashboard')}
            className="bg-blue-500 text-white px-5 py-2 rounded hover:bg-blue-600 transition duration-200"
          >
            Back to Dashboard
          </button>
        </div>
      </div>
    </div>
  );
}

export default Unauthorized;
