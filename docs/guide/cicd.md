# CI/CD Integration

Run deterministic test scripts in CI/CD pipelines without an LLM.

!!! warning "Clicker Binary Required"
    E2E tests require the Vibium clicker binary, which is not yet publicly distributed. The examples below assume you have access to the clicker binary. See [Prerequisites](../getting-started/prerequisites.md) for details.

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

on:
  workflow_dispatch:
    inputs:
      clicker_url:
        description: 'URL to download clicker binary'
        required: true
        type: string

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Download Clicker
        run: |
          curl -L -o clicker "${{ github.event.inputs.clicker_url }}"
          chmod +x clicker
          echo "VIBIUM_CLICKER_PATH=$PWD/clicker" >> $GITHUB_ENV

      - name: Install Vibium CLI
        run: go install github.com/agentplexus/vibium-go/cmd/vibium@latest

      - name: Run E2E Tests
        env:
          VIBIUM_HEADLESS: "1"
        run: |
          vibium run tests/login.json
          vibium run tests/checkout.json
```

!!! note "Manual Trigger"
    Until the clicker is publicly distributed, E2E workflows should use `workflow_dispatch` with a `clicker_url` input rather than automatic triggers on push/PR.

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

      - name: Download Clicker
        run: |
          curl -L -o clicker "${{ github.event.inputs.clicker_url }}"
          chmod +x clicker
          echo "VIBIUM_CLICKER_PATH=$PWD/clicker" >> $GITHUB_ENV

      - name: Setup
        run: go install github.com/agentplexus/vibium-go/cmd/vibium@latest

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
  image: golang:1.22

  variables:
    CLICKER_URL: "$CLICKER_URL"  # Set in CI/CD variables
    VIBIUM_HEADLESS: "1"

  before_script:
    - curl -L -o clicker "$CLICKER_URL"
    - chmod +x clicker
    - export VIBIUM_CLICKER_PATH=$PWD/clicker
    - go install github.com/agentplexus/vibium-go/cmd/vibium@latest

  script:
    - vibium run tests/smoke.json
    - vibium run tests/auth.json

  artifacts:
    when: on_failure
    paths:
      - screenshots/
    expire_in: 1 week

  rules:
    - when: manual  # Manual trigger until clicker is public
```

## CircleCI

```yaml
version: 2.1

jobs:
  e2e:
    docker:
      - image: cimg/go:1.22
    steps:
      - checkout
      - run:
          name: Download Clicker
          command: |
            curl -L -o clicker "$CLICKER_URL"
            chmod +x clicker
            echo "export VIBIUM_CLICKER_PATH=$PWD/clicker" >> $BASH_ENV
      - run:
          name: Install Vibium CLI
          command: go install github.com/agentplexus/vibium-go/cmd/vibium@latest
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
      - e2e:
          # Manual approval until clicker is public
          type: approval
```

## Jenkins

```groovy
pipeline {
    agent any

    parameters {
        string(name: 'CLICKER_URL', description: 'URL to download clicker binary')
    }

    environment {
        VIBIUM_HEADLESS = '1'
    }

    stages {
        stage('Setup') {
            steps {
                sh '''
                    curl -L -o clicker "${CLICKER_URL}"
                    chmod +x clicker
                '''
                sh 'go install github.com/agentplexus/vibium-go/cmd/vibium@latest'
            }
        }

        stage('E2E Tests') {
            environment {
                VIBIUM_CLICKER_PATH = "${WORKSPACE}/clicker"
            }
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
trigger: none  # Manual trigger until clicker is public

parameters:
  - name: clickerUrl
    displayName: 'Clicker Download URL'
    type: string

pool:
  vmImage: 'ubuntu-latest'

steps:
  - task: GoTool@0
    inputs:
      version: '1.22'

  - script: |
      curl -L -o clicker "${{ parameters.clickerUrl }}"
      chmod +x clicker
      echo "##vso[task.setvariable variable=VIBIUM_CLICKER_PATH]$(pwd)/clicker"
    displayName: 'Download Clicker'

  - script: go install github.com/agentplexus/vibium-go/cmd/vibium@latest
    displayName: 'Install Vibium CLI'

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

For WCAG 2.2 accessibility testing in CI/CD, use [agent-a11y](https://github.com/agentplexus/agent-a11y):

```yaml
name: Accessibility

on:
  workflow_dispatch:
    inputs:
      clicker_url:
        description: 'URL to download clicker binary'
        required: true

jobs:
  wcag:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download Clicker
        run: |
          curl -L -o clicker "${{ github.event.inputs.clicker_url }}"
          chmod +x clicker
          echo "VIBIUM_CLICKER_PATH=$PWD/clicker" >> $GITHUB_ENV

      - name: Setup
        run: go install github.com/agentplexus/agent-a11y/cmd/agent-a11y@latest

      - name: Run WCAG 2.2 AA Evaluation
        env:
          VIBIUM_HEADLESS: "1"
        run: |
          agent-a11y vpat https://staging.example.com \
            --format json --output wcag-results.json

      - name: Upload WCAG Results
        uses: actions/upload-artifact@v4
        with:
          name: wcag-results
          path: wcag-results.json
```

agent-a11y combines:

- **Automated testing** (~40% coverage) - axe-core rule-based checks
- **Specialized automation** (~25% coverage) - keyboard, focus, reflow tests
- **LLM-as-a-Judge** (~25% coverage) - semantic evaluation (optional)

See the [agent-a11y documentation](https://github.com/agentplexus/agent-a11y) for details.

## Example: Complete Workflow

See `.github/workflows/e2e.yaml` in this repository for a complete example.
