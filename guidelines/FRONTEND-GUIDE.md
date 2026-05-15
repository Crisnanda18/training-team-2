# Frontend Developer Security Guide

## 🎯 Your Mission

As a **Frontend Developer**, your task is to identify and fix security vulnerabilities in the SecureTask React application. The frontend contains intentional security flaws that need to be discovered and remediated.

---

## 📋 Discovery Checklist

Work through each area systematically:

### 1. XSS (Cross-Site Scripting) Vulnerabilities

**Files to Review:**
- `frontend/src/pages/Dashboard.jsx`
- `frontend/src/pages/Profile.jsx`

**Questions to Ask:**
- How are task descriptions rendered in the UI?
- Is user-generated content sanitized before display?
- What happens if I create a task with `<script>alert('XSS')</script>` in the description?
- Are we using `dangerouslySetInnerHTML` anywhere?

**What to Look For:**
- Use of `dangerouslySetInnerHTML` without sanitization
- User content rendered directly in HTML
- No output encoding on user-generated fields

**Test Payloads:**
```html
<script>alert('XSS')</script>
<img src=x onerror=alert('XSS')>
<svg onload=alert(document.cookie)>
```

**Impact:**
- Attackers can execute JavaScript in victim's browser
- Can steal JWT tokens from localStorage
- Can perform actions on behalf of users
- Can redirect users to malicious sites

---

### 2. Insecure Data Storage

**Files to Review:**
- `frontend/src/utils/storage.js`
- `frontend/src/pages/Login.jsx`
- `frontend/src/services/api.js`

**Questions to Ask:**
- Where are JWT tokens stored?
- Can JavaScript access the authentication token?
- What user data is stored in localStorage/sessionStorage?
- Are passwords stored anywhere in the frontend?

**What to Look For:**
- JWT tokens in localStorage (vulnerable to XSS)
- User passwords stored in localStorage
- Sensitive data in sessionStorage
- Debug information in localStorage

**Impact:**
- XSS attacks can steal tokens and impersonate users
- Sensitive data accessible to malicious scripts
- Session hijacking
- Privacy violations

---

### 3. Hardcoded Credentials & Secrets

**Files to Review:**
- `frontend/src/config.js`
- `frontend/src/services/api.js`

**Questions to Ask:**
- Are there any API keys in the source code?
- Are there hardcoded credentials?
- Can I see secrets in the browser DevTools?
- Are AWS/cloud credentials exposed?

**What to Look For:**
- API keys in `config.js`
- Hardcoded admin credentials
- AWS access keys and secret keys
- Database connection strings

**Impact:**
- Attackers can use exposed API keys
- Can access services with leaked credentials
- Complete system compromise possible

---

### 4. Missing Content Security Policy (CSP)

**Files to Review:**
- `frontend/index.html`
- `frontend/vite.config.js`

**Questions to Ask:**
- Are there CSP meta tags in the HTML?
- Are CSP headers configured in the server?
- Can inline scripts execute?
- Are external scripts restricted?

**What to Look For:**
- No CSP meta tag in `index.html`
- No CSP headers in server configuration
- No restrictions on script sources

**Impact:**
- No protection against XSS attacks
- Malicious scripts can execute freely
- No defense in depth

---

### 5. Authorization Bypass

**Files to Review:**
- `frontend/src/App.jsx`
- `frontend/src/pages/AdminPanel.jsx`

**Questions to Ask:**
- How is authentication checked?
- Can I access admin pages by modifying localStorage?
- Is authorization checked on the server?
- What happens if I manually navigate to `/admin`?

**What to Look For:**
- Client-side only authorization checks
- Role verification using localStorage data
- No server-side authorization validation

**Impact:**
- Users can access admin functions
- Can bypass access controls
- Privilege escalation

---

### 6. Information Disclosure

**Files to Review:**
- `frontend/src/services/api.js`
- `frontend/src/pages/Profile.jsx`
- All page components

**Questions to Ask:**
- What information is logged to console?
- Are there debug panels showing sensitive data?
- Are API responses logged?
- Is sensitive data visible in the UI?

**What to Look For:**
- `console.log()` statements with sensitive data
- Debug information panels
- Full error messages displayed to users
- Passwords or tokens visible in UI

**Impact:**
- Information leakage to attackers
- Helps attackers understand system internals
- Exposes sensitive user data

---

## 🔍 Tools You Can Use

### Testing XSS
1. Open DevTools Console
2. Create a task with: `<script>alert('XSS')</script>`
3. Update profile bio with: `<img src=x onerror=alert(document.cookie)>`
4. Check if scripts execute

### Testing Storage Security
1. Open DevTools → Application → Local Storage
2. Check what's stored after login
3. Look for JWT tokens and user data
4. Check if passwords are visible

### Testing Authorization Bypass
1. Login as regular user
2. Open DevTools → Application → Local Storage
3. Modify `user` object: change role from `user` to `admin`
4. Navigate to `/admin` - can you access it?

### Checking for Hardcoded Secrets
```bash
# Search for API keys and credentials
grep -r "API_KEY" frontend/src/
grep -r "password" frontend/src/
grep -r "secret" frontend/src/
```

---

## 📝 Expected Outcomes (No Direct Answers!)

After fixing vulnerabilities, your code should achieve:

### ✅ XSS Prevention
- All user-generated content sanitized before rendering
- No use of `dangerouslySetInnerHTML` OR sanitized with DOMPurify
- HTML entities encoded in output
- CSP headers configured

### ✅ Secure Data Storage
- JWT tokens stored in httpOnly cookies (requires backend change)
- No sensitive data in localStorage
- User passwords never stored in frontend
- Debug information removed

### ✅ Secrets Management
- No hardcoded API keys in source code
- Secrets moved to environment variables
- Use `.env` files (not committed)
- Access via `import.meta.env.VITE_*`

### ✅ Content Security Policy
- CSP meta tag added to `index.html`
- Restricts script sources
- Prevents inline script execution
- Reports CSP violations

### ✅ Proper Authorization
- Server-side authorization checks (work with backend team)
- Client-side checks for UX only, not security
- Proper error handling when access denied

### ✅ No Information Disclosure
- Remove all `console.log()` with sensitive data
- Hide debug panels in production
- Generic error messages to users
- Detailed errors only in development mode

---

## 🎓 Learning Resources

### XSS Prevention
- OWASP XSS Prevention Cheat Sheet
- DOMPurify library: https://github.com/cure53/DOMPurify
- React security best practices

### Secure Storage
- OWASP Session Management
- httpOnly cookies vs localStorage
- Token storage best practices

### Content Security Policy
- MDN CSP documentation
- CSP header generator tools
- CSP reporting

### Environment Variables in Vite
- Vite environment variables: https://vitejs.dev/guide/env-and-mode.html
- Using `import.meta.env.VITE_*`

---

## ✋ Hints (Not Solutions!)

### For XSS Prevention:
- Install DOMPurify: `npm install dompurify`
- Replace `dangerouslySetInnerHTML={{ __html: content }}` with sanitized version
- Example pattern: `dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(content) }}`
- Or better: avoid rendering HTML, use plain text

### For Secure Storage:
- Remove localStorage usage for tokens
- Work with backend team to implement httpOnly cookies
- Clear all sensitive data from storage utilities
- Use session storage only for non-sensitive UI state

### For Hardcoded Secrets:
- Create `.env` file with: `VITE_API_URL=http://localhost:8080/api`
- Access with: `import.meta.env.VITE_API_URL`
- Add `.env` to `.gitignore`
- Never commit actual secrets

### For CSP:
```html
<!-- Add to index.html <head> -->
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline';">
```

### For Authorization:
- Keep client-side checks for UX (showing/hiding buttons)
- Backend must ALWAYS validate authorization
- Handle 403 errors gracefully
- Don't trust localStorage for security decisions

### For Info Disclosure:
- Remove `console.log()` statements or wrap with `if (import.meta.env.DEV)`
- Remove debug panels showing user objects
- Use generic error messages: "Operation failed" instead of detailed SQL errors

---

## 🧪 How to Verify Your Fixes

### Test XSS is Fixed:
```
1. Create task with: <script>alert('XSS')</script>
2. Script should NOT execute
3. Should see escaped HTML or plain text
```

### Test Secure Storage:
```
1. Login to application
2. Check localStorage - should NOT contain JWT token
3. Check that authentication still works (using cookies)
```

### Test No Hardcoded Secrets:
```bash
# Should find no secrets
grep -r "admin-key" frontend/src/
grep -r "AWS" frontend/src/
```

### Test CSP:
```
1. Open DevTools Console
2. Try: eval('alert("test")')
3. Should be blocked by CSP
```

### Test Authorization:
```
1. Login as regular user
2. Try accessing /admin
3. Backend should return 403/401
4. Even if localStorage modified
```

---

## 🤝 Collaboration Points

- **With Backend Team**: 
  - Request httpOnly cookie implementation for tokens
  - Ensure backend validates all authorization
  - Discuss CSP headers configuration
  
- **With QA Team**: 
  - Provide XSS test payloads
  - Show how to check localStorage
  - Explain authorization bypass technique

---

## 🏆 Success Criteria

You've successfully completed the frontend security fixes when:

1. ✅ XSS payloads don't execute (sanitized)
2. ✅ JWT tokens not in localStorage
3. ✅ No hardcoded secrets in source code
4. ✅ CSP configured and blocking inline scripts
5. ✅ Admin panel requires server-side auth
6. ✅ No sensitive data logged to console
7. ✅ No debug information visible in production
8. ✅ QA team cannot exploit any frontend vulnerabilities

**Remember**: Security is a shared responsibility. Work closely with the backend team to ensure both client and server are secure!
