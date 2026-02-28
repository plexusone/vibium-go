# Smoke Test Plan

## Overview

Basic smoke tests to verify the application is working.

## Test Cases

### 1. Homepage Loads

**Steps:**

1. Navigate to the homepage
2. Verify the page title is correct
3. Verify the main heading is visible

**Expected:** Page loads with correct title and heading

### 2. Navigation Works

**Steps:**

1. Navigate to the homepage
2. Click the main link
3. Verify navigation completes

**Expected:** Link navigation works correctly

## Execution

```bash
vibium run tests/e2e/smoke/homepage.json
vibium run tests/e2e/smoke/navigation.json
```

## CI/CD

These tests run on every push to main.
