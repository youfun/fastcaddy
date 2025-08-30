# FastCaddy Python Examples

This directory contains comprehensive Python examples demonstrating the FastCaddy programming interface.

## Files

- **`advanced_example.py`** - Complete programming interface demonstration
- **`domain_management.py`** - Specialized domain management with object-oriented approach
- **`programming_interface_complete.ipynb`** - Interactive Jupyter notebook with step-by-step examples

## Quick Start

### Prerequisites

1. Caddy server running on `localhost:2019`
2. Python 3.8+ installed
3. Required packages:
   ```bash
   pip install fastcore httpx
   ```

### Running Examples

```bash
# Advanced programming interface example
python3 advanced_example.py

# Domain management specialized example  
python3 domain_management.py

# Interactive notebook (requires Jupyter)
jupyter notebook programming_interface_complete.ipynb
```

## Key Functions Demonstrated

### Core API Functions
- `setup_caddy()` - Initialize FastCaddy environment
- `has_id(domain)` - Check if domain is configured
- `has_path(path)` - Check if configuration path exists
- `add_reverse_proxy(from, to)` - Add reverse proxy
- `add_wildcard_route(domain)` - Add wildcard domain
- `add_sub_reverse_proxy(domain, sub, port)` - Add subdomain
- `del_id(id)` - Delete domain configuration
- `gcfg(path)` - Get configuration data

### Advanced Patterns
- Safe domain addition with conflict resolution
- Batch domain operations
- Configuration validation and verification
- Error handling and recovery
- System status checking

## Example Scenarios

### 1. Basic Domain Management
```python
# Check if domain exists
if has_id("api.example.com"):
    print("Domain already configured")

# Add reverse proxy
add_reverse_proxy("api.example.com", "localhost:8080")

# Verify addition
if has_id("api.example.com"):
    print("Domain successfully added")

# Delete domain
del_id("api.example.com")
```

### 2. Wildcard Domain Setup
```python
# Add wildcard route
add_wildcard_route("dev.local")

# Add subdomains
add_sub_reverse_proxy("dev.local", "app", "8080")
add_sub_reverse_proxy("dev.local", "api", "8081")

# Check subdomain status
if has_id("app.dev.local"):
    print("Subdomain configured successfully")
```

### 3. Batch Operations
```python
domains = {
    "api.example.com": "localhost:8080",
    "web.example.com": "localhost:3000",
    "admin.example.com": "localhost:9000"
}

# Add all domains
for domain, target in domains.items():
    if not has_id(domain):
        add_reverse_proxy(domain, target)
        print(f"Added {domain} -> {target}")
```

### 4. Configuration Validation
```python
# Check system status
if has_path("/apps/http"):
    print("HTTP app configured")
if has_path("/apps/tls"):
    print("TLS app configured")

# Get detailed configuration
config = gcfg("/apps/http/servers")
print(f"HTTP servers: {len(config)}")
```

## Best Practices

### 1. Always Check Before Acting
```python
def safe_add_domain(domain, target):
    if has_id(domain):
        print(f"Domain {domain} already exists")
        return False
    
    try:
        add_reverse_proxy(domain, target)
        return True
    except Exception as e:
        print(f"Failed to add domain: {e}")
        return False
```

### 2. Verify Operations
```python
def verified_delete(domain):
    if not has_id(domain):
        print(f"Domain {domain} doesn't exist")
        return True
    
    del_id(domain)
    
    # Verify deletion
    if has_id(domain):
        print(f"Failed to delete {domain}")
        return False
    
    print(f"Successfully deleted {domain}")
    return True
```

### 3. Handle Errors Gracefully
```python
def robust_domain_operations():
    try:
        setup_caddy(local=True, install_trust=True)
    except Exception as e:
        print(f"Setup failed: {e}")
        return False
    
    # Continue with operations...
    return True
```

## Troubleshooting

### Common Issues

1. **Import Error**: Ensure the FastCaddy module is in Python path
   ```python
   import sys
   sys.path.append('..')  # Adjust path as needed
   ```

2. **Connection Error**: Verify Caddy is running
   ```bash
   curl http://localhost:2019/config/
   ```

3. **Permission Error**: Check network access to Caddy admin API

### Debug Tips

- Use `gcfg("/")` to see complete configuration
- Check `has_path("/apps")` to verify basic setup
- Use try-except blocks around all API calls
- Add time delays between operations if needed

## Interactive Development

The Jupyter notebook (`programming_interface_complete.ipynb`) provides an interactive environment for:

- Learning the API step by step
- Experimenting with configurations
- Testing different scenarios
- Understanding the complete workflow

It's the recommended starting point for new users.

## Integration Examples

### Flask Web Application
```python
from flask import Flask, request, jsonify
from fastcaddy.core import *

app = Flask(__name__)

@app.route('/domains', methods=['POST'])
def add_domain():
    data = request.json
    domain = data['domain']
    target = data['target']
    
    if has_id(domain):
        return jsonify({'error': 'Domain already exists'}), 400
    
    try:
        add_reverse_proxy(domain, target)
        return jsonify({'success': True})
    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/domains/<domain>', methods=['DELETE'])
def delete_domain(domain):
    if not has_id(domain):
        return jsonify({'error': 'Domain not found'}), 404
    
    try:
        del_id(domain)
        return jsonify({'success': True})
    except Exception as e:
        return jsonify({'error': str(e)}), 500
```

### Configuration Management Script
```python
#!/usr/bin/env python3
import json
import argparse
from fastcaddy.core import *

def load_config(filename):
    with open(filename) as f:
        return json.load(f)

def apply_config(config):
    for domain, target in config.get('domains', {}).items():
        if not has_id(domain):
            add_reverse_proxy(domain, target)
            print(f"Added {domain}")

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('config', help='Configuration file')
    args = parser.parse_args()
    
    config = load_config(args.config)
    apply_config(config)
```

This provides a solid foundation for using FastCaddy Python API in various scenarios.