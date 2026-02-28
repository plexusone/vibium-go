# Test Suite

## Structure

```
tests/
├── e2e/                    # Deterministic test scripts
│   ├── smoke/              # Basic smoke tests
│   │   ├── homepage.json
│   │   └── navigation.json
│   ├── auth/               # Authentication tests
│   │   ├── login.json
│   │   └── logout.json
│   └── a11y/               # Accessibility tests (WCAG 2.2 AA)
│       ├── wcag22aa.json
│       └── login-a11y.json
├── plans/                  # Human-readable test plans
│   ├── smoke.md
│   └── auth.md
└── README.md
```

## Running Tests

### All Tests

```bash
# Run all smoke tests
for f in tests/e2e/smoke/*.json; do vibium run "$f" --headless; done

# Run all auth tests
for f in tests/e2e/auth/*.json; do vibium run "$f" --headless; done

# Run all accessibility tests (WCAG 2.2 AA)
for f in tests/e2e/a11y/*.json; do vibium run "$f" --headless; done
```

### Individual Tests

```bash
vibium run tests/e2e/smoke/homepage.json --headless
vibium run tests/e2e/auth/login.json --headless
```

### With Debug Output

```bash
VIBIUM_DEBUG=1 vibium run tests/e2e/smoke/homepage.json
```

## Writing New Tests

### From Markdown Plan

1. Write a test plan in `tests/plans/`
2. Use LLM with MCP to execute the plan
3. Start recording: `start_recording`
4. Execute steps
5. Export: `export_script`
6. Save to `tests/e2e/<category>/`

### Manual

Create a JSON file following the schema:

```json
{
  "name": "Test Name",
  "description": "What this tests",
  "version": 1,
  "steps": [
    {"action": "navigate", "url": "https://example.com"},
    {"action": "assertTitle", "expected": "Example"}
  ]
}
```

See `script/vibium-script.schema.json` for full schema.

## Test Sites

| Site | URL | Use Case |
|------|-----|----------|
| Example.com | https://example.com | Basic smoke tests |
| The Internet | https://the-internet.herokuapp.com | Interactive patterns |

## CI/CD

Tests run automatically via GitHub Actions:

- On push to main
- On pull requests
- Manual trigger via workflow_dispatch

See `.github/workflows/e2e.yaml`
