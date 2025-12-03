# Logger Package

A flexible logging package for Go applications that provides styled terminal output using charmbracelet/log with
optional structured logging capabilities.

## Quick Start

```go
import "github.com/EncoreDigitalGroup/golib/logger"

// Basic usage - styled terminal output
logger.Info("Application started")
logger.Error("Something went wrong", "error", "connection failed")
logger.Debug("Debug information")
```

## Features

### Core Functionality (Stable)

- **Styled terminal output** using charmbracelet/log
- **Multiple log levels**: Info, Error, Debug, Warn, Print
- **Formatted logging**: `Infof`, `Errorf`, etc.
- **Key-value pairs**: Structured context in terminal output
- **Custom styling**: Error highlighting and colors
- **Package-level functions**: Use without creating instances
- **Logger instances**: Create multiple loggers with different configurations

### Experimental Features

- **Structured logging** with Go's standard `log/slog`
- **Dual backend** for simultaneous terminal + structured output
- **Environment-based configuration** for different deployment scenarios

See the [experimental features documentation](experimental-features.md).

## Basic Usage

### Package Functions

```go
// Different log levels
logger.Info("User logged in", "user_id", 123)
logger.Error("Database error", "error", err.Error())
logger.Debug("Cache miss", "key", "user:123")
logger.Warn("Rate limit approaching", "remaining", 10)

// Formatted logging
logger.Infof("Processing %d files", fileCount)
logger.Errorf("Failed to connect to %s", hostname)
```

### Logger Instances

```go
// Create custom logger instances
appLogger := logger.New()
dbLogger := logger.New()

appLogger.Info("Application started")
dbLogger.Error("Connection failed")
```

## Output Examples

### Terminal Output

```
INFO Application started
ERROR Something went wrong error="connection failed"
WARN Rate limit approaching remaining=10
```

The terminal output includes:

- **Colored level indicators**
- **Styled error messages** with red background
- **Key-value formatting** for structured data
- **Readable timestamps** and formatting

## Configuration

### Default Logger

The package provides a default logger instance that's ready to use:

```go
// Uses the default styled terminal logger
logger.Info("Ready to go!")
```

### Custom Styling

The logger uses charmbracelet/log with custom styling:

- Error messages have red background with white text
- Key-value pairs are highlighted
- Clean, readable terminal output

## Use Cases

### CLI Applications

Perfect for command-line tools where users need clear, readable output:

```go
logger.Info("Starting backup process...")
logger.Infof("Processing file %d of %d", current, total)
logger.Error("Backup failed", "file", filename, "error", err.Error())
```

### Development

Great for development environments where you want styled terminal output:

```go
logger.Debug("Cache configuration", "ttl", "5m", "size", 1000)
logger.Info("Server starting", "port", 8080)
logger.Warn("Deprecated API used", "endpoint", "/old-api")
```

### Interactive Applications

Ideal for applications that need to display information to users in real-time:

```go
logger.Info("✓ Database connected")
logger.Info("✓ Cache initialized")
logger.Error("✗ External service unavailable")
```

## API Reference

### Log Levels

| Function               | Description                |
|------------------------|----------------------------|
| `Print()` / `Printf()` | Plain output without level |
| `Info()` / `Infof()`   | Informational messages     |
| `Error()` / `Errorf()` | Error conditions           |
| `Debug()` / `Debugf()` | Debug information          |
| `Warn()` / `Warnf()`   | Warning conditions         |

### Logger Methods

All package functions are also available as methods on logger instances:

```go
logger := logger.New()
logger.Info("Instance method")
logger.Errorf("Error: %v", err)
```

## Thread Safety

The logger is thread-safe and can be used concurrently from multiple goroutines.

## Performance

- Minimal overhead for terminal output
- Efficient string formatting
- No file I/O in standard configuration
- Optimized for interactive use

## Dependencies

- `github.com/charmbracelet/log` - Styled terminal logging
- `github.com/charmbracelet/lipgloss` - Terminal styling (via charmbracelet/log)

## Backwards Compatibility

The core logger API is stable and maintains backwards compatibility. All existing code will continue to work without
changes.

For experimental features and their stability guarantees, see [Experimental Features](experimental-features.md).