# vibium-rpa Roadmap

This document outlines planned features and improvements for vibium-rpa.

## v0.2 - Enhanced Activities

- [ ] **http.put** - HTTP PUT requests
- [ ] **http.delete** - HTTP DELETE requests
- [ ] **http.patch** - HTTP PATCH requests
- [ ] **data.transform** - Data transformation with expressions
- [ ] **data.filter** - Filter arrays based on conditions
- [ ] **data.map** - Map array elements
- [ ] **control.while** - While loops
- [ ] **control.break** - Break out of loops
- [ ] **control.return** - Early return from workflow
- [ ] **util.email** - Send email notifications
- [ ] **util.slack** - Send Slack messages
- [ ] **file.copy** - Copy files
- [ ] **file.move** - Move/rename files
- [ ] **file.list** - List directory contents

## v0.3 - Recorder

- [ ] Browser extension for recording user actions
- [ ] Convert recordings to YAML workflows
- [ ] Screenshot-based element suggestions
- [ ] Automatic selector generation
- [ ] Recording playback and editing

## v0.4 - Server Mode

- [ ] REST API for workflow management
  - [ ] CRUD operations for workflows
  - [ ] Execution endpoints
  - [ ] Status and results retrieval
- [ ] Workflow scheduling (cron-like)
  - [ ] Schedule definitions in YAML
  - [ ] Timezone support
  - [ ] Execution history
- [ ] Execution history and logs
  - [ ] Persistent storage
  - [ ] Log aggregation
  - [ ] Metrics collection
- [ ] Web dashboard
  - [ ] Workflow editor
  - [ ] Execution monitoring
  - [ ] Results visualization
  - [ ] Screenshot gallery

## v0.5 - Distributed Execution

- [ ] Agent-based architecture
  - [ ] Worker registration
  - [ ] Health monitoring
  - [ ] Automatic reconnection
- [ ] Remote browser pools
  - [ ] Browser farm integration
  - [ ] Cloud browser providers
- [ ] Load balancing
  - [ ] Round-robin scheduling
  - [ ] Priority queues
  - [ ] Resource-based routing
- [ ] Centralized orchestration
  - [ ] Workflow coordination
  - [ ] State management
  - [ ] Failure recovery

## v0.6 - AI Integration

- [ ] LLM-assisted workflow generation
  - [ ] Natural language to YAML
  - [ ] Step suggestions
  - [ ] Error explanation
- [ ] Smart selectors
  - [ ] Visual element matching
  - [ ] Accessibility-based selection
  - [ ] Multiple fallback strategies
- [ ] Self-healing selectors
  - [ ] Automatic selector repair
  - [ ] Change detection
  - [ ] Notification on changes
- [ ] Natural language step descriptions
  - [ ] Human-readable logs
  - [ ] Report generation
  - [ ] Documentation extraction

## v0.7 - Enterprise Features

- [ ] Authentication providers
  - [ ] SSO integration
  - [ ] OAuth2/OIDC
  - [ ] LDAP/Active Directory
- [ ] Role-based access control
- [ ] Audit logging
- [ ] Secrets management
  - [ ] Vault integration
  - [ ] Encrypted variables
- [ ] Compliance reporting
- [ ] SLA monitoring

## v0.8 - Testing Integration

- [ ] Visual regression testing
  - [ ] Baseline management
  - [ ] Diff generation
  - [ ] Threshold configuration
- [ ] Accessibility testing (WCAG)
  - [ ] Integration with vibium-wcag
  - [ ] Automated reports
- [ ] Performance testing
  - [ ] Page load metrics
  - [ ] Core Web Vitals
  - [ ] Custom timing marks
- [ ] Test case generation
  - [ ] From user stories
  - [ ] From API specs
  - [ ] From recorded sessions

## Future Considerations

- Mobile automation support
- Desktop automation (native apps)
- OCR-based element detection
- PDF form filling
- Database activities (SQL queries)
- Message queue integration (Kafka, RabbitMQ)
- Kubernetes operator for workflow scheduling
- Multi-language workflow authoring (Python, JavaScript)

---

## Contributing

We welcome contributions! If you'd like to work on any of these features:

1. Check if there's an existing issue for the feature
2. Create an issue if one doesn't exist
3. Discuss the implementation approach
4. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for more details.
