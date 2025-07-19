import '@picocss/pico/css/pico.min.css'
import './styles.css'

// Social media specific functionality
class SocialApp {
  constructor() {
    this.initializeApp()
  }

  private initializeApp(): void {
    this.setupCharacterCount()
    this.setupLikeButtons()
    this.setupPostForm()
  }

  private setupCharacterCount(): void {
    const textarea = document.querySelector('#post-content') as HTMLTextAreaElement
    const charCount = document.querySelector('.char-count') as HTMLElement
    
    if (textarea && charCount) {
      const maxLength = 280
      
      textarea.addEventListener('input', () => {
        const remaining = maxLength - textarea.value.length
        charCount.textContent = `${remaining} characters remaining`
        
        if (remaining < 0) {
          charCount.style.color = 'var(--pico-color-red-500)'
        } else if (remaining < 20) {
          charCount.style.color = 'var(--pico-color-orange-500)'
        } else {
          charCount.style.color = 'var(--pico-color-grey-500)'
        }
      })
    }
  }

  private setupLikeButtons(): void {
    document.addEventListener('click', (e) => {
      const target = e.target as HTMLElement
      if (target.classList.contains('like-btn')) {
        this.handleLike(target)
      }
    })
  }

  private handleLike(button: HTMLElement): void {
    const postId = button.dataset.postId
    if (!postId) return

    fetch(`/like/${postId}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    .then(response => response.json())
    .then(data => {
      if (data.success) {
        const likeCount = button.nextElementSibling as HTMLElement
        if (likeCount) {
          likeCount.textContent = `${data.likes} likes`
        }
        
        if (data.liked) {
          button.classList.add('liked')
          button.innerHTML = '♥'
        } else {
          button.classList.remove('liked')
          button.innerHTML = '♡'
        }
      }
    })
    .catch(error => {
      console.error('Error liking post:', error)
    })
  }

  private setupPostForm(): void {
    const form = document.querySelector('#post-form') as HTMLFormElement
    if (!form) return

    form.addEventListener('submit', (e) => {
      e.preventDefault()
      
      const formData = new FormData(form)
      const content = formData.get('content') as string
      
      if (!content.trim()) {
        this.showMessage('Post content cannot be empty', 'error')
        return
      }
      
      if (content.length > 280) {
        this.showMessage('Post content too long (max 280 characters)', 'error')
        return
      }

      fetch('/post', {
        method: 'POST',
        body: formData
      })
      .then(response => {
        if (response.ok) {
          // Refresh posts or add new post to feed
          location.reload()
        } else {
          throw new Error('Failed to create post')
        }
      })
      .catch(error => {
        console.error('Error creating post:', error)
        this.showMessage('Failed to create post', 'error')
      })
    })
  }

  private showMessage(message: string, type: 'success' | 'error'): void {
    const messageEl = document.createElement('div')
    messageEl.className = `message message-${type}`
    messageEl.textContent = message
    messageEl.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      padding: 1rem;
      border-radius: 4px;
      color: white;
      z-index: 1000;
      background: ${type === 'success' ? 'var(--pico-color-green-500)' : 'var(--pico-color-red-500)'};
    `
    
    document.body.appendChild(messageEl)
    
    setTimeout(() => {
      messageEl.remove()
    }, 5000)
  }
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
  new SocialApp()
})