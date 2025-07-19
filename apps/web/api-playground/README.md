# API Playground

Interactive developer tools built with the Go Echo TypeScript Stack - a high-performance web application combining Echo framework, TypeScript, Vite, and Pico.css.

## Tech Stack

- **Backend Framework**: [Echo](https://echo.labstack.com/) - High-performance, extensible web framework for Go
- **Frontend Language**: [TypeScript](https://www.typescriptlang.org/) - Type-safe client-side development
- **Build Tool**: [Vite](https://vitejs.dev/) - Lightning-fast frontend build tool with HMR
- **Styling**: [Pico.css](https://picocss.com/) - Minimal CSS framework
- **Asset Embedding**: Go embed for single-binary deployment
- **Database**: [SQLite](https://www.sqlite.org/) (via modernc.org/sqlite) - CGO-free embedded database

## Features

### Text Analysis & Manipulation
- **Text Analysis**: Character, word, sentence, and paragraph counting
- **Text Transformation**: Reverse text, title case conversion, vowel counting

### Random Generators
- **Cryptographically Secure Random Numbers**: Range-based number generation
- **Random String Generation**: Alphanumeric string creation
- **UUID Generation**: RFC 4122 compliant UUID v4

### Encoding & Hashing
- **SHA256 Hashing**: Secure hash generation
- **Base64 Encoding/Decoding**: Text encoding utilities

### Fun APIs
- **Programming Jokes**: Curated collection of developer humor
- **Inspirational Quotes**: Motivational quotes for developers
- **Tech Facts**: Interesting technology trivia

### Time & Date Utilities
- **Current Time**: Multi-timezone time display
- **Unix Timestamps**: Epoch time conversion
- **Date Formats**: Multiple date/time format examples

## Architecture

### Echo Framework Integration
- **Template Renderer**: Custom Echo renderer for HTML templates
- **Middleware Stack**: Logger, recovery, CORS, security, and request ID
- **Static File Serving**: Embedded assets served via Echo's StaticFS
- **Structured Logging**: Comprehensive request/response logging

### Frontend Pipeline
- **TypeScript**: Type-safe client-side development
- **Vite Build**: Optimized bundling with code splitting
- **Hot Module Replacement**: Instant development feedback
- **Pico.css**: Semantic CSS with dark theme support

### Production Deployment
- **Single Binary**: All assets embedded using go:embed
- **Zero Dependencies**: Self-contained executable
- **High Performance**: Sub-millisecond response times

## Development

### Prerequisites
- Go 1.24+
- Node.js 18+
- npm or yarn

### Setup
```bash
# Install Go dependencies
go mod tidy

# Install frontend dependencies
npm install

# Build frontend assets
npm run build

# Build application
go build -o api-playground

# Run application
./api-playground
```

### Development Workflow

#### Backend Development with Live Reload
```bash
# Using Mage (recommended)
mage runapiwithair

# Manual with Air
air
```

#### Frontend Development with HMR
```bash
# Terminal 1: Start Go backend
./api-playground

# Terminal 2: Start Vite dev server
npm run dev
# Access via http://localhost:3000 (proxies to backend at :8080)
```

#### Full-Stack Development
```bash
# Terminal 1: Backend with live reload
mage runapiwithair

# Terminal 2: Frontend with HMR
mage runapifrontenddev
```

## API Endpoints

### Text Operations
- `POST /analyze` - Text analysis (characters, words, sentences, paragraphs)
- `POST /reverse` - Reverse text
- `POST /titlecase` - Convert to title case
- `POST /count-vowels` - Count vowels in text

### Generators
- `POST /random-number` - Generate random number in range
- `POST /random-string` - Generate random alphanumeric string
- `POST /uuid` - Generate UUID v4

### Encoding
- `POST /sha256` - Generate SHA256 hash
- `POST /base64-encode` - Base64 encode text
- `POST /base64-decode` - Base64 decode text

### Fun APIs
- `GET /joke` - Random programming joke
- `GET /quote` - Random inspirational quote
- `GET /fact` - Random tech fact

### Time & Date
- `POST /current-time` - Current time in specified timezone
- `POST /timestamp` - Unix timestamp
- `GET /date-formats` - Various date format examples

## Configuration

### Environment Variables
- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Logging level (default: info)

### Build Configuration
- `vite.config.ts` - Frontend build configuration
- `tsconfig.json` - TypeScript compiler options
- `package.json` - NPM dependencies and scripts

## Performance

- **Response Time**: Sub-millisecond for most operations
- **Memory Usage**: ~15MB baseline
- **Binary Size**: ~12MB (includes all assets)
- **Concurrent Users**: 10,000+ on modest hardware

## Security Features

- **CORS Protection**: Configurable cross-origin resource sharing
- **Security Headers**: HSTS, CSP, and other security headers
- **Input Validation**: Server-side validation and sanitization
- **XSS Protection**: Automatic escaping via html/template
- **Secure Random**: Cryptographically secure random generation

## Deployment

### Single Binary
```bash
# Build for production
npm run build
go build -ldflags="-s -w" -o api-playground

# Deploy binary
./api-playground
```

### Docker
```dockerfile
FROM scratch
COPY api-playground /
EXPOSE 8080
CMD ["/api-playground"]
```

### Reverse Proxy (Caddy)
```
api-playground.example.com {
    reverse_proxy localhost:8080
}
```

## License

MIT License - see root repository LICENSE file.