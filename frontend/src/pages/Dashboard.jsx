import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { getTasks, createTask, updateTask, deleteTask, searchTasks } from '../services/api';
import { getUserData, removeToken, clearUserData } from '../utils/storage';

function Dashboard() {
  const [tasks, setTasks] = useState([]);
  const [searchResults, setSearchResults] = useState(null);
  const [newTask, setNewTask] = useState({ title: '', description: '', priority: 'medium' });
  const [searchTerm, setSearchTerm] = useState('');
  const [user, setUser] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const userData = getUserData();
    setUser(userData);
    loadTasks();
  }, []);

  const loadTasks = async () => {
    try {
      const response = await getTasks();
      setTasks(response.data);
    } catch (error) {
      console.error('Failed to load tasks:', error);
    }
  };

  const handleCreateTask = async (e) => {
    e.preventDefault();
    try {
      await createTask(newTask);
      setNewTask({ title: '', description: '', priority: 'medium' });
      loadTasks();
    } catch (error) {
      console.error('Failed to create task:', error);
    }
  };

  const handleDeleteTask = async (id) => {
    try {
      await deleteTask(id);
      loadTasks();
    } catch (error) {
      console.error('Failed to delete task:', error);
    }
  };

  // VULNERABILITY #1: No sanitization before sending search query (SQL Injection on backend)
  const handleSearch = async (e) => {
    e.preventDefault();
    try {
      const response = await searchTasks(searchTerm);
      setSearchResults(response.data);
    } catch (error) {
      console.error('Search failed:', error);
      alert('Search failed: ' + (error.response?.data?.error || 'Unknown error'));
    }
  };

  const handleLogout = () => {
    removeToken();
    clearUserData();
    navigate('/login');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation */}
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">SecureTask</h1>
              <span className="ml-4 text-sm text-red-600">⚠️ Vulnerable Training App</span>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                {/* VULNERABILITY #5: Displaying sensitive user data from localStorage */}
                Welcome, {user?.name} ({user?.email})
              </span>
              <Link to="/profile" className="text-blue-500 hover:underline">
                Profile
              </Link>
              {user?.role === 'admin' && (
                <Link to="/admin" className="text-blue-500 hover:underline">
                  Admin Panel
                </Link>
              )}
              <button
                onClick={handleLogout}
                className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Search Box */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">Search Tasks</h2>
          <form onSubmit={handleSearch} className="flex gap-2">
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search tasks... (Try: ' OR '1'='1)"
              className="flex-1 px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
            />
            <button
              type="submit"
              className="bg-blue-500 text-white px-6 py-2 rounded hover:bg-blue-600"
            >
              Search
            </button>
          </form>
          
          {searchResults && (
            <div className="mt-4 p-4 bg-gray-50 rounded">
              <h3 className="font-semibold mb-2">Search Results:</h3>
              <pre className="text-sm overflow-auto">
                {JSON.stringify(searchResults, null, 2)}
              </pre>
              <button
                onClick={() => setSearchResults(null)}
                className="mt-2 text-blue-500 hover:underline text-sm"
              >
                Clear Results
              </button>
            </div>
          )}
        </div>

        {/* Create Task Form */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-lg font-semibold mb-4">Create New Task</h2>
          <form onSubmit={handleCreateTask} className="space-y-4">
            <div>
              <input
                type="text"
                value={newTask.title}
                onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
                placeholder="Task title"
                className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              {/* VULNERABILITY #3: No sanitization - XSS possible */}
              <textarea
                value={newTask.description}
                onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
                placeholder="Task description (Try: <script>alert('XSS')</script>)"
                className="w-full px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
                rows="3"
              />
            </div>
            <div className="flex gap-4 items-center">
              <select
                value={newTask.priority}
                onChange={(e) => setNewTask({ ...newTask, priority: e.target.value })}
                className="px-4 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
              >
                <option value="low">Low Priority</option>
                <option value="medium">Medium Priority</option>
                <option value="high">High Priority</option>
              </select>
              <button
                type="submit"
                className="bg-green-500 text-white px-6 py-2 rounded hover:bg-green-600"
              >
                Create Task
              </button>
            </div>
          </form>
        </div>

        {/* Task List */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold mb-4">My Tasks</h2>
          <div className="space-y-4">
            {tasks.length === 0 ? (
              <p className="text-gray-500">No tasks yet. Create one above!</p>
            ) : (
              tasks.map((task) => (
                <div
                  key={task.id}
                  className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition"
                >
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-gray-900">{task.title}</h3>
                      {/* VULNERABILITY #3: Rendering unsanitized HTML - XSS attack vector! */}
                      <div 
                        className="text-gray-600 mt-2"
                        dangerouslySetInnerHTML={{ __html: task.description }}
                      />
                      <div className="mt-2 flex gap-2">
                        <span className={`text-xs px-2 py-1 rounded ${
                          task.priority === 'high' ? 'bg-red-100 text-red-800' :
                          task.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' :
                          'bg-green-100 text-green-800'
                        }`}>
                          {task.priority}
                        </span>
                        <span className="text-xs px-2 py-1 rounded bg-gray-100 text-gray-800">
                          {task.status}
                        </span>
                      </div>
                    </div>
                    <button
                      onClick={() => handleDeleteTask(task.id)}
                      className="ml-4 bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
