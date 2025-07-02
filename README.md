# Nuclio OPA Client

A Go client library for Open Policy Agent (OPA) with support for HTTP-based policy queries.

## Features

- 🚀 **Multiple Client Types**: HTTP, Mock, and No-op clients
- 🔄 **Retry Logic**: Built-in retry mechanism for HTTP requests
- 📊 **Batch Queries**: Query permissions for multiple resources at once
- 🛡️ **Override Support**: Bypass policy checks with override headers
- 🔧 **Configurable**: Flexible configuration options
- 🧪 **Well Tested**: Comprehensive test coverage
- 📝 **Structured Logging**: Integration with nuclio logger

## Installation

```bash
go get github.com/nuclio/opa
```

## Quick Start

```go
package main

import (
    "context"
    "time"
    
    "github.com/nuclio/logger"
    "github.com/nuclio/opa"
)

func main() {
    // Create configuration
    config := &opa.Config{ 
        ClientKind:           opa.ClientKindHTTP,
        Address:              "http://localhost:8181",
        PermissionQueryPath:  "/v1/data/authz/allow",
        PermissionFilterPath: "/v1/data/authz/filter_allowed",
        RequestTimeout:       10,
        Verbose:              false,
    }
    
    // Create client
    logger := // your logger instance
    client := opa.CreateOpaClient(logger, config)
    
    // Query single permission
    allowed, err := client.QueryPermissions(
        "resource1",
        opa.ActionRead,
        &opa.PermissionOptions{
            MemberIds: []string{"user123"},
        },
    )
    
    // Query multiple permissions
    permissions, err := client.QueryPermissionsMultiResources(
        context.Background(),
        []string{"resource1", "resource2"},
        opa.ActionRead,
        &opa.PermissionOptions{
            MemberIds: []string{"user123"},
        },
    )
}
```

## Configuration

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `ClientKind` | `ClientKind` | Type of client (`http`, `nop`, `mock`) | `nop` |
| `Address` | `string` | OPA server URL | - |
| `PermissionQueryPath` | `string` | Single permission query endpoint | - |
| `PermissionFilterPath` | `string` | Multi-resource query endpoint | - |
| `RequestTimeout` | `int` | HTTP timeout in seconds | 10 |
| `Verbose` | `bool` | Enable verbose logging | `false` |
| `OverrideHeaderValue` | `string` | Value for bypass functionality | - |

## Client Types

### HTTP Client
Production client that communicates with OPA over HTTP.

### No-op Client  
Always returns `true` for all permission checks. Useful for development/testing.

### Mock Client
Test client using `testify/mock` for unit testing.

## Actions

Supported actions: `read`, `create`, `update`, `delete`

## Contributing

### Prerequisites
- Go 1.23+
- Make

### Format Code
```bash
make fmt
```

### Testing
```bash
make test
make test-coverage
```

### Linting
```bash
make lint
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
