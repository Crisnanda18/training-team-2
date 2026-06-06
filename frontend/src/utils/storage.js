// VULNERABILITY #5: Insecure storage of sensitive data

let tokenMemory = null;

// VULNERABILITY: Storing JWT token in localStorage (vulnerable to XSS)
export const setToken = (token) => {
  {/* 
    localStorage.setItem('token', token); // Should use httpOnly cookies!
    */}
  
  // document.cookie = `token=${token}; path=/; secure; samesite=strict;`;
  tokenMemory = token;
};

export const getToken = () => {
  // return localStorage.getItem('token');

  // return document.cookie.split('; ').find(row => row.startsWith('token='))?.split('=')[1] || null;
  return tokenMemory;
};

export const removeToken = () => {
  // document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1945 17:08:45 GMT';
  tokenMemory = null;
};

{/*
  // VULNERABILITY #5: Storing user data including sensitive info in localStorage
  export const setUserData = (user) => {
  
    const { password, ...sanitizedUser } = user;
    localStorage.setItem('user', JSON.stringify(sanitizedUser));
    
    sessionStorage.setItem('currentUser', JSON.stringify(sanitizedUser));
  };
*/}

// tdk dipakai (ngga dihapus untuk penjelasan)
// export const getUserData = () => {
//   const user = localStorage.getItem('user');
//   return user ? JSON.parse(user) : null;
// };

// tdk dipakai (ngga dihapus untuk penjelasan)
// export const clearUserData = () => {
//   // VULNERABILITY: Not clearing all sensitive data
//   // localStorage.clear() would be better, but this leaves traces
  
//   // document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1945 17:08:45 GMT';
//   localStorage.removeItem('user');
//   sessionStorage.removeItem('currentUser');
// };

// VULNERABILITY #5: Storing sensitive settings in localStorage (tidak dipakai)
// export const saveSettings = (settings) => {
//   localStorage.setItem('appSettings', JSON.stringify(settings));
//   // document.cookie
// };

{/*
  // VULNERABILITY: Exposing internal debug data
export const saveDebugInfo = (info) => {
  localStorage.setItem('debugInfo', JSON.stringify({
    ...info,
    timestamp: new Date().toISOString(),
    userAgent: navigator.userAgent
  }));
};
   */}
