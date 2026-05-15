# SecureTask - Security Training Project 🔐

A deliberately vulnerable task management application designed for security training. This project contains **intentional security vulnerabilities** based on OWASP Top 10 that need to be identified and fixed by development teams.

## ⚠️ WARNING
This application contains intentional security vulnerabilities. **DO NOT deploy this to production or any public environment.** It is designed solely for educational purposes.

## 🎯 Training Objectives

Learn to identify and fix 5 critical security vulnerabilities:

1. **SQL Injection** - Unsafe database queries
2. **Unsecured API Design** - Missing authentication/authorization, lack of input validation
3. **XSS (Cross-Site Scripting)** - Unsanitized user input in HTML rendering
4. **Hardcoded Credentials** - Secrets exposed in code and Git commits
5. **Unsecured Data Storage** - Sensitive data stored insecurely

## 🏗️ Project Structure

```
securetask/
├── backend/          # Golang API (contains vulnerabilities)
├── frontend/         # React application (contains vulnerabilities)
├── qa/              # QA testing materials and templates
├── guidelines/      # Role-specific guides for each team
└── README.md        # This file
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- Git

### 1. Clone and Setup Database

```bash
# Start PostgreSQL database
docker-compose up -d

# Database will be available at localhost:5432
# Database: securetask
# User: taskuser
# Password: taskpass123
```

### 2. Start Backend (Golang)

```bash
cd backend
go mod download
go run main.go

# Server runs at http://localhost:8080
```

### 3. Start Frontend (React)

```bash
cd frontend
npm install
npm run dev

# Application runs at http://localhost:5173
```

### 4. Access the Application

Open your browser and navigate to `http://localhost:5173`

**Default Accounts:**
- Regular User: `user@example.com` / `password123`
- Admin User: `admin@example.com` / `admin123`

## 👥 Team Roles & Guides

### Backend Developers
📖 Read: `guidelines/BACKEND-GUIDE.md`

Your mission: Identify and fix backend security vulnerabilities including SQL injection, authentication issues, and hardcoded credentials.

### Frontend Developers
📖 Read: `guidelines/FRONTEND-GUIDE.md`

Your mission: Identify and fix frontend security vulnerabilities including XSS attacks, insecure storage, and exposed secrets.

### QA Team
📖 Read: `guidelines/QA-GUIDE.md`

Your mission: Test the application for vulnerabilities, document findings, verify fixes, and ensure security standards are met.

## 📚 Features (All Vulnerable!)

- ✅ User Authentication (Login/Register)
- ✅ Task Management (Create, Read, Update, Delete)
- ✅ Search Tasks
- ✅ User Profile Management
- ✅ Admin Panel
- ❌ All features contain security vulnerabilities

## 🎓 Training Workflow

### Phase 1: Discovery (Week 1)
- Each team explores the codebase
- Identify vulnerabilities in your domain
- Document findings

### Phase 2: Analysis (Week 1-2)
- Understand the impact of each vulnerability
- Research best practices
- Plan fixes

### Phase 3: Implementation (Week 2-3)
- Backend team fixes API vulnerabilities
- Frontend team fixes client-side vulnerabilities
- QA team creates test cases

### Phase 4: Verification (Week 3-4)
- Cross-team code reviews
- QA verifies all fixes
- Final security assessment

## 🛠️ Tech Stack

**Backend:**
- Go 1.21
- Gin Web Framework
- GORM (ORM)
- PostgreSQL
- JWT Authentication

**Frontend:**
- React 18
- Vite
- React Router
- Axios
- Tailwind CSS

## 📖 Documentation

- `guidelines/BACKEND-GUIDE.md` - Backend security guide
- `guidelines/FRONTEND-GUIDE.md` - Frontend security guide
- `guidelines/QA-GUIDE.md` - QA testing guide
- `guidelines/EXPECTED-OUTCOMES.md` - What success looks like
- `guidelines/TRAINER-SOLUTIONS.md` - Complete solutions (for trainers only)

## 🔍 What Makes This Project Effective?

- ✅ **Realistic**: Based on real-world vulnerabilities
- ✅ **Hands-on**: Learn by fixing actual code
- ✅ **Collaborative**: Requires team coordination
- ✅ **Comprehensive**: Covers OWASP Top 10 issues
- ✅ **Measurable**: Clear before/after states

## 📞 Support

For questions or issues with setup, contact your training coordinator.

## 📄 License

This project is for educational purposes only.

---

**Remember**: Never deploy this application to production. It contains intentional security vulnerabilities for learning purposes.
