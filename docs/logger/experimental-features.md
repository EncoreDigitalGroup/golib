# Experimental Features

This package contains experimental features that are subject to change in future versions. Experimental features do
not follow semantic versioning guarantees and may be modified or removed without notice.

## Enabling Experimental Features

To use experimental features, you must build with the `experimental` build tag:

```bash
# Build with experimental features
go build -tags experimental

# Run tests with experimental features
go test -tags experimental

# Install with experimental features
go install -tags experimental
```

## Available Experimental Features

### 1. Slog Integration (`slog.go`)

Structured logging integration with Go's standard `log/slog` package:

```go
// Enable JSON structured logging
logger.SetSlog(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
logger.Info("User login", "user_id", 123, "ip", "192.168.1.1")

// Create slog-only logger
jsonLogger := logger.NewSlogJSON(logFile)
```

### 2. Dual Backend (`dual.go`)

Simultaneous output to both terminal (charm) and structured logging:

```go
// CLI applications: pretty terminal + debug file
logger.ConfigureForCLI("debug.log")

// Microservices: environment-based configuration
logger.ConfigureForMicroservice(isDev, logOutput)

// Manual dual output
logger := logger.NewDualLoggerJSON(logFile)
logger.Info("Message") // Both terminal styling AND JSON file
```

## Stability

⚠️ **Warning**: Experimental features may:

- Change API without notice
- Be removed in future versions
- Have incomplete documentation or testing
- Contain bugs or performance issues

Use experimental features only for:

- Testing and evaluation
- Non-production environments
- When you accept the risk of API changes

## Migration Path

When experimental features become stable:

1. The `experimental` build tag will be removed
2. APIs may change based on feedback
3. Features will be included in standard builds
4. Comprehensive documentation will be provided

## Feedback

Please provide feedback on experimental features through GitHub issues, including:

- API design suggestions
- Bug reports
- Performance observations
- Use case requirements