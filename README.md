# go-log

A high-performance, asynchronous, resilient, and multi-channel logging library for Go. Designed for applications that require non-blocking operations, featuring native support for multiple destinations and extensible formatting dispatch.

---

## Features

- Non-Blocking Architecture: All write operations are processed in the background using buffered channels to ensure zero latency on the main application flow.
- Resilient File Handling: The internal file wrapper automatically detects external deletions, reopens descriptors, and ensures data integrity.
- Multi-Provider Orchestration: Dispatch the same log record to multiple destinations (Console, File, Network Streams) simultaneously.
- Flexible Formatting: Decoupled formatting logic, allowing for easy implementation of Plain Text, NDJSON (JSON Lines), or custom protocols.
- Precise Timestamps: Millisecond-precision timestamps with a clock abstraction for deterministic unit testing.

---

## Installation

```sh
go get github.com/Rafael24595/go-log
```

---

## Quick Start

Configure a logger that writes to the Console (Text) and a File (JSON Lines) at the same time:

```go
package main

import (
	"context"
    "fmt"

	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/format/json"
	"github.com/Rafael24595/go-log/log/provider/console"
	"github.com/Rafael24595/go-log/log/provider/file"
	"github.com/Rafael24595/go-log/log/provider/multi"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Define your providers
	p := multi.New(
		console.New(), // Standard text output
		file.FileProvider{
			Path:    "./logs",
			Session: "api-v1",
			Format:  &json.JsonLineFormat, // Structured JSON output
		},
	)

	// 2. Initialize the global logger
	if err := log.DefaultFromProvider(ctx, p); err != nil {
		panic(err)
	}

	// Ensure all buffers are flushed on exit
	defer log.OnClose()

	// 3. Start logging
	log.Message("Server started on port 8080")
	log.Custom("AUDIT", "User 123 accessed the dashboard")
	log.Error(fmt.Errorf("database connection timeout"))
}
```

---

## Advanced Features

### Bootstrap Logging

One of the most powerful features of this library is the **Bootstrap mechanism**. If you attempt to log messages before a `Default` logger is configured, the system doesn't discard them. Instead:
1. It buffers the records in an internal **Bootstrap** memory.
2. Once `DefaultFromProvider` is called, all previously buffered records are automatically flushed into the new logger.
3. This ensures you never lose critical startup logs.

## Extensibility: Custom Implementations

The library is designed to be extended without modifying the core. You can implement your own transport layers and data formats.

### 1. Custom Log Categories

You are not limited to `Message`, `Warning`, or `Error`. You can define your own domain-specific categories:

```go
// Custom category logging
log.Custom("DATABASE", "Connection established")
log.Customf("SECURITY", "Unauthorized access attempt from IP: %s", ip)
```

The system automatically normalizes these categories (e.g., converting them to uppercase) to maintain a consistent format in your files or console.

### 2. Custom Providers

A `Provider` is responsible for building and initializing a `Log` engine. To create a custom one (e.g., an HTTP Webhook provider), implement the `Provider` interface:

```go
type WebhookProvider struct {
    URL string
}

func (p WebhookProvider) Build(ctx context.Context) (log.Log, error) {
    // 1. Initialize your custom engine/client
    client := myHttpClient.New(p.URL)
    
    // 2. Return an instance that satisfies the log.Log interface
    // You can use the internal 'engine' to handle async queuing
    return myCustomLogger{client: client}, nil
}
```

### 3. Custom Log Implementations

If the default asynchronous engine doesn't fit your needs, you can implement the `Log` interface entirely. This is useful for creating loggers with custom filtering, real-time alerting, or specialized buffering logic.

```go
type MyCustomLogger struct {
    // your internal state (e.g., mutexes, external clients, filters)
}

// Implement the Log interface:
func (m *MyCustomLogger) Name() logger.Logger { return "CustomTarget" }
func (m *MyCustomLogger) Message(msg string) record.Record { /* ... */ }
func (m *MyCustomLogger) Error(err error) record.Record { /* ... */ }
// ... implement all methods from the Log interface
```

### 4. Custom Formats

Since the Format system uses a struct with function pointers (vtable), you can create stateful formatters by injecting instance methods. This is useful for CSV headers, counters, or complex encodings:

```go
type AuditFormatter struct {
    AppID string
}

func (f *AuditFormatter) Serialize(records ...record.Record) (string, error) {
    // Custom logic using f.AppID
    return formattedString, nil
}

// Usage:
myFormat := format.Format{
    Extension: "audit",
    Format:    auditInst.Serialize, // Injecting the method
}
```

## Standard Library Bridge (io.Writer)

One of the most versatile features of **Go-Log** is its ability to act as an `io.Writer`. This allows you to plug the asynchronous engine into any standard Go component.

### Dynamic Proxying

The writer returned by `WriterFromString()` or `WriterFromCategory()` does not store a static reference to the logger. Instead, it always resolves to the **current active logger**. 

1. **Phase 1 (Bootstrap)**: If an HTTP server writes to the proxy before you initialize your final provider, the logs go to the Bootstrap buffer.
2. **Phase 2 (Handover)**: As soon as `DefaultFromProvider` is called, the **very next byte** written to the proxy will go directly to your new destination (File, Console, etc.).

### Usage Examples

#### 1. Redirecting the Standard `log` Package

```go
package main

import (
	"log"

	go_log "github.com/Rafael24595/go-log/log"
)

func main() {
	// Redirect all log.Print() calls to our system under the "LEGACY" category
	log.SetOutput(go_log.WriterFromString("LEGACY"))

	log.Println("This standard call is now asynchronous!")
}
```

#### 2. Capturing HTTP Server Errors

```go
package main

import (
	"log"
	"net/http"

	go_log "github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/model/record"
)

const HTTP_INTERNAL record.Category = "HTTP-INTERNAL"

func main() {
	server := &http.Server{
		Addr: ":8080",
		// Capture internal server errors into our "HTTP-INTERNAL" category
		ErrorLog: log.New(go_log.WriterFromCategory(HTTP_INTERNAL), "", 0),
	}
}

```

---

## Architecture

The library is built on decoupled components to ensure maximum extensibility:
- Engine: The core concurrent processor managing the background goroutine and record history.
- Providers: High-level abstractions to build specific loggers (Console, File, Stream, Multi).
- Formats: Serialization logic separated from transport. It uses an injectable function approach (VTable) to allow for both stateless and stateful formatters.

---

## Thread Safety & Reliability

- Concurrent-Safe: All providers use internal Mutexes and Atomic Booleans to manage state across multiple goroutines.
- Graceful Shutdown: The Close() method ensures that all pending records in the channel are processed before the engine shuts down.
- Atomic Writes: Includes utilities for atomic file writes (write-then-rename) to prevent file corruption during system crashes.
