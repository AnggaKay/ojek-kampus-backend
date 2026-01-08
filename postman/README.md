# 📬 Postman Testing Guide

## 🚀 Quick Start

### 1. Import ke Postman
1. Buka Postman
2. Click **Import** (top left)
3. Drag & drop kedua file:
   - `Ojek_Kampus_API.postman_collection.json`
   - `Ojek_Kampus_Local.postman_environment.json`
4. Select environment: **Ojek Kampus - Local** (top right)

### 2. Start Server
```bash
go run cmd/api/main.go
```

Server akan jalan di: `http://localhost:8080`

---

## ✅ Testing Checklist

### A. Health Check
- [ ] **GET** `/health` → Status 200, response success

**Expected Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "ok",
    "timestamp": "2026-01-08T21:33:25+07:00",
    "service": "ojek-kampus-backend"
  }
}
```

---

### B. Happy Path - Register & Login Flow

#### 1. Register Passenger
- [ ] **POST** `/api/auth/register/passenger`
- [ ] Status: **201 Created**
- [ ] Response contains: `user`, `access_token`, `refresh_token`
- [ ] Phone normalized: `081234567890` → `+6281234567890`
- [ ] Token auto-saved ke environment variables

**Request Body:**
```json
{
  "phone_number": "081234567890",
  "password": "password123",
  "full_name": "John Doe",
  "email": "john@example.com"
}
```

**Expected Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "id": 1,
      "phone_number": "+6281234567890",
      "email": "john@example.com",
      "full_name": "John Doe",
      "role": "PASSENGER",
      "status": "ACTIVE",
      "phone_verified": false
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "9f75567beee7cfccac4d...",
    "expires_in": 900
  }
}
```

#### 2. Login
- [ ] **POST** `/api/auth/login`
- [ ] Status: **200 OK**
- [ ] Same response structure as register
- [ ] Token auto-saved ke environment

**Request Body:**
```json
{
  "phone_number": "081234567890",
  "password": "password123"
}
```

#### 3. Get Profile (Protected)
- [ ] **GET** `/api/auth/me`
- [ ] Header: `Authorization: Bearer {{access_token}}`
- [ ] Status: **200 OK**
- [ ] Response contains user info

**Expected Response:**
```json
{
  "success": true,
  "message": "Profile retrieved",
  "data": {
    "user_id": 1,
    "role": "PASSENGER"
  }
}
```

#### 4. Refresh Token
- [ ] **POST** `/api/auth/refresh`
- [ ] Status: **200 OK**
- [ ] Response contains new `access_token`
- [ ] New token auto-saved

**Request Body:**
```json
{
  "refresh_token": "{{refresh_token}}"
}
```

#### 5. Logout
- [ ] **POST** `/api/auth/logout`
- [ ] Status: **200 OK**
- [ ] Refresh token revoked in database

**Request Body:**
```json
{
  "refresh_token": "{{refresh_token}}"
}
```

---

### C. Validation Testing

#### 1. Invalid Phone Format
- [ ] Phone: `"123"` → Status **400**, error message

#### 2. Weak Password
- [ ] Password: `"123"` (< 8 char) → Status **400**
- [ ] Password: `"abcdefgh"` (no number) → Status **400**
- [ ] Password: `"12345678"` (no letter) → Status **400**

#### 3. Missing Required Fields
- [ ] Missing `phone_number` → Status **400**
- [ ] Missing `password` → Status **400**
- [ ] Missing `full_name` → Status **400**

#### 4. Duplicate Phone Number
- [ ] Register dengan phone yang sudah ada → Status **400**
- [ ] Error: `"phone number already registered"`

#### 5. Duplicate Email
- [ ] Register dengan email yang sudah ada → Status **400**
- [ ] Error: `"email already registered"`

---

### D. Authentication Testing

#### 1. Missing Authorization Header
- [ ] **GET** `/api/auth/me` tanpa header → Status **401**
- [ ] Error: `"Missing authorization token"`

#### 2. Invalid Token Format
- [ ] Header: `Authorization: InvalidToken` → Status **401**
- [ ] Error: `"Invalid authorization format"`

#### 3. Invalid Token
- [ ] Header: `Authorization: Bearer fake_token` → Status **401**
- [ ] Error: `"Invalid or expired token"`

#### 4. Wrong Password
- [ ] Login dengan password salah → Status **401**
- [ ] Error: `"invalid phone number or password"`

#### 5. Non-existent User
- [ ] Login dengan phone yang tidak terdaftar → Status **401**

---

### E. Edge Cases

#### 1. Phone Number Normalization
Test berbagai format phone:
- [ ] `081234567890` → `+6281234567890` ✅
- [ ] `6281234567890` → `+6281234567890` ✅
- [ ] `+6281234567890` → `+6281234567890` ✅
- [ ] `0812-3456-7890` → `+6281234567890` ✅

#### 2. Email Optional
- [ ] Register tanpa email → Success ✅
- [ ] Register dengan email → Success ✅

#### 3. Refresh Token Expiry
- [ ] Use expired refresh token → Status **401**
- [ ] Error: `"token has expired"`

#### 4. Revoked Token
- [ ] Logout → Token revoked
- [ ] Try to refresh with revoked token → Status **401**

---

## 📊 Expected Status Codes

| Scenario | Status Code | Success |
|----------|-------------|---------|
| Register success | 201 Created | true |
| Login success | 200 OK | true |
| Profile retrieved | 200 OK | true |
| Token refreshed | 200 OK | true |
| Logout success | 200 OK | true |
| Validation error | 400 Bad Request | false |
| Unauthorized | 401 Unauthorized | false |
| Forbidden | 403 Forbidden | false |
| Not found | 404 Not Found | false |
| Server error | 500 Internal Error | false |

---

## 🔧 Troubleshooting

### Server tidak jalan
```bash
# Check port 8080
netstat -ano | findstr :8080

# Kill process
taskkill /F /PID <PID>

# Run server
go run cmd/api/main.go
```

### Token tidak auto-saved
1. Check tab **Tests** di Postman request
2. Check **Console** (bottom left Postman) untuk log
3. Manually copy token ke environment variables

### 401 Unauthorized terus
1. Check token di environment: `{{access_token}}`
2. Token expire setelah 15 menit → Login lagi atau refresh
3. Check header format: `Bearer <token>` (ada spasi)

---

## 📤 Export untuk Frontend

Setelah testing selesai, berikan ke frontend team:

### 1. API Documentation
```markdown
Base URL: http://localhost:8080

Endpoints:
- POST /api/auth/register/passenger
- POST /api/auth/login
- POST /api/auth/refresh
- POST /api/auth/logout
- GET  /api/auth/me (requires Bearer token)
```

### 2. Example cURL
```bash
# Register
curl -X POST http://localhost:8080/api/auth/register/passenger \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"081234567890","password":"password123","full_name":"John Doe"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"081234567890","password":"password123"}'

# Get Profile
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 3. Response Schema
Berikan file:
- `internal/dto/auth_request.go` - Request schemas
- `internal/dto/auth_response.go` - Response schemas

---

## ✅ Success Criteria

Semua test cases berikut harus **PASS**:
- [x] Register passenger berhasil
- [x] Login berhasil
- [x] Token auto-saved
- [x] Protected route accessible dengan token
- [x] Validation errors handled
- [x] Phone normalization working
- [x] Password hashing (tidak return plain text)
- [x] Refresh token working
- [x] Logout revoke token

**Setelah semua checklist ✅, API siap untuk frontend integration!** 🎉
