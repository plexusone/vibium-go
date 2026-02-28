# CI/CD Integration

Run deterministic test scripts in CI/CD pipelines without an LLM.

## Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         Development Workflow                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   Developer writes         LLM explores &         Script saved           │
│   Markdown test plan  ──▶  records actions   ──▶  to repo               │
│                            (with MCP)                                    │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│                           CI/CD Pipeline                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   git push  ──▶  CI runs  ──▶  vibium run test.json  ──▶  Pass/Fail    │
│                  headless         (no LLM needed)                        │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Benefits

| Benefit | Description |
|---------|-------------|
| **No LLM costs** | Scripts run without API calls |
| **Deterministic** | Same inputs → same outputs |
| **Fast** | No LLM latency |
| **Auditable** | Scripts are version-controlled |
| **Parallelizable** | Run multiple scripts concurrently |

## GitHub Actions

### Basic Workflow

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install Vibium
        run: |
          npm install -g vibium
          go install github.com/grokify/vibium-go/cmd/vibium@latest

      - name: Run E2E Tests
        env:
          VIBIUM_HEADLESS: "1"
        run: |
          vibium run tests/login.json
          vibium run tests/checkout.json
```

### Matrix Strategy

Run tests in parallel:

```yaml
jobs:
  e2e:
    runs-on: ubuntu-latest
    strategy:
      fail-fast: false
      matrix:
        test: [smoke, auth, checkout, search]

    steps:
      - uses: actions/checkout@v4

      - name: Setup
        run: |
          npm install -g vibium
          go install github.com/grokify/vibium-go/cmd/vibium@latest

      - name: Run ${{ matrix.test }} tests
        env:
          VIBIUM_HEADLESS: "1"
        run: |
          for script in tests/${{ matrix.test }}/*.json; do
            vibium run "$script"
          done
```

### Upload Artifacts on Failure

```yaml
      - name: Run tests
        run: vibium run tests/e2e.json

      - name: Upload screenshots on failure
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: screenshots
          path: screenshots/
          retention-days: 7
```

### Scheduled Runs

```yaml
on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours
  workflow_dispatch:        # Manual trigger
```

## GitLab CI

```yaml
e2e:
  image: node:20

  before_script:
    - npm install -g vibium
    - apt-get update && apt-get install -y golang
    - go install github.com/grokify/vibium-go/cmd/vibium@latest

  script:
    - export VIBIUM_HEADLESS=1
    - vibium run tests/smoke.json
    - vibium run tests/auth.json

  artifacts:
    when: on_failure
    paths:
      - screenshots/
    expire_in: 1 week
```

## CircleCI

```yaml
version: 2.1

jobs:
  e2e:
    docker:
      - image: cimg/go:1.22-node
    steps:
      - checkout
      - run:
          name: Install Vibium
          command: |
            npm install -g vibium
            go install github.com/grokify/vibium-go/cmd/vibium@latest
      - run:
          name: Run E2E Tests
          environment:
            VIBIUM_HEADLESS: "1"
          command: |
            vibium run tests/smoke.json
            vibium run tests/auth.json
      - store_artifacts:
          path: screenshots
          destination: screenshots

workflows:
  test:
    jobs:
      - e2e
```

## Jenkins

```groovy
pipeline {
    agent any

    environment {
        VIBIUM_HEADLESS = '1'
    }

    stages {
        stage('Setup') {
            steps {
                sh 'npm install -g vibium'
                sh 'go install github.com/grokify/vibium-go/cmd/vibium@latest'
            }
        }

        stage('E2E Tests') {
            steps {
                sh 'vibium run tests/smoke.json'
                sh 'vibium run tests/auth.json'
            }
        }
    }

    post {
        failure {
            archiveArtifacts artifacts: 'screenshots/**', fingerprint: true
        }
    }
}
```

## Azure Pipelines

```yaml
trigger:
  - main

pool:
  vmImage: 'ubuntu-latest'

steps:
  - task: NodeTool@0
    inputs:
      versionSpec: '20.x'

  - task: GoTool@0
    inputs:
      version: '1.22'

  - script: |
      npm install -g vibium
      go install github.com/grokify/vibium-go/cmd/vibium@latest
    displayName: 'Install Vibium'

  - script: |
      export VIBIUM_HEADLESS=1
      vibium run tests/smoke.json
      vibium run tests/auth.json
    displayName: 'Run E2E Tests'

  - task: PublishBuildArtifacts@1
    condition: failed()
    inputs:
      pathToPublish: 'screenshots'
      artifactName: 'screenshots'
```

## Test Organization

Recommended structure:

```
tests/
├── e2e/
│   ├── smoke/
│   │   ├── homepage.json
│   │   └── navigation.json
│   ├── auth/
│   │   ├── login.json
│   │   ├── logout.json
│   │   └── password-reset.json
│   ├── checkout/
│   │   ├── add-to-cart.json
│   │   └── purchase.json
│   └── search/
│       └── basic-search.json
├── plans/
│   ├── smoke.md           # Markdown test plans
│   ├── auth.md
│   └── checkout.md
└── README.md
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VIBIUM_HEADLESS` | Run headless | `false` |
| `VIBIUM_DEBUG` | Enable debug logs | `false` |
| `VIBIUM_CLICKER_PATH` | Path to clicker | Auto-detect |
| `VIBIUM_TIMEOUT` | Default timeout | `30s` |

## Best Practices

### 1. Use Headless Mode

Always set `VIBIUM_HEADLESS=1` in CI:

```yaml
env:
  VIBIUM_HEADLESS: "1"
```

### 2. Set Appropriate Timeouts

CI environments may be slower:

```json
{
  "name": "CI Test",
  "timeout": "60s",
  "steps": [...]
}
```

### 3. Capture Screenshots on Failure

Add screenshot steps for debugging:

```json
{
  "steps": [
    {"action": "navigate", "url": "https://example.com"},
    {"action": "screenshot", "file": "screenshots/step1.png"},
    {"action": "click", "selector": "#submit"},
    {"action": "screenshot", "file": "screenshots/step2.png"}
  ]
}
```

### 4. Use `continueOnError` for Non-Critical Steps

```json
{
  "action": "click",
  "selector": "#optional-banner-close",
  "continueOnError": true
}
```

### 5. Parallelize Independent Tests

Use matrix strategies to run tests concurrently.

### 6. Version Control Test Scripts

- Store scripts in Git alongside code
- Review script changes in PRs
- Track test evolution over time

## Debugging CI Failures

### Enable Debug Logging

```yaml
env:
  VIBIUM_DEBUG: "1"
```

### Download Artifacts

Screenshots and logs uploaded as artifacts help debug failures.

### Run Locally

Reproduce CI failures locally:

```bash
VIBIUM_HEADLESS=1 vibium run tests/failing-test.json
```

## Accessibility Testing in CI/CD

For WCAG 2.2 accessibility testing in CI/CD, use [vibium-wcag](https://github.com/agentplexus/vibium-wcag):

```yaml
name: Accessibility

on: [push, pull_request]

jobs:
  wcag:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup
        run: |
          npm install -g vibium
          go install github.com/agentplexus/vibium-wcag/cmd/vibium-wcag@latest

      - name: Run WCAG 2.2 AA Evaluation
        env:
          VIBIUM_HEADLESS: "1"
        run: |
          vibium-wcag evaluate https://staging.example.com \
            --format json --output wcag-results.json

      - name: Upload WCAG Results
        uses: actions/upload-artifact@v4
        with:
          name: wcag-results
          path: wcag-results.json
```

vibium-wcag combines:

- **Automated testing** (~40% coverage) - axe-core rule-based checks
- **Specialized automation** (~25% coverage) - keyboard, focus, reflow tests
- **LLM-as-a-Judge** (~25% coverage) - semantic evaluation (optional)

See the [vibium-wcag documentation](https://github.com/agentplexus/vibium-wcag) for details.

## Example: Complete Workflow

See `.github/workflows/e2e.yaml` in this repository for a complete example.
