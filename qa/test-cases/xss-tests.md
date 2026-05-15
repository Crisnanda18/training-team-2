# XSS (Cross-Site Scripting) Test Cases

## Test Case 1: Stored XSS in Task Description
**Test ID**: XSS-001  
**Priority**: High  
**Category**: XSS - Stored  

**Objective**: Verify that task descriptions are vulnerable to stored XSS

**Preconditions**:
- User is logged in
- User can create tasks

**Test Steps**:
1. Navigate to Dashboard
2. Click "Create Task"
3. Title: "Test XSS"
4. Description: `<script>alert('XSS')</script>`
5. Click "Create Task"
6. Observe the dashboard

**Expected Result (Vulnerable)**:
- Alert box appears
- JavaScript executes

**Expected Result (Secure)**:
- Script tags displayed as text or removed
- No JavaScript execution

**Payload**: `<script>alert('XSS')</script>`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 2: Image Tag XSS
**Test ID**: XSS-002  
**Priority**: High  

**Objective**: Test XSS using image tag with onerror

**Test Steps**:
1. Create task with description: `<img src=x onerror=alert('XSS')>`
2. View task on dashboard

**Expected Result (Vulnerable)**:
- Alert executes when image fails to load

**Expected Result (Secure)**:
- HTML is escaped or sanitized
- No alert

**Payload**: `<img src=x onerror=alert('XSS')>`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 3: SVG XSS
**Test ID**: XSS-003  
**Priority**: High  

**Objective**: Test XSS using SVG tag

**Test Steps**:
1. Create task with: `<svg onload=alert('XSS')></svg>`
2. View task

**Expected Result (Vulnerable)**:
- Alert executes on page load

**Expected Result (Secure)**:
- SVG tag sanitized or escaped

**Payload**: `<svg onload=alert('XSS')></svg>`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 4: XSS in Profile Bio
**Test ID**: XSS-004  
**Priority**: High  

**Objective**: Test XSS in user profile bio field

**Test Steps**:
1. Navigate to Profile page
2. Update bio with: `<img src=x onerror=alert(document.cookie)>`
3. Save profile
4. Refresh page

**Expected Result (Vulnerable)**:
- Alert shows cookie/localStorage data
- Proves token theft possibility

**Expected Result (Secure)**:
- HTML escaped
- No script execution

**Payload**: `<img src=x onerror=alert(document.cookie)>`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 5: Event Handler XSS
**Test ID**: XSS-005  
**Priority**: Medium  

**Objective**: Test various HTML event handlers

**Test Payloads**:
- `<body onload=alert('XSS')>`
- `<input onfocus=alert('XSS') autofocus>`
- `<marquee onstart=alert('XSS')>`
- `<div onmouseover=alert('XSS')>Hover me</div>`

**Test Steps**:
1. Try each payload in task description
2. Observe behavior

**Expected Result (Vulnerable)**:
- Events trigger JavaScript

**Expected Result (Secure)**:
- All sanitized

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 6: JavaScript Protocol XSS
**Test ID**: XSS-006  
**Priority**: Medium  

**Objective**: Test javascript: protocol in href

**Test Steps**:
1. Create task with: `<a href="javascript:alert('XSS')">Click</a>`
2. Click the link

**Expected Result (Vulnerable)**:
- Alert executes on click

**Expected Result (Secure)**:
- Link is sanitized or disabled

**Payload**: `<a href="javascript:alert('XSS')">Click</a>`

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 7: Data URI XSS
**Test ID**: XSS-007  
**Priority**: Medium  

**Objective**: Test data URI for XSS

**Payload**: `<a href="data:text/html,<script>alert('XSS')</script>">Click</a>`

**Test Steps**:
1. Insert payload in description
2. Click link if rendered

**Expected Result (Vulnerable)**:
- XSS executes in new context

**Expected Result (Secure)**:
- Data URIs blocked or sanitized

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## Test Case 8: DOM-Based XSS
**Test ID**: XSS-008  
**Priority**: High  

**Objective**: Check if URL parameters are reflected unsafely

**Test Steps**:
1. Try accessing: `http://localhost:5173/?search=<script>alert('XSS')</script>`
2. Check if parameter is rendered without sanitization

**Expected Result (Vulnerable)**:
- Script from URL executes

**Expected Result (Secure)**:
- URL parameters sanitized

**Status**: [ ] Pass [ ] Fail [ ] Blocked

---

## XSS Test Automation Script

```javascript
// Run in browser console
const xssPayloads = [
  '<script>alert("XSS")</script>',
  '<img src=x onerror=alert("XSS")>',
  '<svg onload=alert("XSS")>',
  '<iframe src="javascript:alert(\'XSS\')">',
  '<body onload=alert("XSS")>'
];

console.log('Testing XSS payloads...');
xssPayloads.forEach((payload, index) => {
  console.log(`Payload ${index + 1}: ${payload}`);
  // Instructions: Manually test each payload
});
```
