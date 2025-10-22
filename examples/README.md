# FastCaddy Examples

This directory contains comprehensive examples demonstrating the FastCaddy programming interface.

## Directory Structure

```
examples/
├── advanced/                    # Advanced examples
│   └── main.go                 # Comprehensive programming interface demo
├── basic/                      # Basic examples
│   └── main.go                 # Basic usage demonstration
├── domain-management/          # Domain management examples
│   └── main.go                 # Domain management specialized demo
└── test-utils/                 # Test utilities
    └── main.go                 # Testing utilities
```

## Examples

### 1. Basic Examples (`examples/basic/main.go`)
- Environment setup
- Adding reverse proxies
- Adding wildcard routes
- Adding subdomain reverse proxies
- Basic configuration checking

### 2. Advanced Examples (`examples/advanced/main.go`)
- Batch domain operations
- Domain status checking
- Configuration management
- Deletion and cleanup operations
- Complete workflow demonstrations

### 3. Domain Management (`examples/domain-management/main.go`)
- Safe domain addition with conflict checking
- Safe domain deletion with verification
- Domain configuration updates
- Batch domain management
- Wildcard domain handling
- Best practices demonstration

## Key Programming Interface Features Demonstrated

### Domain Status Checking
- `fc.HasID(domain)`, `fc.HasPath(path)`

### Domain Configuration
- `fc.AddReverseProxy()`, `fc.AddWildcardRoute()`, `fc.AddSubReverseProxy()`

### Domain Deletion
- `fc.DeleteRoute(id)`

### Configuration Management
- `fc.GetConfig(path)`, `fc.PutConfig()`

## Running the Examples

```bash
# Basic example
cd examples/basic
go run main.go

# Advanced example
cd examples/advanced
go run main.go

# Domain management example
cd examples/domain-management
go run main.go
```

## Prerequisites

- Go 1.24+ installed
- Caddy server running on `localhost:2019`
- FastCaddy Go module available

## Common Use Cases Covered

1. **Domain Existence Checking**
   - Check if a domain is already configured before adding
   - Verify deletion success after removing domains

2. **Safe Domain Operations**
   - Add domains with conflict resolution
   - Update existing domain configurations
   - Delete domains with verification

3. **Batch Operations**
   - Add multiple domains efficiently
   - Check status of multiple domains
   - Delete multiple domains safely

4. **Wildcard Domain Management**
   - Set up wildcard routes for subdomains
   - Add specific subdomain configurations
   - Manage subdomain lifecycle

5. **Configuration Validation**
   - Check system component status
   - Verify configuration integrity
   - Troubleshoot configuration issues

6. **Production Scenarios**
   - Microservices domain setup
   - Development environment switching
   - Deployment verification workflows

## Best Practices Demonstrated

1. **Always Check Before Acting**
   - Verify domain status before operations
   - Check system prerequisites

2. **Handle Errors Gracefully**
   - Proper error handling and reporting
   - Fallback strategies for failures

3. **Verify Operations**
   - Confirm success after modifications
   - Validate configuration integrity

4. **Use Batch Operations**
   - Efficient handling of multiple domains
   - Consistent error handling across operations

5. **Clean Up Resources**
   - Remove test configurations
   - Avoid configuration pollution

6. **Document Configurations**
   - Clear naming conventions
   - Comprehensive logging and feedback

## Troubleshooting

### Common Issues

1. **Caddy Server Not Running**
   - Ensure Caddy is running on `localhost:2019`
   - Check Caddy admin API is accessible

2. **Permission Issues**
   - Verify network access to Caddy admin API
   - Check firewall settings

3. **Configuration Conflicts**
   - Use force update options when appropriate
   - Clean up conflicting configurations

4. **Import Errors (Python)**
   - Ensure Python path includes FastCaddy module
   - Install required dependencies

For more help, see the main FastCaddy documentation or raise an issue in the repository.