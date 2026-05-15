# SecureTask Frontend

React frontend for the SecureTask security training project.

## ⚠️ Security Warning
This frontend contains **intentional security vulnerabilities** for training purposes. DO NOT use in production!

## Setup

### Prerequisites
- Node.js 18+ 
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production (don't actually deploy this!)
npm run build
```

The application will run on `http://localhost:5173`

## Features

- **Login/Register** - User authentication
- **Dashboard** - Task management with create, view, delete
- **Search** - Search tasks (vulnerable to SQL injection on backend)
- **Profile** - Update user profile and bio
- **Admin Panel** - View all users (no proper authorization)

## Known Vulnerabilities (For Training)

### 1. XSS (Cross-Site Scripting)
- Task descriptions rendered with `dangerouslySetInnerHTML`
- User bio field not sanitized
- No output encoding on user-generated content

**Test Payloads:**
```html
<script>alert('XSS')</script>
<img src=x onerror=alert('XSS')>
<svg onload=alert('XSS')>
```

### 2. Insecure Data Storage
- JWT tokens stored in localStorage (vulnerable to XSS)
- User data including passwords stored in localStorage
- Session data stored in sessionStorage
- Debug information exposed in localStorage

### 3. Hardcoded Credentials
- API keys hardcoded in `src/config.js`
- Default credentials exposed
- AWS credentials (example) in source code

### 4. Missing Authorization
- Admin panel accessible with client-side check only
- No server-side role verification
- Can bypass checks by modifying localStorage

### 5. No CSP (Content Security Policy)
- No CSP headers configured
- Inline scripts allowed
- No protection against XSS

### 6. Information Disclosure
- Sensitive data logged to console
- Debug information visible in UI
- Detailed error messages exposed

### 7. SQL Injection (Backend)
- Search functionality sends unsanitized input to backend
- Backend uses string concatenation for SQL queries

## Pages

### Login (`/login`)
- Stores JWT in localStorage
- Saves full user object including password

### Register (`/register`)
- No password strength validation
- No client-side input validation

### Dashboard (`/dashboard`)
- XSS vulnerability in task descriptions
- SQL injection via search
- Displays sensitive user data

### Profile (`/profile`)
- XSS in bio field
- No authorization check for profile updates
- Can update any user's profile

### Admin Panel (`/admin`)
- Client-side authorization only
- Displays all users with plain text passwords
- No server-side access control

## Technology Stack

- React 18
- Vite
- React Router
- Axios
- Tailwind CSS

## Testing Vulnerabilities

### Test XSS:
1. Create a task with description: `<script>alert('XSS')</script>`
2. Update profile bio with: `<img src=x onerror=alert('XSS')>`

### Test SQL Injection:
1. Search for: `' OR '1'='1`
2. Search for: `'; DROP TABLE tasks; --`

### Test Insecure Storage:
1. Login and check localStorage in DevTools
2. See JWT token and full user data including password

### Test Authorization Bypass:
1. Login as regular user
2. Modify localStorage: `role` from `user` to `admin`
3. Access `/admin` endpoint

## What Should Be Fixed?

- Use DOMPurify for HTML sanitization
- Store tokens in httpOnly cookies
- Remove hardcoded credentials
- Implement proper CSP headers
- Add server-side authorization checks
- Never store passwords in frontend
- Validate and sanitize all inputs
- Remove debug information from production
