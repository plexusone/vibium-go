# Authentication Test Plan

## Overview

Tests for user authentication flows.

## Test Cases

### 1. Login with Valid Credentials

**Steps:**

1. Navigate to login page
2. Enter valid username
3. Enter valid password
4. Click login button
5. Verify redirect to secure area
6. Verify success message

**Expected:** User is logged in and sees success message

### 2. Logout

**Steps:**

1. Log in (prerequisite)
2. Click logout button
3. Verify redirect to login page
4. Verify logout message

**Expected:** User is logged out and sees logout message

### 3. Login with Invalid Credentials (TODO)

**Steps:**

1. Navigate to login page
2. Enter invalid credentials
3. Click login button
4. Verify error message

**Expected:** User sees error message and remains on login page

## Test Site

Using [The Internet](https://the-internet.herokuapp.com/) for testing.

**Credentials:**

- Username: `tomsmith`
- Password: `SuperSecretPassword!`

## Execution

```bash
vibium run tests/e2e/auth/login.json
vibium run tests/e2e/auth/logout.json
```

## CI/CD

These tests run on every push and PR.
