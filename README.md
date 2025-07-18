<p align="center">
  <img src="https://github.com/dunamismax/go-web/blob/main/docs/images/go-logo.png" alt="Go Standard Library Web Stack Logo" width="400" />
</p>

<p align="center">
  <a href="https://github.com/dunamismax/go-stdlib">
    <img src="https://readme-typing-svg.demolab.com/?font=Fira+Code&size=24&pause=1000&color=00ADD8&center=true&vCenter=true&width=800&lines=The+Ultimate+Go+Standard+Library+Web+Stack;Go+%2B+http.ServeMux+%2B+HTMX+%2B+html/template;Single-Binary+Deployment;SQLite+%2B+Vanilla+CSS;Mage+Build+System" alt="Typing SVG" />
  </a>
</p>

<p align="center">
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8.svg?logo=go" alt="Go Version"></a>
  <a href="https://pkg.go.dev/net/http"><img src="https://img.shields.io/badge/Router-http.ServeMux-00ADD8.svg?logo=go" alt="Standard Library Router"></a>
  <a href="https://htmx.org/"><img src="https://img.shields.io/badge/HTMX-2.0+-3366CC.svg?logo=htmx" alt="HTMX Version"></a>
  <a href="https://pkg.go.dev/html/template"><img src="https://img.shields.io/badge/Templates-html/template-00ADD8.svg?logo=go" alt="Standard Library Templates"></a>
  <a href="https://sqlite.org/"><img src="https://img.shields.io/badge/SQLite-3.0+-003B57.svg?logo=sqlite" alt="SQLite Version"></a>
  <a href="https://magefile.org/"><img src="https://img.shields.io/badge/Mage-1.15+-purple.svg?logo=go" alt="Mage Version"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="MIT License"></a>
</p>

---

## About

A monorepo showcasing **The Ultimate Go Standard Library Web Stack** - built for maximum simplicity and stability using only Go's standard library. Produces a single, self-contained binary with embedded assets and hypermedia-driven interfaces.

**Core Philosophy:**

- Standard Library First with minimal external dependencies
- Single-Binary Deployment with embedded assets
- Type-Safe Templates with XSS protection
- Zero-Latency SQLite database (CGO-free)
- Vanilla CSS without framework overhead

## Tech Stack

| Layer         | Technology                                                          | Purpose                                 |
| ------------- | ------------------------------------------------------------------- | --------------------------------------- |
| **Backend**   | [Go](https://go.dev/doc/) + [net/http](https://pkg.go.dev/net/http) | HTTP server with method-aware routing   |
| **Database**  | [SQLite](https://www.sqlite.org/docs.html)                          | Embedded database (CGO-free)            |
| **Frontend**  | [HTMX](https://htmx.org/docs/)                                      | Dynamic interactions without JavaScript |
| **Templates** | [html/template](https://pkg.go.dev/html/template)                   | Type-safe HTML with XSS protection      |
| **Styling**   | Vanilla CSS                                                         | Direct styling control                  |
| **Build**     | [Mage](https://magefile.org/)                                       | Go-based build automation               |

## Quick Start

1. Install Mage: `go install github.com/magefile/mage@latest`
2. Clone and initialize:

   ```bash
   git clone https://github.com/dunamismax/go-stdlib.git
   cd go-stdlib
   mage dev:init
   ```

3. Start applications:

   ```bash
   mage dev:start
   # API Playground: http://localhost:8080
   # GoSocial: http://localhost:8081
   ```

## Applications

### API Playground (Port 8080)

Interactive API testing platform with text analysis, random generators, hash/encoding tools, and time utilities.

### GoSocial (Port 8081)

Social media platform with secure authentication, real-time feeds, SQLite database, and responsive design.

<p align="center">
  <img src="https://github.com/dunamismax/go-web/blob/main/docs/images/gopher-mage.svg" alt="Gopher Mage" width="150" />
</p>

## Commands

```bash
# Development
mage dev:init        # Initialize environment
mage dev:start       # Start all applications
mage build:all       # Build all applications
mage test:all        # Run all tests

# Individual apps
mage runapi          # Run API playground only
mage runsocial       # Run GoSocial only

# Production
mage prod:release    # Create production release
mage prod:caddy      # Start Caddy reverse proxy
```

## Package Architecture

- **Database** (`pkg/database`): SQLite management with migrations and connection pooling
- **Middleware** (`pkg/middleware`): HTTP middleware for structured logging, CORS, and rate limiting
- **Utils** (`pkg/utils`): Response helpers, text processing, and random generation

## Production Deployment

Each application builds to a single binary containing the Go executable, static assets, HTML templates, and SQLite schema.

```bash
mage prod:release
./build/api-playground
./build/go-social
```

## Security & Performance

**Security:** Input validation, HTML escaping, parameterized queries, secure headers, bcrypt hashing, HTTP-only cookies

**Performance:** Compiled binary, embedded assets, zero network latency, minimal dependencies, goroutine concurrency

## Contributing

1. Fork the repository
2. Create a feature branch
3. Run tests: `mage test:all`
4. Format code: `mage dev:fmt`
5. Submit a pull request

<p align="center">
  <a href="https://buymeacoffee.com/dunamismax" target="_blank">
    <img src="https://github.com/dunamismax/go-web/blob/main/docs/images/buy-coffee-go.gif" alt="Buy Me A Coffee" style="height: 150px !important;" />
  </a>
</p>

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

<p align="center">
  <img src="https://github.com/dunamismax/go-web/blob/main/docs/images/gopher-running-jumping.gif" alt="Gopher Running and Jumping" width="400" />
</p>
