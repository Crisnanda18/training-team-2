// VULNERABILITY #4: Hardcoded API configuration and secrets
// These should be in environment variables!

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL;

// VULNERABILITY #4: Hardcoded API key in source code
export const ADMIN_API_KEY = import.meta.env.VITE_ADMIN_API_KEY; // Should NEVER be in frontend code!


// VULNERABILITY: Debug mode left enabled
export const VITE_DEBUG_MODE = import.meta.env.VITE_DEBUG_MODE === 'true'; // Should be false in production

export const APP_CONFIG = {
  name: 'SecureTask',
  version: '1.0.0',
  // VULNERABILITY #4: AWS credentials hardcoded (even if fake)
  aws: {
    accessKey: import.meta.env.VITE_AWS_ACCESS_KEY,
    secretKey: import.meta.env.VITE_AWS_SECRET_KEY,
    bucket: import.meta.env.VITE_AWS_BUCKET_NAME
  }
};
