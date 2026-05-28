// FIXED: Hardcoded secrets removed
// API URL diambil dari environment variable (.env file)
// Vite replace import.meta.env.VITE_* saat build time

export const API_BASE_URL = import.meta.env.VITE_API_URL || '/api';

// REMOVED: ADMIN_API_KEY
// Alasan: API key statis di frontend bisa dilihat siapapun lewat DevTools → Sources.
// Admin authorization seharusnya lewat JWT token (sudah ada di interceptor).

// REMOVED: DEFAULT_CREDENTIALS
// Alasan: Hardcoded credentials di source code = publik. Siapapun yang buka
// bundled JS bisa baca email/password. Test accounts didokumentasikan di README saja.

// REMOVED: DEBUG_MODE
// Alasan: Gunakan import.meta.env.DEV bawaan Vite.
// Otomatis true di development, false di production build.

// REMOVED: AWS credentials
// Alasan: Cloud credentials TIDAK PERNAH boleh di frontend code.
// Kalau butuh upload ke S3, gunakan presigned URL dari backend.

// Non-sensitive app config
export const APP_CONFIG = {
  name: 'SecureTask',
  version: '1.0.0',
};
