package components

// HTMX helper functions for generating HTML attributes
func HxPost(url string) string {
	return `hx-post="` + url + `"`
}

func HxTarget(target string) string {
	return `hx-target="` + target + `"`
}

func HxInclude(include string) string {
	return `hx-include="` + include + `"`
}

func HxSwap(swap string) string {
	return `hx-swap="` + swap + `"`
}

func HxOn(event string, code string) string {
	return `hx-on:` + event + `="` + code + `"`
}

func HxGet(url string) string {
	return `hx-get="` + url + `"`
}