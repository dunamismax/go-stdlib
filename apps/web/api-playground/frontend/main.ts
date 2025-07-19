// Import Pico.css for styling
import '@picocss/pico/css/pico.min.css'
import './styles.css'

// Import HTMX
import 'htmx.org'

// TypeScript types for better development experience
declare global {
  interface Window {
    htmx: any;
  }
}

// Custom application logic for the API Playground
class APIPlayground {
  constructor() {
    this.init()
  }

  private init(): void {
    console.log('API Playground initialized with TypeScript + Vite')
    
    // Add any custom interactivity here that goes beyond HTMX
    this.setupCustomFeatures()
  }

  private setupCustomFeatures(): void {
    // Example: Custom loading states for better UX
    document.addEventListener('htmx:beforeRequest', (event: Event) => {
      const target = (event as any).detail.elt
      if (target) {
        target.setAttribute('aria-busy', 'true')
      }
    })

    document.addEventListener('htmx:afterRequest', (event: Event) => {
      const target = (event as any).detail.elt
      if (target) {
        target.setAttribute('aria-busy', 'false')
      }
    })

    // Example: Copy to clipboard functionality
    this.setupCopyButtons()
  }

  private setupCopyButtons(): void {
    document.addEventListener('click', (event: Event) => {
      const target = event.target as HTMLElement
      if (target.classList.contains('copy-btn')) {
        const resultElement = target.previousElementSibling as HTMLElement
        if (resultElement && resultElement.textContent) {
          navigator.clipboard.writeText(resultElement.textContent.trim())
            .then(() => {
              target.textContent = 'Copied!'
              setTimeout(() => {
                target.textContent = 'Copy'
              }, 2000)
            })
            .catch(() => {
              console.error('Failed to copy text')
            })
        }
      }
    })
  }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new APIPlayground()
})