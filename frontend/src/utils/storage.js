// VULNERABILITY #5: Insecure storage of sensitive data

// VULNERABILITY: Storing JWT token in localStorage (vulnerable to XSS)
export const setToken = (token) => {
  localStorage.setItem('token', token); // Should use httpOnly cookies!
};

export const getToken = () => {
  return localStorage.getItem('token');
};

export const removeToken = () => {
  localStorage.removeItem('token');
};

// VULNERABILITY #5: Storing user data including sensitive info in localStorage
export const setUserData = (user) => {
  // Storing full user object including password!
  localStorage.setItem('user', JSON.stringify(user));
  
  // VULNERABILITY: Also storing in sessionStorage
  sessionStorage.setItem('currentUser', JSON.stringify(user));
};

export const getUserData = () => {
  const user = localStorage.getItem('user');
  return user ? JSON.parse(user) : null;
};

export const clearUserData = () => {
  localStorage.removeItem('user');
  sessionStorage.removeItem('currentUser');
  
  // VULNERABILITY: Not clearing all sensitive data
  // localStorage.clear() would be better, but this leaves traces
};

// VULNERABILITY #5: Storing sensitive settings in localStorage
export const saveSettings = (settings) => {
  localStorage.setItem('appSettings', JSON.stringify(settings));
};

// VULNERABILITY: Exposing internal debug data
export const saveDebugInfo = (info) => {
  localStorage.setItem('debugInfo', JSON.stringify({
    ...info,
    timestamp: new Date().toISOString(),
    userAgent: navigator.userAgent
  }));
};
