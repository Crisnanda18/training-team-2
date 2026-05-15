# Backend Developer Security Guide

## 🎯 Your Mission

As a **Backend Developer**, your task is to identify and fix security vulnerabilities in the SecureTask API. The backend contains intentional security flaws that need to be discovered and remediated.

---

## 📋 Discovery Checklist

Work through each area systematically:

### 1. SQL Injection Vulnerabilities

**Files to Review:**
- `backend/database/database.go`
- `backend/handlers/tasks.go`

**Questions to Ask:**
- How is the search functionality implemented?
- Are database queries using parameterized statements?
- Is user input being directly concatenated into SQL queries?
- What happens if I search for: `' OR '1'='1` ?

**What to Look For:**
- String concatenation in SQL queries
- Use of `fmt.Sprintf()` with user input in SQL
- Raw SQL execution without parameter binding

**Impact:**
- Attackers can read all data from database
- Can modify or delete data
- Can potentially execute commands on database server

---

### 2. Authentication & Authorization Issues

**Files to Review:**
- `backend/main.go` (route definitions)
- `backend/handlers/auth.go`
- `backend/handlers/users.go`

**Questions to Ask:**
- Which endpoints require authentication middleware?
- Are there routes exposed publicly that should be protected?
- Can users access resources they don't own?
- Is there proper role-based access control for admin endpoints?

**What to Look For:**
- Routes without authentication middleware
- Missing authorization checks (e.g., checking if user owns the resource)
- Admin endpoints accessible without role verification
- No rate limiting on sensitive endpoints (login)

**Impact:**
- Unauthorized access to sensitive data
- Ability to modify other users' data
- Privilege escalation attacks
- Brute force attacks on login

---

### 3. Input Validation & Sanitization

**Files to Review:**
- `backend/handlers/tasks.go`
- `backend/handlers/users.go`
- `backend/models/user.go`
- `backend/models/task.go`

**Questions to Ask:**
- Is user input validated before processing?
- Are there minimum/maximum length checks?
- Is data sanitized before storage?
- What fields accept HTML/script content?

**What to Look For:**
- No input validation on API requests
- Missing data sanitization
- Fields that could contain XSS payloads (description, bio)
- No password complexity requirements

**Impact:**
- XSS attacks via stored malicious content
- Data integrity issues
- Weak password acceptance

---

### 4. Hardcoded Credentials & Secrets

**Files to Review:**
- `backend/main.go`
- `backend/database/database.go`
- `backend/handlers/auth.go`
- `backend/.env`

**Questions to Ask:**
- Are there hardcoded passwords or secrets in the code?
- Is the JWT secret strong and from environment variables?
- Are database credentials in the source code?
- Is the `.env` file committed to Git?

**What to Look For:**
- Hardcoded database connection strings
- JWT secrets in constants
- API keys in source code
- Weak secrets (e.g., "supersecret123")

**Impact:**
- Complete system compromise
- Unauthorized access to all data
- Ability to forge authentication tokens

---

### 5. Sensitive Data Exposure

**Files to Review:**
- `backend/models/user.go`
- `backend/handlers/auth.go`
- `backend/handlers/users.go`

**Questions to Ask:**
- Are passwords being returned in API responses?
- Is the password field excluded from JSON serialization?
- Are passwords hashed before storage?
- What password hashing algorithm is used (if any)?

**What to Look For:**
- Password field without `json:"-"` tag
- Plain text password storage
- Passwords in API responses
- No use of bcrypt or similar

**Impact:**
- Password leakage
- Account compromise
- Compliance violations (GDPR, etc.)

---

## 🔍 Tools You Can Use

### Testing SQL Injection
```bash
# Test search endpoint
curl "http://localhost:8080/api/tasks/search?q=' OR '1'='1"
curl "http://localhost:8080/api/tasks/search?q=' UNION SELECT * FROM users--"
```

### Testing Authentication
```bash
# Try accessing protected endpoints without auth
curl http://localhost:8080/api/admin/users
curl -X DELETE http://localhost:8080/api/tasks/1
```

### Reviewing Database Queries
```bash
# Search for SQL string concatenation
grep -r "fmt.Sprintf" backend/
grep -r "SELECT \* FROM" backend/
```

---

## 📝 Expected Outcomes (No Direct Answers!)

After fixing vulnerabilities, your code should achieve:

### ✅ SQL Injection Prevention
- All database queries use parameterized statements
- User input is never directly concatenated into SQL
- GORM methods are used correctly (e.g., `Where("field = ?", value)`)

### ✅ Proper Authentication & Authorization
- All sensitive endpoints protected with authentication middleware
- Authorization checks verify user owns the resource
- Admin endpoints verify role = "admin"
- Rate limiting implemented on login endpoint

### ✅ Input Validation
- All API inputs validated for type, length, format
- Password strength requirements enforced
- Appropriate error messages without revealing sensitive info

### ✅ Secrets Management
- All secrets moved to environment variables
- Strong, randomly generated JWT secret
- `.env` file NOT committed to Git (add to `.gitignore`)
- Database credentials loaded from environment

### ✅ Secure Password Handling
- Passwords hashed using bcrypt before storage
- Password field excluded from JSON responses (`json:"-"` tag)
- Password comparison uses bcrypt.CompareHashAndPassword()
- Minimum password length enforced (e.g., 8 characters)

---

## 🎓 Learning Resources

### SQL Injection Prevention
- OWASP SQL Injection: https://owasp.org/www-community/attacks/SQL_Injection
- Parameterized queries in Go
- GORM security best practices

### Authentication & Authorization
- JWT best practices
- RBAC (Role-Based Access Control)
- OWASP Authentication Cheat Sheet

### Password Security
- Bcrypt documentation for Go: `golang.org/x/crypto/bcrypt`
- OWASP Password Storage Cheat Sheet
- Password hashing algorithms

### Environment Variables in Go
- Using `os.Getenv()` in Go
- Loading `.env` files with `godotenv` package

---

## ✋ Hints (Not Solutions!)

### For SQL Injection:
- Look at the `SearchTasks` function - how is the query built?
- Replace string concatenation with `DB.Where("field LIKE ?", "%"+value+"%")`
- Remove the `ExecuteRawSQL` function or ensure it's never called with user input

### For Authentication:
- Check `main.go` route definitions - which routes lack `authorized.Use()`?
- Add middleware to public routes that should be protected
- In `GetAllUsers`, add a check: `if role != "admin" { return error }`

### For Hardcoded Secrets:
- Use `os.Getenv("JWT_SECRET")` instead of hardcoded strings
- Load `.env` using `godotenv` package
- Never commit `.env` files to Git

### For Password Security:
- Import `golang.org/x/crypto/bcrypt`
- Hash on register: `hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)`
- Compare on login: `err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))`
- Add `json:"-"` tag to Password field in User model

---

## 📊 How to Verify Your Fixes

### Test SQL Injection is Fixed:
```bash
# This should NOT return all tasks
curl "http://localhost:8080/api/tasks/search?q=' OR '1'='1"
```

### Test Authentication is Required:
```bash
# These should return 401 Unauthorized
curl http://localhost:8080/api/admin/users
curl -X DELETE http://localhost:8080/api/tasks/1
```

### Test Password Hashing:
```bash
# Register a user and check the database
# Password should be hashed, not plain text
```

### Test Secrets Not in Code:
```bash
# Search for hardcoded secrets
grep -r "supersecret" backend/
grep -r "taskpass123" backend/
# Should only find them in .env file
```

---

## 🤝 Collaboration Points

- **With Frontend Team**: Ensure API returns appropriate error messages without exposing internal details
- **With QA Team**: Provide them with test credentials and explain authentication flow
- **Cross-team**: Discuss how to handle XSS - backend should not trust frontend sanitization

---

## 🏆 Success Criteria

You've successfully completed the backend security fixes when:

1. ✅ No SQL injection possible in any endpoint
2. ✅ All sensitive endpoints require authentication
3. ✅ Admin endpoints check user role
4. ✅ Passwords are hashed with bcrypt
5. ✅ No hardcoded secrets in source code
6. ✅ Password fields not returned in API responses
7. ✅ Input validation on all endpoints
8. ✅ QA team cannot exploit any backend vulnerabilities

**Remember**: Don't just fix - understand WHY each vulnerability is dangerous and HOW your fix prevents exploitation!
