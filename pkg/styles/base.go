package styles

import "strings"

// NewBaseStyles creates a new stylesheet with base styles
func NewBaseStyles() string {
	var styles strings.Builder

	// Global styles
	styles.WriteString(`
* {
	box-sizing: border-box;
	margin: 0;
	padding: 0;
}

body {
	font-family: system-ui, -apple-system, sans-serif;
	background-color: #0f0f0f;
	color: #e4e4e4;
	line-height: 1.6;
}

.container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 2rem;
}

.card {
	background: #1a1a1a;
	border-radius: 8px;
	padding: 1.5rem;
	margin: 1rem 0;
	border: 1px solid #333;
}

input, textarea, select {
	background: #2a2a2a;
	border: 1px solid #444;
	color: #e4e4e4;
	padding: 0.75rem;
	border-radius: 4px;
	font-size: 1rem;
	width: 100%;
	margin: 0.5rem 0;
}

button {
	background: #2563eb;
	color: white;
	border: none;
	padding: 0.75rem 1.5rem;
	border-radius: 4px;
	cursor: pointer;
	font-size: 1rem;
	margin: 0.5rem 0.5rem 0.5rem 0;
	transition: background-color 0.2s ease;
}

button:hover {
	background: #1d4ed8;
}

button:disabled {
	background: #6b7280;
	cursor: not-allowed;
}

.result {
	background: #2a2a2a;
	border-radius: 4px;
	padding: 1rem;
	margin: 1rem 0;
	border: 1px solid #444;
	word-wrap: break-word;
}

.grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
	gap: 1rem;
	margin: 1rem 0;
}

h1, h2, h3, h4, h5, h6 {
	color: #f3f4f6;
	margin-bottom: 1rem;
}

p {
	margin-bottom: 1rem;
}

a {
	color: #3b82f6;
	text-decoration: none;
}

a:hover {
	color: #1d4ed8;
	text-decoration: underline;
}
`)

	return styles.String()
}

// AddSocialStyles adds social media specific styles
func AddSocialStyles(baseStyles string) string {
	var styles strings.Builder
	styles.WriteString(baseStyles)

	styles.WriteString(`
.post {
	background: #1a1a1a;
	border-radius: 8px;
	padding: 1.5rem;
	margin: 1rem 0;
	border: 1px solid #333;
}

.post-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 1rem;
}

.post-author {
	font-weight: bold;
	color: #f3f4f6;
}

.post-time {
	color: #9ca3af;
	font-size: 0.875rem;
}

.post-content {
	margin-bottom: 1rem;
	line-height: 1.6;
}

.post-actions {
	display: flex;
	gap: 1rem;
}

.like-btn {
	background: transparent;
	color: #9ca3af;
	border: 1px solid #4b5563;
	padding: 0.5rem 1rem;
	border-radius: 4px;
	cursor: pointer;
	font-size: 0.875rem;
	transition: all 0.2s ease;
}

.like-btn:hover {
	color: #ef4444;
	border: 1px solid #ef4444;
}

.like-btn.liked {
	color: #ef4444;
	border: 1px solid #ef4444;
	background: rgba(239, 68, 68, 0.1);
}

.nav {
	background: #1a1a1a;
	padding: 1rem 0;
	margin-bottom: 2rem;
	border-bottom: 1px solid #333;
}

.nav-content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 2rem;
}

.nav-links {
	display: flex;
	gap: 1rem;
}

.form-group {
	margin-bottom: 1rem;
}

.form-label {
	display: block;
	margin-bottom: 0.5rem;
	color: #f3f4f6;
	font-weight: 500;
}

.error {
	color: #ef4444;
	font-size: 0.875rem;
	margin-top: 0.25rem;
	margin-bottom: 1rem;
}

.success {
	color: #10b981;
	font-size: 0.875rem;
	margin-top: 0.25rem;
	margin-bottom: 1rem;
}

.w-full {
	width: 100%;
}

.max-w-md {
	max-width: 28rem;
}

.max-w-2xl {
	max-width: 42rem;
}

.text-center {
	text-align: center;
}

.text-2xl {
	font-size: 1.5rem;
}

.text-3xl {
	font-size: 1.875rem;
}

.text-xl {
	font-size: 1.25rem;
}

.text-lg {
	font-size: 1.125rem;
}

.text-sm {
	font-size: 0.875rem;
}

.font-bold {
	font-weight: 700;
}

.font-semibold {
	font-weight: 600;
}

.mb-4 {
	margin-bottom: 1rem;
}

.mb-6 {
	margin-bottom: 1.5rem;
}

.mb-3 {
	margin-bottom: 0.75rem;
}

.mt-4 {
	margin-top: 1rem;
}

.space-y-4 > * + * {
	margin-top: 1rem;
}

.flex {
	display: flex;
}

.gap-4 {
	gap: 1rem;
}

.gap-2 {
	gap: 0.5rem;
}

.justify-center {
	justify-content: center;
}

.justify-between {
	justify-content: space-between;
}

.items-center {
	align-items: center;
}

.items-start {
	align-items: flex-start;
}

.text-gray-400 {
	color: #9ca3af;
}

.text-blue-400 {
	color: #60a5fa;
}

.text-blue-300 {
	color: #93c5fd;
}

.text-red-400 {
	color: #f87171;
}

.hover\\:text-blue-300:hover {
	color: #93c5fd;
}

.hover\\:text-red-400:hover {
	color: #f87171;
}

.btn {
	display: inline-block;
	background: #2563eb;
	color: white;
	border: none;
	padding: 0.75rem 1.5rem;
	border-radius: 4px;
	cursor: pointer;
	font-size: 1rem;
	margin: 0.5rem 0.5rem 0.5rem 0;
	transition: background-color 0.2s ease;
	text-decoration: none;
}

.btn:hover {
	background: #1d4ed8;
	text-decoration: none;
}

.btn-secondary {
	background: #6b7280;
}

.btn-secondary:hover {
	background: #4b5563;
}

.textarea {
	min-height: 120px;
	resize: vertical;
}
`)

	return styles.String()
}
