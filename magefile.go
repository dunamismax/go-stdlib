//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Dev.Start

// Project structure
const (
	rootDir          = "."
	appsDir          = "./apps"
	webDir           = "./apps/web"
	pkgDir           = "./pkg"
	buildDir         = "./build"
	scriptsDir       = "./scripts"
	apiPlaygroundDir = "./apps/web/api-playground"
	goSocialDir      = "./apps/web/go-social"
	goHyperDocsDir   = "./apps/web/gohyperdocs"
	componentsDir    = "./pkg/components"
	databaseDir      = "./pkg/database"
	middlewareDir    = "./pkg/middleware"
	stylesDir        = "./pkg/styles"
	utilsDir         = "./pkg/utils"
)

// Build outputs
const (
	apiPlaygroundBin = "./build/api-playground"
	goSocialBin      = "./build/go-social"
	goHyperDocsBin   = "./build/gohyperdocs"
)

// Colors for output
const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[0;33m"
	colorBlue   = "\033[0;34m"
	colorPurple = "\033[0;35m"
	colorCyan   = "\033[0;36m"
	colorWhite  = "\033[0;37m"
	colorReset  = "\033[0m"
)

// Helper functions
func colorPrint(color, message string) {
	fmt.Printf("%s%s%s\n", color, message, colorReset)
}

func printSection(title string) {
	colorPrint(colorCyan, fmt.Sprintf("🚀 %s", title))
}

func printSuccess(message string) {
	colorPrint(colorGreen, fmt.Sprintf("✅ %s", message))
}

func printError(message string) {
	colorPrint(colorRed, fmt.Sprintf("❌ %s", message))
}

func printWarning(message string) {
	colorPrint(colorYellow, fmt.Sprintf("⚠️  %s", message))
}

// Development namespace
type Dev mg.Namespace

// Build namespace
type Build mg.Namespace

// Test namespace
type Test mg.Namespace

// Tools namespace
type Tools mg.Namespace

// Production namespace
type Prod mg.Namespace

// Help displays available targets
func Help() {
	colorPrint(colorCyan, "Go Standard Library Web Stack - The Ultimate Hypermedia-Driven Web Stack")
	fmt.Println("===========================================================================")
	fmt.Println()
	colorPrint(colorBlue, "Development Commands:")
	fmt.Println("  mage dev:init        Initialize development environment")
	fmt.Println("  mage dev:start       Start development environment")
	fmt.Println("  mage dev:startwithair Start with Air live reloading")
	fmt.Println("  mage dev:deps        Download all dependencies")
	fmt.Println("  mage dev:tidy        Tidy all Go modules")
	fmt.Println("  mage dev:fmt         Format all Go code")
	fmt.Println("  mage dev:lint        Run linter on all code")
	fmt.Println()
	colorPrint(colorBlue, "Build Commands:")
	fmt.Println("  mage build:all       Build all applications")
	fmt.Println("  mage build:web       Build web applications")
	fmt.Println("  mage build:api       Build API playground")
	fmt.Println("  mage build:social    Build GoSocial")
	fmt.Println("  mage build:docs      Build GoHyperDocs")
	fmt.Println()
	colorPrint(colorBlue, "Test Commands:")
	fmt.Println("  mage test:all        Run all tests")
	fmt.Println("  mage test:quick      Quick test (no cache)")
	fmt.Println("  mage test:security   Run security checks")
	fmt.Println()
	colorPrint(colorBlue, "Production Commands:")
	fmt.Println("  mage prod:caddy      Start Caddy reverse proxy")
	fmt.Println("  mage prod:release    Create production release")
	fmt.Println()
	colorPrint(colorBlue, "Maintenance Commands:")
	fmt.Println("  mage clean           Clean build artifacts")
	fmt.Println("  mage format          Format all Go code and tidy modules")
	fmt.Println("  mage tools:install   Install development tools")
	fmt.Println("  mage tools:upgrade   Upgrade all dependencies")
	fmt.Println("  mage status          Show project status")
	fmt.Println()
	colorPrint(colorBlue, "Individual App Commands:")
	fmt.Println("  mage runapi          Run API playground only")
	fmt.Println("  mage runsocial       Run GoSocial only")
	fmt.Println("  mage rundocs         Run GoHyperDocs only")
	fmt.Println("  mage runapiwithair   Run API playground with Air live reload")
	fmt.Println("  mage runsocialwithair Run GoSocial with Air live reload")
	fmt.Println("  mage rundocswithair  Run GoHyperDocs with Air live reload")
}

// Clean removes build artifacts and caches
func Clean() error {
	printSection("Cleaning build artifacts...")

	if err := sh.Rm(buildDir); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := sh.Run("go", "clean", "-cache"); err != nil {
		return err
	}

	printSuccess("Clean complete")
	return nil
}

// Status shows project status and statistics
func Status() error {
	colorPrint(colorCyan, "Go Standard Library Web Stack - The Ultimate Hypermedia-Driven Web Stack")
	fmt.Println("==========================================================================")
	fmt.Println()

	// Project structure
	colorPrint(colorBlue, "📁 Project Structure:")
	webFiles, _ := countGoFiles(webDir)
	pkgFiles, _ := countGoFiles(pkgDir)

	fmt.Printf("  Web Apps:     %d Go files\n", webFiles)
	fmt.Printf("  Packages:     %d Go files\n", pkgFiles)
	fmt.Println()

	// Tech Stack
	colorPrint(colorBlue, "🔧 Tech Stack:")
	fmt.Println("  Backend:      Go + net/http")
	fmt.Println("  Frontend:     HTMX + html/template")
	fmt.Println("  Database:     SQLite (CGO-free)")
	fmt.Println("  Styling:      Vanilla CSS")
	fmt.Println("  Build:        Mage")
	fmt.Println("  Live Reload:  Air")
	fmt.Println()

	// Applications
	colorPrint(colorBlue, "📱 Applications:")
	fmt.Println("  API Playground:  Interactive API testing (Port 8080)")
	fmt.Println("  GoSocial:        Social media platform (Port 8081)")
	fmt.Println()

	// Build status
	colorPrint(colorBlue, "🏗️  Build Status:")
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		fmt.Println("  No builds found. Run 'mage build:all' to build all applications.")
	} else {
		files, _ := filepath.Glob(filepath.Join(buildDir, "*"))
		fmt.Printf("  Built binaries: %d\n", len(files))
		for _, file := range files {
			if info, err := os.Stat(file); err == nil && !info.IsDir() {
				fmt.Printf("    %s\n", filepath.Base(file))
			}
		}
	}
	fmt.Println()

	// Line count
	colorPrint(colorBlue, "📊 Line Count:")
	totalLines, _ := countTotalLines(".")
	fmt.Printf("  Total Go code: %d lines\n", totalLines)

	return nil
}

// countGoFiles counts Go files in a directory
func countGoFiles(dir string) (int, error) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			count++
		}
		return nil
	})
	return count, err
}

// countTotalLines counts total lines in Go files
func countTotalLines(dir string) (int, error) {
	totalLines := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			if lines, err := countLinesInFile(path); err == nil {
				totalLines += lines
			}
		}
		return nil
	})
	return totalLines, err
}

// countLinesInFile counts lines in a single file
func countLinesInFile(filename string) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	return len(strings.Split(string(content), "\n")), nil
}

// Init initializes the development environment
func (Dev) Init() error {
	printSection("Initializing Go Chi monorepo...")

	mg.Deps(Tools.Install)
	mg.Deps(Dev.Deps)
	mg.Deps(Dev.Tidy)

	printSuccess("Initialization complete! The Ultimate Hypermedia-Driven Web Stack is ready.")
	return nil
}

// Deps downloads all dependencies
func (Dev) Deps() error {
	printSection("Downloading dependencies...")

	if err := sh.Run("go", "work", "sync"); err != nil {
		return err
	}

	// Get all modules with go.mod files
	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  📁 Downloading dependencies for %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "mod", "download"); err != nil {
				return err
			}
		}
	}

	printSuccess("Dependencies downloaded")
	return nil
}

// Tidy tidies all Go modules
func (Dev) Tidy() error {
	printSection("Tidying Go modules...")

	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  📁 Tidying %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "mod", "tidy"); err != nil {
				return err
			}
		}
	}

	printSuccess("All modules tidied")
	return nil
}

// Fmt formats all Go code
func (Dev) Fmt() error {
	printSection("Formatting Go code...")

	if err := sh.Run("gofmt", "-w", "-s", "."); err != nil {
		return err
	}

	printSuccess("Code formatted")
	return nil
}

// Lint runs linter on all Go code
func (Dev) Lint() error {
	printSection("Running linter...")

	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  🔍 Vetting %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "vet", "./..."); err != nil {
				printWarning(fmt.Sprintf("Vetting failed for %s: %v", dir, err))
				continue
			}
		}
	}

	printSuccess("Linting complete")
	return nil
}

// Start starts the development environment
func (Dev) Start() error {
	printSection("Starting The Ultimate Hypermedia-Driven Web Stack...")

	mg.Deps(Build.Web)

	colorPrint(colorCyan, "🌐 Applications running:")
	colorPrint(colorGreen, "  API Playground: http://localhost:8080")
	colorPrint(colorGreen, "  GoSocial: http://localhost:8081")
	colorPrint(colorGreen, "  GoHyperDocs: http://localhost:8082")
	colorPrint(colorYellow, "Press Ctrl+C to stop all services")

	return runConcurrently([]string{apiPlaygroundBin, goSocialBin, goHyperDocsBin})
}

// StartWithAir starts development with Air live reloading
func (Dev) StartWithAir() error {
	printSection("Starting with Air live reloading...")

	// Check if Air is installed
	if _, err := exec.LookPath("air"); err != nil {
		printError("Air not installed. Run 'mage tools:install' first")
		return err
	}

	colorPrint(colorCyan, "🌐 Applications with live reloading:")
	colorPrint(colorGreen, "  API Playground: http://localhost:8080")
	colorPrint(colorGreen, "  GoSocial: http://localhost:8081")
	colorPrint(colorGreen, "  GoHyperDocs: http://localhost:8082")
	colorPrint(colorYellow, "Press Ctrl+C to stop all services")

	// Start all apps with Air
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// Start API Playground with Air
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sh.RunWith(map[string]string{"PWD": apiPlaygroundDir}, "air"); err != nil {
			errChan <- fmt.Errorf("API Playground failed: %w", err)
		}
	}()

	// Start GoSocial with Air
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sh.RunWith(map[string]string{"PWD": goSocialDir}, "air"); err != nil {
			errChan <- fmt.Errorf("GoSocial failed: %w", err)
		}
	}()

	// Start GoHyperDocs with Air
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := sh.RunWith(map[string]string{"PWD": goHyperDocsDir}, "air"); err != nil {
			errChan <- fmt.Errorf("GoHyperDocs failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// runConcurrently runs multiple commands concurrently
func runConcurrently(commands []string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(commands))

	for _, cmd := range commands {
		wg.Add(1)
		go func(command string) {
			defer wg.Done()
			if err := sh.Run(command); err != nil {
				errChan <- err
			}
		}(cmd)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// All builds all applications
func (Build) All() error {
	printSection("Building all applications...")

	mg.Deps(Clean)
	mg.Deps(Build.Web)

	printSuccess("All applications built successfully!")
	return Build{}.showBuildResults()
}

// Web builds all web applications
func (Build) Web() error {
	mg.Deps(Build.API, Build.Social, Build.Docs)
	return nil
}

// API builds the API playground web app
func (Build) API() error {
	printSection("Building API Playground...")

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}

	if err := sh.Run("go", "build", "-o", apiPlaygroundBin, apiPlaygroundDir); err != nil {
		return err
	}

	printSuccess(fmt.Sprintf("API Playground built: %s", apiPlaygroundBin))
	return nil
}

// Social builds the GoSocial web app
func (Build) Social() error {
	printSection("Building GoSocial...")

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}

	if err := sh.Run("go", "build", "-o", goSocialBin, goSocialDir); err != nil {
		return err
	}

	printSuccess(fmt.Sprintf("GoSocial built: %s", goSocialBin))
	return nil
}

// Docs builds the GoHyperDocs web app
func (Build) Docs() error {
	printSection("Building GoHyperDocs...")

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}

	if err := sh.Run("go", "build", "-o", goHyperDocsBin, goHyperDocsDir); err != nil {
		return err
	}

	printSuccess(fmt.Sprintf("GoHyperDocs built: %s", goHyperDocsBin))
	return nil
}

// showBuildResults displays build results
func (Build) showBuildResults() error {
	files, err := filepath.Glob(filepath.Join(buildDir, "*"))
	if err != nil {
		return err
	}

	fmt.Println()
	colorPrint(colorGreen, "🎉 Build Results:")
	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			fmt.Printf("  %s\n", file)
		}
	}

	return nil
}

// All runs all tests
func (Test) All() error {
	printSection("Running tests...")

	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  🧪 Running tests in %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "test", "./..."); err != nil {
				printWarning(fmt.Sprintf("Tests failed for %s: %v", dir, err))
				continue
			}
		}
	}

	printSuccess("All tests completed")
	return nil
}

// Quick runs quick tests without cache
func (Test) Quick() error {
	printSection("Running quick tests...")
	return sh.Run("go", "test", "-count=1", "./...")
}

// Security runs security checks
func (Test) Security() error {
	printSection("Running security checks...")

	// Check if gosec is available
	if _, err := exec.LookPath("gosec"); err != nil {
		printWarning("gosec not installed. Run 'mage tools:install' first")
		return nil
	}

	if err := sh.Run("gosec", "./..."); err != nil {
		return err
	}

	printSuccess("Security check complete")
	return nil
}

// Install installs development tools
func (Tools) Install() error {
	printSection("Installing development tools...")

	// Install Air
	if _, err := exec.LookPath("air"); err != nil {
		fmt.Println("  📦 Installing Air live reload...")
		if err := sh.Run("go", "install", "github.com/air-verse/air@latest"); err != nil {
			return err
		}
	}

	// Install golangci-lint
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		fmt.Println("  📦 Installing golangci-lint...")
		cmd := "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2"
		if err := sh.RunV("sh", "-c", cmd); err != nil {
			return err
		}
	}

	// Install gosec
	if _, err := exec.LookPath("gosec"); err != nil {
		fmt.Println("  📦 Installing gosec...")
		if err := sh.Run("go", "install", "github.com/securego/gosec/v2/cmd/gosec@latest"); err != nil {
			return err
		}
	}

	printSuccess("Development tools installed")
	return nil
}

// Upgrade upgrades all dependencies
func (Tools) Upgrade() error {
	printSection("Upgrading dependencies...")

	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  📁 Upgrading %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "get", "-u", "./..."); err != nil {
				return err
			}
		}
	}

	mg.Deps(Dev.Tidy)
	printSuccess("Dependencies upgraded")
	return nil
}

// Caddy starts Caddy reverse proxy
func (Prod) Caddy() error {
	printSection("Starting Caddy reverse proxy...")

	colorPrint(colorCyan, "Make sure your backend services are running first")
	colorPrint(colorCyan, "Access via:")
	colorPrint(colorGreen, "  API Playground: https://api-playground.localhost")
	colorPrint(colorGreen, "  GoSocial: https://go-social.localhost")

	// Check if caddy is available
	if _, err := exec.LookPath("caddy"); err != nil {
		printError("Caddy not installed. Install with: brew install caddy")
		return err
	}

	return sh.Run("caddy", "run", "--config", "Caddyfile")
}

// Release creates a production release
func (Prod) Release() error {
	printSection("Creating production release...")

	mg.Deps(Test.All, Dev.Lint, Test.Security, Build.All)

	releaseDir := filepath.Join(buildDir, "release")
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return err
	}

	// Copy binaries to release directory
	files, err := filepath.Glob(filepath.Join(buildDir, "*"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if info, err := os.Stat(file); err == nil && !info.IsDir() {
			if err := sh.Copy(filepath.Join(releaseDir, filepath.Base(file)), file); err != nil {
				return err
			}
		}
	}

	// Create tarball
	tarFile := filepath.Join(buildDir, "go-stdlib-release.tar.gz")
	if err := sh.RunV("tar", "-czf", tarFile, "-C", releaseDir, "."); err != nil {
		return err
	}

	printSuccess(fmt.Sprintf("Production release created: %s", tarFile))
	return nil
}

// Format formats all Go code in the monorepo
func Format() error {
	printSection("Formatting all Go code in monorepo...")

	// Format all Go files recursively with gofmt
	if err := sh.Run("gofmt", "-w", "-s", "."); err != nil {
		return fmt.Errorf("failed to format Go code: %w", err)
	}

	// Also run go mod tidy on all modules to clean up dependencies
	modules := []string{
		apiPlaygroundDir,
		goSocialDir,
		goHyperDocsDir,
		componentsDir,
		databaseDir,
		middlewareDir,
		stylesDir,
		utilsDir,
	}

	for _, dir := range modules {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			fmt.Printf("  📁 Tidying %s\n", dir)
			if err := sh.RunWith(map[string]string{"PWD": dir}, "go", "mod", "tidy"); err != nil {
				printWarning(fmt.Sprintf("Failed to tidy %s: %v", dir, err))
				continue
			}
		}
	}

	printSuccess("All Go code formatted and modules tidied")
	return nil
}

// Watch watches for changes and rebuilds (requires entr)
func Watch() error {
	printSection("Watching for changes...")
	printWarning("This requires 'entr' to be installed")

	cmd := exec.Command("find", ".", "-name", "*.go")
	cmd2 := exec.Command("entr", "-r", "mage", "build:all")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd2.Stdin = pipe
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd2.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return cmd2.Wait()
}

// RunAPI runs API playground only
func RunAPI() error {
	mg.Deps(Build.API)
	colorPrint(colorCyan, "🌐 API Playground: http://localhost:8080")
	return sh.Run(apiPlaygroundBin)
}

// RunSocial runs GoSocial only
func RunSocial() error {
	mg.Deps(Build.Social)
	colorPrint(colorCyan, "🌐 GoSocial: http://localhost:8081")
	return sh.Run(goSocialBin)
}

// RunAPIWithAir runs API playground with Air live reloading
func RunAPIWithAir() error {
	printSection("Starting API Playground with Air...")

	// Check if Air is installed
	if _, err := exec.LookPath("air"); err != nil {
		printError("Air not installed. Run 'mage tools:install' first")
		return err
	}

	colorPrint(colorCyan, "🌐 API Playground with live reload: http://localhost:8080")
	return sh.RunWith(map[string]string{"PWD": apiPlaygroundDir}, "air")
}

// RunSocialWithAir runs GoSocial with Air live reloading
func RunSocialWithAir() error {
	printSection("Starting GoSocial with Air...")

	// Check if Air is installed
	if _, err := exec.LookPath("air"); err != nil {
		printError("Air not installed. Run 'mage tools:install' first")
		return err
	}

	colorPrint(colorCyan, "🌐 GoSocial with live reload: http://localhost:8081")
	return sh.RunWith(map[string]string{"PWD": goSocialDir}, "air")
}

// RunDocs runs GoHyperDocs only
func RunDocs() error {
	mg.Deps(Build.Docs)
	colorPrint(colorCyan, "🌐 GoHyperDocs: http://localhost:8082")
	return sh.Run(goHyperDocsBin)
}

// RunDocsWithAir runs GoHyperDocs with Air live reloading
func RunDocsWithAir() error {
	printSection("Starting GoHyperDocs with Air...")

	// Check if Air is installed
	if _, err := exec.LookPath("air"); err != nil {
		printError("Air not installed. Run 'mage tools:install' first")
		return err
	}

	colorPrint(colorCyan, "🌐 GoHyperDocs with live reload: http://localhost:8082")
	return sh.RunWith(map[string]string{"PWD": goHyperDocsDir}, "air")
}

// Docs generates documentation
func Docs() error {
	printSection("Generating documentation...")

	docsDir := "./docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	readme := `# Go Chi Monorepo - The Ultimate Hypermedia-Driven Web Stack

## Applications

### Web Applications
- **API Playground** (Port 8080): Interactive API testing with Go + HTMX
- **GoSocial** (Port 8081): Social media platform with Go + HTMX + SQLite

## Tech Stack

- **Backend**: Go + Chi Router
- **Frontend**: HTMX + Gomponents
- **Database**: SQLite with embedded assets
- **Styling**: Vanilla CSS
- **Build**: Mage

## Usage

` + "```bash" + `
mage help           # Show all available commands
mage dev:init       # Initialize development environment
mage dev:start      # Start development mode
mage build:all      # Build all applications
mage test:all       # Run all tests
` + "```" + `
`

	if err := os.WriteFile(filepath.Join(docsDir, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}

	printSuccess("Documentation generated")
	return nil
}
