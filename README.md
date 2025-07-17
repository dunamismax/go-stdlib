<p align="center">
  <img src="go.png" alt="Go Standard Library Web Stack Logo" width="600" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-stdlib">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=800&lines=The+Ultimate+Go+Standard+Library+Web+Stack;Go+%2B+http.ServeMux+%2B+HTMX+%2B+html/template;Single-Binary+Deployment;SQLite+%2B+Vanilla+CSS;Mage+Build+System" alt="Typing SVG" />
  </a>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.23+-00ADD8.svg?logo=go" alt="Go Version"></a>
  <a href="https://pkg.go.dev/net/http"><img src="https://img.shields.io/badge/Router-http.ServeMux-00ADD8.svg?logo=go" alt="Standard Library Router"></a>
  <a href="https://htmx.org/"><img src="https://img.shields.io/badge/HTMX-2.0+-3366CC.svg?logo=htmx" alt="HTMX Version"></a>
  <a href="https://pkg.go.dev/html/template"><img src="https://img.shields.io/badge/Templates-html/template-00ADD8.svg?logo=go" alt="Standard Library Templates"></a>
  <a href="https://sqlite.org/"><img src="https://img.shields.io/badge/SQLite-3.0+-003B57.svg?logo=sqlite" alt="SQLite Version"></a>
  <a href="https://magefile.org/"><img src="https://img.shields.io/badge/Mage-1.15+-purple.svg?logo=go" alt="Mage Version"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License"></a>
</p>

---

## About This Stack

This monorepo showcases **The Ultimate Go Standard Library Web Stack** - architected for maximum simplicity, stability, and robustness by exclusively using the Go standard library for all backend logic. It produces a single, self-contained, and dependency-free binary that serves a rich, hypermedia-driven interface, representing the purest form of a Go web application.

### Core Philosophy

- **Standard Library First**: Pure Go standard library with minimal external dependencies
- **Hypermedia-Driven**: Server-rendered HTML with HTMX for dynamic interactions
- **Single-Binary Deployment**: Everything embedded via `//go:embed`
- **Type-Safe Templates**: Go's built-in `html/template` with automatic XSS protection
- **Zero-Latency Database**: SQLite embedded directly in the binary (CGO-free)
- **No Build Steps**: Pure Go with embedded assets
- **Vanilla CSS**: Direct styling control without framework overhead

## Tech Stack

| Layer                    | Technology                                                      | Purpose                                          |
| ------------------------ | --------------------------------------------------------------- | ------------------------------------------------ |
| **Backend Language**     | [Go](https://go.dev/doc/)                                       | Fast, compiled, statically-typed backend         |
| **Asset Embedding**      | [Go Embed](https://pkg.go.dev/embed)                            | Bundle all static assets into binary             |
| **Backend Router**       | [net/http](https://pkg.go.dev/net/http)                         | Standard library HTTP router (http.ServeMux)    |
| **Database**             | [SQLite](https://www.sqlite.org/docs.html)                      | Zero-latency embedded database (CGO-free)       |
| **Frontend Interaction** | [HTMX](https://htmx.org/docs/)                                  | Dynamic interfaces without JavaScript frameworks |
| **Markup & Templates**   | [html/template](https://pkg.go.dev/html/template)               | Type-safe HTML templates with XSS protection    |
| **Styling**              | [Vanilla CSS](https://developer.mozilla.org/en-US/docs/Web/CSS) | Direct styling control, embedded via Go          |
| **Build Orchestrator**   | [Mage](https://magefile.org/)                                   | Go-based build automation                        |

## Quick Start

### Prerequisites

- Go 1.23 or higher
- [Mage](https://magefile.org/) for build automation

### Installation

1. Install Mage:

   ```bash
   go install github.com/magefile/mage@latest
   ```

2. Clone and initialize:

   ```bash
   git clone https://github.com/dunamismax/go-stdlib.git
   cd go-stdlib
   mage dev:init
   ```

3. Start the stack:

   ```bash
   mage dev:start
   # API Playground: http://localhost:8080
   # GoSocial: http://localhost:8081
   ```

## Applications

### API Playground (Port 8080)

Interactive API testing platform showcasing Go + HTMX capabilities.

**Features:**

- Text analysis with word frequency and reading time
- Random generators for numbers, strings, dice, and UUIDs
- Hash and encoding tools (SHA256, Base64, URL encoding)
- Time utilities with multiple formats

**Implementation:**

- Pure Go standard library with `http.ServeMux`
- Go's `html/template` for safe HTML rendering
- Embedded CSS styling via `//go:embed`
- HTMX for dynamic interactions
- Type-safe form handling

### GoSocial (Port 8081)

Social media platform demonstrating full-stack Go development using only standard library.

**Features:**

- Secure bcrypt-based authentication
- Real-time social feed with like/unlike functionality
- Embedded SQLite database with proper migrations (CGO-free)
- Responsive design with dark theme
- XSS protection via `html/template`

**Implementation:**

- Go standard library with `http.ServeMux`
- SQLite with CGO-free driver (`modernc.org/sqlite`)
- HTMX for dynamic interactions
- `html/template` for type-safe HTML generation
- Vanilla CSS embedded via `//go:embed`
- HTTP-only cookies for session management

## Development Commands

### Essential Commands

```bash
mage dev:init        # Initialize development environment
mage dev:start       # Start all applications
mage build:all       # Build all applications
mage test:all        # Run all tests
mage status          # Show project status
mage clean           # Clean build artifacts
```

### Build Commands

```bash
mage build:web       # Build web applications
mage build:api       # Build API playground
mage build:social    # Build GoSocial
```

### Development Tools

```bash
mage dev:fmt         # Format all Go code
mage dev:lint        # Run linter on all code
mage dev:deps        # Download all dependencies
mage dev:tidy        # Tidy all Go modules
```

### Production Commands

```bash
mage prod:release    # Create production release
mage prod:caddy      # Start Caddy reverse proxy
```

### Individual Applications

```bash
mage runapi          # Run API playground only
mage runsocial       # Run GoSocial only
```

## Package Architecture

### Database (`pkg/database`)

SQLite management with embedded migrations, connection pooling, and query optimization using a CGO-free driver.

### Middleware (`pkg/middleware`)

Basic HTTP middleware for logging, CORS, timeouts, and error recovery using standard library patterns.

### Utils (`pkg/utils`)

Common utilities including response helpers, text processing, and random generation.

## Production Deployment

Each application builds to a single binary containing:

- Go executable
- All static assets (CSS, JS, images)
- HTML templates
- SQLite database schema

```bash
mage prod:release    # Build production binaries
./build/api-playground  # Run API playground
./build/go-social       # Run GoSocial
```

### Reverse Proxy

Use Caddy for automatic HTTPS and reverse proxying:

```bash
mage dev:start       # Start applications
mage prod:caddy      # Start Caddy proxy
```

## Security Features

- Input validation and sanitization
- Context-aware HTML escaping via `html/template`
- Parameterized SQL queries preventing SQL injection
- Secure headers (CORS, content-type)
- Bcrypt password hashing
- HTTP-only cookies with secure expiration
- Path validation for file operations

## Performance Benefits

- Compiled binary with no runtime interpretation overhead
- Embedded assets eliminate filesystem lookups
- SQLite provides zero network latency
- Minimal dependencies for faster startup
- Efficient connection pooling
- Go's goroutines handle thousands of concurrent connections
- Standard library optimizations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `mage test:all`
5. Format code: `mage dev:fmt`
6. Submit a pull request

## Troubleshooting

### Common Issues

**Build Problems:**

```bash
mage clean
mage build:all
```

**Development Issues:**

```bash
mage dev:init
mage dev:deps
mage dev:tidy
```

**Database Issues:**

```bash
# Check SQLite files
ls -la apps/web/go-social/data/
# Reset database (removes all data)
rm apps/web/go-social/data/*.db
```

## Support This Project

If you find this Ultimate Go Standard Library Web Stack valuable, consider supporting its development:

<p align="center">
  <a href="https://www.buymeacoffee.com/dunamismax" target="_blank">
    <img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" />
  </a>
</p>

## Connect

<p align="center">
  <a href="https://twitter.com/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Twitter-%231DA1F2.svg?&style=for-the-badge&logo=twitter&logoColor=white" alt="Twitter"></a>
  <a href="https://bsky.app/profile/dunamismax.bsky.social" target="_blank"><img src="https://img.shields.io/badge/Bluesky-blue?style=for-the-badge&logo=bluesky&logoColor=white" alt="Bluesky"></a>
  <a href="https://reddit.com/user/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Reddit-%23FF4500.svg?&style=for-the-badge&logo=reddit&logoColor=white" alt="Reddit"></a>
  <a href="https://discord.com/users/dunamismax" target="_blank"><img src="https://img.shields.io/badge/Discord-dunamismax-7289DA.svg?style=for-the-badge&logo=discord&logoColor=white" alt="Discord"></a>
  <a href="https://signal.me/#p/+dunamismax.66" target="_blank"><img src="https://img.shields.io/badge/Signal-dunamismax.66-3A76F0.svg?style=for-the-badge&logo=signal&logoColor=white" alt="Signal"></a>
</p>

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <strong>The Ultimate Go Standard Library Web Stack</strong><br>
  <sub>Go + http.ServeMux + HTMX + html/template + SQLite + Vanilla CSS + Mage</sub>
</p>