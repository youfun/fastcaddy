# FastCaddy

FastCaddy is a Go library and command-line tool for managing Caddy server configurations via the Caddy Admin API. It provides an easy-to-use interface for setting up reverse proxies, managing TLS certificates, and handling wildcard domains.

## Features

- 🚀 **High Performance**: Written in Go for optimal performance
- 🖥️ **Command-Line Interface**: Easy-to-use CLI for common operations
- 📚 **Programming Interface**: Use as a Go library in your applications
- 🔐 **TLS Support**: ACME (Let's Encrypt) and internal certificate configuration
- 🌐 **Route Management**: Complete reverse proxy and wildcard domain support
- 🛠️ **Flexible Configuration**: Direct access to Caddy's Admin API

## Installation

### From Source

```bash
git clone https://github.com/youfun/fastcaddy.git
cd fastcaddy
go build -o fastcaddy ./cmd/fastcaddy
```

### Using Go Install

```bash
go install github.com/youfun/fastcaddy/cmd/fastcaddy@latest
```

## Command-Line Usage

### Setup Caddy Configuration

#### Local Development (Using Internal Certificates)
```bash
./fastcaddy setup --local
```

#### Production (Using Let's Encrypt + Cloudflare)
```bash
export CADDY_CF_TOKEN="your-cloudflare-token"
./fastcaddy setup --cf-token $CADDY_CF_TOKEN
```

#### Install Root Certificate to System Trust Store
```bash
./fastcaddy setup --local --install-trust
```

### Manage Reverse Proxies

#### Add Simple Reverse Proxy
```bash
./fastcaddy add-proxy --from api.example.com --to localhost:8080
./fastcaddy add-proxy --from web.example.com --to 127.0.0.1:3000
```

#### Delete Reverse Proxy
```bash
./fastcaddy del-proxy --id api.example.com
```

### Wildcard Subdomain Support

#### Add Wildcard Domain
```bash
./fastcaddy add-wildcard --domain example.com
```

#### Add Subdomain Reverse Proxy
```bash
# Single port
./fastcaddy add-sub-proxy --domain example.com --subdomain api --ports 8080

# Multiple ports
./fastcaddy add-sub-proxy --domain example.com --subdomain web --ports 3000,3001

# Specify target host
./fastcaddy add-sub-proxy --domain example.com --subdomain db --ports 5432 --host 192.168.1.10
```

### Check Status
```bash
./fastcaddy status
```

## Programming Interface

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/youfun/fastcaddy"
)

func main() {
    // Create FastCaddy client
    fc := fastcaddy.New()
    
    // Setup local development environment
    err := fc.SetupCaddy("", "srv0", true, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Add reverse proxy
    err = fc.AddReverseProxy("api.localhost", "localhost:8080")
    if err != nil {
        log.Fatal(err)
    }
    
    // Add wildcard domain
    err = fc.AddWildcardRoute("localhost")
    if err != nil {
        log.Fatal(err)
    }
    
    // Add subdomain reverse proxy
    err = fc.AddSubReverseProxy("localhost", "web", []string{"3000"}, "localhost")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Caddy configuration completed!")
}
```

### Advanced Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/youfun/fastcaddy"
    "github.com/youfun/fastcaddy/pkg/types"
)

func main() {
    fc := fastcaddy.New()
    
    // Directly manipulate routes
    route := types.Route{
        ID: "custom-route",
        Match: []types.RouteMatch{
            {Host: []string{"custom.example.com"}},
        },
        Handle: []types.Handler{
            {
                Handler: "reverse_proxy",
                Upstreams: []types.Upstream{
                    {Dial: "backend1:8080"},
                    {Dial: "backend2:8080"},
                },
            },
        },
        Terminal: true,
    }
    
    err := fc.Routes.AddRoute(route)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check configuration status
    if fc.HasPath("/apps/http/servers") {
        fmt.Println("HTTP server configured")
    }
    
    // Get configuration
    config, err := fc.GetConfig("/apps/http/servers/srv0")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server configuration: %+v\n", config)
}
```

## Installing Caddy

This project helps you use the Caddy API rather than a Caddyfile. To use the API with automatic HTTPS, you need to install a plugin for your domain management service. We use Cloudflare, so we'll document that here. For other domain services, see the Caddy docs for other plugins.

### Installing Caddy with Cloudflare DNS Plugin

To install caddy, we'll use a tool called `xcaddy`. First install Go:

- Mac: `brew install go`
- Linux: `sudo apt install golang`

Note: If you are not on the latest Ubuntu, you'll need to setup the backport repo before installing go:

```sh
sudo add-apt-repository -y ppa:longsleep/golang-backports
sudo apt update
```

Now install xcaddy:

```sh
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

Then compile caddy with the Cloudflare plugin:

```sh
mkdir -p ~/go/bin
cd ~/go/bin
./xcaddy build --with github.com/caddy-dns/cloudflare
```

This gives you a `~/go/bin/caddy` binary:

```sh
./caddy version
./caddy run
```

### Run Caddy Securely on Start

If you're using a server or running caddy frequently, you'll want it to run on start. Run from this repo root:

```sh
./setup_service.sh
```

If all went well, you should see output like this:

```sh
● caddy.service - Caddy
     Loaded: loaded (/etc/systemd/system/caddy.service; enabled; preset: enabled)
     Active: active (running) since Sat 2024-11-09 05:06:47 UTC; 2 days ago
       Docs: https://caddyserver.com/docs/
   Main PID: 138140 (caddy)
      Tasks: 29 (limit: 154166)
     Memory: 19.3M (peak: 28.8M)
        CPU: 3min 37.216s
     CGroup: /system.slice/caddy.service
             └─138140 /usr/bin/caddy run --environ
```

## Project Structure

```
fastcaddy/
├── cmd/
│   └── fastcaddy/          # Command-line tool
├── internal/
│   ├── api/               # Caddy API interaction
│   ├── config/            # Configuration management
│   ├── tls/               # TLS configuration
│   ├── routes/            # Route management
│   └── utils/             # Utility functions
├── pkg/
│   └── types/             # Public type definitions
├── examples/              # Example applications
├── fastcaddy.go           # Main client interface
├── go.mod                 # Go module definition
└── go.sum                 # Dependency checksums
```

## Environment Variables

- `CADDY_CF_TOKEN`: Cloudflare API token
- `CLOUDFLARE_API_TOKEN`: Alternative Cloudflare API token

## Examples

See the `examples/` directory for comprehensive examples:

- **basic/** - Basic usage demonstration
- **advanced/** - Advanced programming interface
- **domain-management/** - Domain management patterns

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go conventions
2. Add appropriate tests
3. Update documentation

## License

See [LICENSE](LICENSE) file for details.
