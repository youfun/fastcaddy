# FastCaddy

FastCaddy is a library that provides both Python and Go interfaces for managing Caddy server configurations via the Caddy Admin API.

## Repository Structure

This repository now contains both Python and Go versions of FastCaddy:

```
fastcaddy/
├── python-version/          # Python implementation
│   ├── fastcaddy/          # Python package
│   ├── nbs/                # Jupyter notebooks (source)
│   ├── examples/           # Python examples
│   ├── setup.py            # Python package setup
│   └── ...                 # Python-specific files
├── examples/               # Go examples
├── internal/               # Go internal packages
├── cmd/                    # Go command line tools
├── pkg/                    # Go public packages
├── fastcaddy.go           # Main Go library file
├── go.mod                 # Go module definition
└── ...                    # Go-specific files
```

## Quick Start

### Python Version

1. **Installation**
   ```bash
   cd python-version
   pip install fastcore httpx
   pip install -e .
   ```

2. **Basic Usage**
   ```python
   from fastcaddy.core import *
   
   # Initialize FastCaddy
   setup_caddy(local=True, install_trust=True)
   
   # Add a reverse proxy
   add_reverse_proxy("api.example.com", "localhost:8080")
   
   # Check if domain is configured
   if has_id("api.example.com"):
       print("Domain configured successfully")
   
   # Delete domain
   del_id("api.example.com")
   ```

### Go Version

1. **Installation**
   ```bash
   go mod tidy
   go build
   ```

2. **Basic Usage**
   ```go
   package main
   
   import "github.com/youfun/fastcaddy"
   
   func main() {
       fc := fastcaddy.New()
       
       // Initialize FastCaddy
       installTrust := true
       fc.SetupCaddy("", "srv0", true, &installTrust)
       
       // Add a reverse proxy
       fc.AddReverseProxy("api.example.com", "localhost:8080")
       
       // Check if domain is configured
       if fc.HasID("api.example.com") {
           println("Domain configured successfully")
       }
       
       // Delete domain
       fc.DeleteRoute("api.example.com")
   }
   ```

## Programming Interface Examples

Both versions provide comprehensive examples demonstrating advanced usage:

### Python Examples (`python-version/examples/`)
- **`advanced_example.py`** - Complete programming interface demonstration
- **`domain_management.py`** - Specialized domain management with best practices
- **`programming_interface_complete.ipynb`** - Interactive Jupyter notebook tutorial

### Go Examples (`examples/`)
- **`basic/main.go`** - Basic usage demonstration  
- **`advanced/main.go`** - Comprehensive programming interface demo
- **`domain-management/main.go`** - Specialized domain management examples

### Key Features Demonstrated

#### Domain Management
- **Check domain status**: `has_id(domain)` (Python) / `HasID(domain)` (Go)
- **Add domains**: `add_reverse_proxy(from, to)` (Python) / `AddReverseProxy(from, to)` (Go)  
- **Delete domains**: `del_id(domain)` (Python) / `DeleteRoute(domain)` (Go)
- **Batch operations**: Multiple domain handling patterns

#### Wildcard Domains
- **Add wildcard routes**: `add_wildcard_route(domain)` (Python) / `AddWildcardRoute(domain)` (Go)
- **Add subdomains**: `add_sub_reverse_proxy(domain, sub, port)` (Python) / `AddSubReverseProxy(domain, sub, ports, host)` (Go)

#### Configuration Management
- **Check paths**: `has_path(path)` (Python) / `HasPath(path)` (Go)
- **Get config**: `gcfg(path)` (Python) / `GetConfig(path)` (Go)
- **System status**: Environment and configuration validation

## Best Practices

### 1. Always Check Before Acting
```python
# Python
if not has_id("example.com"):
    add_reverse_proxy("example.com", "localhost:8080")
```

```go
// Go
if !fc.HasID("example.com") {
    fc.AddReverseProxy("example.com", "localhost:8080")
}
```

### 2. Verify Operations
```python
# Python
del_id("example.com")
if not has_id("example.com"):
    print("Successfully deleted")
```

```go
// Go
fc.DeleteRoute("example.com")
if !fc.HasID("example.com") {
    fmt.Println("Successfully deleted")
}
```

### 3. Handle Errors Gracefully
```python
# Python
try:
    add_reverse_proxy("example.com", "localhost:8080")
except Exception as e:
    print(f"Failed: {e}")
```

```go
// Go
if err := fc.AddReverseProxy("example.com", "localhost:8080"); err != nil {
    log.Printf("Failed: %v", err)
}
```

## Prerequisites

### Common Requirements
- Caddy server running with Admin API enabled (typically `localhost:2019`)
- Network access to Caddy Admin API

### Python Specific
- Python 3.8+
- Dependencies: `fastcore`, `httpx`

### Go Specific  
- Go 1.24+
- Module dependencies (auto-installed with `go mod tidy`)

## Documentation

- **Python Examples**: See `python-version/examples/README.md`
- **Go Examples**: See `examples/README.md`
- **Original Python Documentation**: Available in `python-version/nbs/` notebooks

## Cloudflare Setup (for SSL)

FastCaddy can integrate with Cloudflare for automatic SSL certificate management. See the notebooks in `python-version/nbs/` for detailed Cloudflare token setup instructions.

## Contributing

When contributing:
1. Python changes go in `python-version/`
2. Go changes go in the root directory and related subdirectories
3. Add examples for new features in the appropriate `examples/` directory
4. Update documentation to reflect changes

## Migration Guide

If you were using the previous version that had Python files in the root:
1. Python code is now in `python-version/`
2. Import paths remain the same: `from fastcaddy.core import *`
3. Installation: `cd python-version && pip install -e .`

The Go implementation provides equivalent functionality with Go idioms and patterns.

## License

[MIT License](LICENSE)