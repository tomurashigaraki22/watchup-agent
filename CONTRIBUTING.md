# Contributing to Watchup Server Agent

Thank you for your interest in contributing to the Watchup Server Agent!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/watchup-agent.git`
3. Create a branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Test your changes
6. Commit: `git commit -m "Add your feature"`
7. Push: `git push origin feature/your-feature-name`
8. Open a Pull Request

## Development Setup

```bash
# Install dependencies
go mod tidy

# Run locally
go run cmd/agent/main.go config.yaml

# Build
go build -o watchup-agent cmd/agent/main.go

# Run tests (when available)
go test ./...
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Add comments for exported functions
- Keep functions small and focused

## Commit Messages

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove)
- Reference issues when applicable

Examples:
- `Add disk usage monitoring`
- `Fix memory leak in collector`
- `Update documentation for VPS installation`

## Pull Request Process

1. Update documentation if needed
2. Add tests for new features
3. Ensure all tests pass
4. Update CHANGELOG.md
5. Request review from maintainers

## Reporting Issues

- Use GitHub Issues
- Provide clear description
- Include steps to reproduce
- Add logs if applicable
- Specify OS and Go version

## Questions?

Open a discussion on GitHub or email support@watchup.site
