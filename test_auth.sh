#!/bin/bash

# Test script for JatiStore Authentication API
BASE_URL="http://localhost:8080/api/v1"

echo "🧪 Testing JatiStore Authentication API"
echo "========================================"

# Test 1: Register a new user
echo ""
echo "1. Testing user registration..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@jatistore.com",
    "password": "admin123",
    "role": "admin"
  }')

echo "Register Response: $REGISTER_RESPONSE"

# Test 2: Login
echo ""
echo "2. Testing user login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }')

echo "Login Response: $LOGIN_RESPONSE"

# Extract token from login response
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo ""
    echo "✅ Login successful! Token: ${TOKEN:0:20}..."
    
    # Test 3: Get profile (protected route)
    echo ""
    echo "3. Testing protected route (get profile)..."
    PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/auth/profile" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")
    
    echo "Profile Response: $PROFILE_RESPONSE"
    
    # Test 4: Try to access products without token (should fail)
    echo ""
    echo "4. Testing unauthorized access (should fail)..."
    UNAUTHORIZED_RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
      -H "Content-Type: application/json")
    
    echo "Unauthorized Response: $UNAUTHORIZED_RESPONSE"
    
    # Test 5: Access products with token (should succeed)
    echo ""
    echo "5. Testing authorized access to products..."
    AUTHORIZED_RESPONSE=$(curl -s -X GET "$BASE_URL/products" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json")
    
    echo "Authorized Response: $AUTHORIZED_RESPONSE"
    
else
    echo ""
    echo "❌ Login failed! Cannot proceed with protected route tests."
fi

echo ""
echo "🏁 Authentication tests completed!" 