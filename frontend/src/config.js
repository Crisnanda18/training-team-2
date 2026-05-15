// VULNERABILITY #4: Hardcoded API configuration and secrets
// These should be in environment variables!

export const API_BASE_URL = 'http://localhost:8080/api';

// VULNERABILITY #4: Hardcoded API key in source code
export const ADMIN_API_KEY = 'admin-key-12345'; // Should NEVER be in frontend code!

// VULNERABILITY #4: Hardcoded credentials
export const DEFAULT_CREDENTIALS = {
  admin: {
    email: 'admin@example.com',
    password: 'admin123'
  },
  user: {
    email: 'user@example.com', 
    password: 'password123'
  }
};

// VULNERABILITY: Debug mode left enabled
export const DEBUG_MODE = true;

export const APP_CONFIG = {
  name: 'SecureTask',
  version: '1.0.0',
  // VULNERABILITY #4: AWS credentials hardcoded (even if fake)
  aws: {
    accessKey: 'AKIAIOSFODNN7EXAMPLE',
    secretKey: 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY',
    bucket: 'securetask-uploads'
  }
};
