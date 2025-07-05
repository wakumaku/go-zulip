# Contributing to go-zulip

Thank you for your interest in contributing to go-zulip! This document provides complete guidelines for contributing to this project, including development setup, testing, and integration with Zulip.

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/wakumaku/go-zulip.git
   cd go-zulip
   ```

2. **Build development environment**
   ```bash
   make dev-setup
   ```

3. **Start development environment**
   ```bash
   make dev-up
   ```

4. **Run tests**
   ```bash
   make test
   ```

## Development Environment Setup

The project uses a Docker-based development environment that provides:
- Go 1.24 with all required tools (`golangci-lint`, `gofumpt`, `air`)
- A complete Zulip server instance for integration testing
- Automatic file watching and test execution
- Consistent development environment across different machines

### Initial Setup

1. **Build the development environment**
   ```bash
   make dev-setup
   ```
   This will:
   - Copy `docker-compose-dev-env.example.yml` to `docker-compose-dev-env.yml` if needed
   - Build the Docker images with the latest Go version and tools
   - Set up the complete development stack

2. **Start the development environment**
   ```bash
   make dev-up
   ```
   This starts all services including:
   - Zulip server (accessible at `https://localhost`)
   - PostgreSQL database
   - Redis cache
   - RabbitMQ message broker
   - Memcached
   - go-zulip development container

3. **View logs** (optional)
   ```bash
   make dev-logs
   ```

### Setting Up Zulip for Integration Testing

Once the development environment is running, you'll need to configure Zulip for integration testing:

1. **Create a new Zulip organization**
   ```bash
   docker exec -u zulip go-zulip-zulip-1 /home/zulip/deployments/current/manage.py generate_realm_creation_link
   ```
   This will output a link like:
   ```
   https://localhost/new/jj46mgxl7cbxazmoi7evtxjx
   ```

2. **Register your organization**
   - Go to `https://localhost` (accept "Proceed to localhost (unsafe)")
   - Visit the registration link from step 1
   - Create your organization and admin user

3. **Get your API Key**
   - Go to `https://localhost`
   - Navigate to: Profile > Settings > Account & Privacy > API Key
   - Copy the API key

4. **Configure environment variables**
   Edit `docker-compose-dev-env.yml` and add your credentials:
   ```yaml
   environment:
     - ZULIP_EMAIL=your_email@test.com
     - ZULIP_API_KEY=your_api_key_here
     - ZULIP_SITE=https://localhost
   ```

5. **Grant API permissions**
   ```bash
   docker exec -u zulip go-zulip-zulip-1 /home/zulip/deployments/current/manage.py change_user_role your_email@test.com can_create_users -r 2
   ```

6. **Restart the development container**
   ```bash
   make dev-down
   make dev-up
   ```

## Development Workflow

### Available Make Targets

| Command | Description |
|---------|-------------|
| `make dev-setup` | Build development environment (Docker compose build) |
| `make dev-up` | Start development environment (Docker compose up) |
| `make dev-down` | Stop development environment |
| `make dev-logs` | Show development environment logs |
| `make dev-exec CMD="command"` | Execute a command in the development container |
| `make test` | Run all tests (unit tests only by default) |
| `make test-unit` | Run unit tests |
| `make test-integration` | Run integration tests (requires running dev env) |
| `make test-coverage` | Run tests with coverage report |
| `make lint` | Run linters |
| `make fmt` | Format code |
| `make build` | Build the library (verify it compiles) |
| `make tidy` | Tidy go modules |
| `make ci` | Run CI pipeline (with Docker) |
| `make ci-local` | Run CI pipeline locally (without Docker) |
| `make verify` | Verify the project is ready for commit |
| `make clean` | Clean generated files |

### Docker vs Local Development

All development commands automatically detect if the Docker environment is running:
- **With Docker**: Commands run inside the container with the exact Go version and tools
- **Without Docker**: Commands fall back to running locally

This ensures you can always work on the project, even without Docker.

### Daily Development Workflow

1. **Start your development session**
   ```bash
   make dev-up
   ```

2. **Make your changes** to the Go code

3. **Watch the automatic test run on every file save or run tests frequently**
   ```bash
   make dev-logs           # Show logs and watch for changes
   ```
   Or
   ```bash
   make test-unit          # Quick unit tests
   make test-integration   # Full integration tests with Zulip
   ```

4. **Check code quality**
   ```bash
   make lint              # Run linters
   make fmt               # Format code
   ```

5. **Verify everything before committing**
   ```bash
   make verify            # Run complete verification
   ```

6. **Stop development environment when done**
   ```bash
   make dev-down
   ```

### File Watching and Auto-testing

The development environment includes `air` (file watcher) that automatically:
- Runs linters on code changes
- Executes tests when `.go` files are modified
- Provides immediate feedback during development

## Code Style and Standards

This project follows standard Go conventions:

- **Formatting**: Use `make fmt` (runs `gofumpt`)
- **Linting**: Use `make lint` (runs `golangci-lint` with project configuration)
- **Testing**: Add tests for all new functionality
- **Documentation**: Update documentation and examples as needed
- **Error Handling**: Follow Go best practices for error handling
- **API Consistency**: Follow existing patterns for request/response structures

### Go Module Management

The project uses Go 1.24 with the modern `tool` directive for tool management:
- Tools are declared in `go.mod` under the `tool` section
- All tools are available via `go tool <tool-name>`
- Use `make tidy` to ensure dependencies are up to date

## Running Tests

### Unit Tests
```bash
make test-unit
```
- Fast execution
- No external dependencies
- Use mocks for external services
- Run in both Docker and local environments

### Integration Tests
```bash
make test-integration
```
- Require a running Zulip server
- Test real API interactions
- Verify end-to-end functionality
- Must run with `make dev-up` environment

### Coverage Report
```bash
make test-coverage
```
- Generates `coverage.out` and `coverage.html`
- View detailed coverage in your browser
- Aim for high coverage on new code

## Code Quality Checks

Before submitting a pull request, ensure your code passes all checks:

```bash
make verify
```

This comprehensive check runs:
- `make ci` (which includes tidy, lint, and tests)
- All quality gates required for the project

For GitHub Actions compatibility, use:
```bash
make ci-local
```

## Pull Request Process

1. **Fork the repository** and create a feature branch from `main`
2. **Set up your development environment**
   ```bash
   make dev-setup
   make dev-up
   ```
3. **Make your changes** following the code style guidelines
4. **Add comprehensive tests** for any new functionality
5. **Test with integration environment**
   ```bash
   make test-integration
   ```
6. **Verify code quality**
   ```bash
   make verify
   ```
7. **Update documentation** if necessary
8. **Submit a pull request** with a clear description of your changes

### Pull Request Checklist

- [ ] Code follows project conventions
- [ ] Tests are added for new functionality
- [ ] All tests pass (`make test`)
- [ ] Integration tests pass (`make test-integration`)
- [ ] Code is properly formatted (`make fmt`)
- [ ] Linter checks pass (`make lint`)
- [ ] Documentation is updated
- [ ] Changes are described in the PR description

## Adding New API Endpoints

When adding support for new Zulip API endpoints:

1. **Choose the appropriate package** (e.g., `messages/`, `channels/`, `users/`)
2. **Follow existing patterns** for request/response structures
3. **Add comprehensive tests** including:
   - Unit tests with mocks
   - Integration tests with real Zulip API
4. **Update documentation** and examples
5. **Test with both Docker and local environments**

### Example Structure for New Endpoint

```go
// 1. Define request/response structures
type CreateSomethingRequest struct {
    Name string `json:"name"`
}

type CreateSomethingResponse struct {
    ID int `json:"id"`
}

// 2. Implement the function
func (c *Client) CreateSomething(req CreateSomethingRequest) (*CreateSomethingResponse, error) {
    // Implementation
}

// 3. Add unit tests
func TestCreateSomething(t *testing.T) {
    // Test implementation
}

// 4. Add integration tests
func TestCreateSomethingIntegration(t *testing.T) {
    // Integration test with real Zulip
}
```

## Directory Structure

```
.
├── channels/                    # Channel-related operations
├── messages/                    # Message operations  
├── users/                      # User management
├── realtime/                   # Real-time events
├── examples/                   # Usage examples
├── test/                       # Integration tests
├── docs/                       # Documentation
├── docker-compose-*.yml        # Docker environment configuration
├── go.mod                      # Go module with tool dependencies
├── golangci.yml               # Linter configuration
└── Makefile                   # Development workflow automation
```

## Troubleshooting

### Common Issues

1. **Docker environment not starting**
   ```bash
   make dev-down
   make dev-setup
   make dev-up
   ```

2. **Integration tests failing**
   - Ensure Zulip is properly configured (see Setup section)
   - Check that environment variables are set in `docker-compose-dev-env.yml`
   - Verify API permissions are granted

3. **Tools not found**
   ```bash
   make tidy
   go tool -n  # List available tools
   ```

4. **Permission errors**
   ```bash
   docker exec -u zulip go-zulip-zulip-1 /home/zulip/deployments/current/manage.py change_user_role your_email@test.com can_create_users -r 2
   ```

### Getting Help

- **View logs**: `make dev-logs`
- **Execute commands in container**: `make dev-exec CMD="bash"`
- **Check container status**: `docker compose -p go-zulip ps`
- **Restart environment**: `make dev-down && make dev-up`

## Reporting Issues

When reporting issues:

1. **Use the issue template** if available
2. **Include minimal reproduction steps**
3. **Provide relevant environment information**:
   - Operating system
   - Docker version
   - Go version (if running locally)
4. **Include error messages and stack traces**
5. **Mention if you're using Docker or local development**

## Questions?

Feel free to open an issue for questions or discussions about the project. For development environment issues, include:
- Output of `make dev-logs`
- Your `docker-compose-dev-env.yml` configuration (without sensitive data)
- Steps you've already tried
