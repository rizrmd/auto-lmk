package main

import (
	"html/template"
	"log"
)

func main() {
	// Define all the functions that should be available
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"formatPrice": func(price int64) string {
			return "Rp " + "formatted_price"
		},
		"formatMileage": func(km int) string {
			return "formatted_mileage"
		},
		// Missing functions that cause errors:
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" {
				return def
			}
			return val
		},
		"js": func(v interface{}) string {
			return "js_escaped"
		},
	}

	// Try to parse templates like in PageHandler
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
		"templates/layouts/base.html",
		"templates/components/button.html",
		"templates/components/card.html",
		"templates/components/input.html",
		"templates/components/nav.html",
		"templates/components/hero.html",
		"templates/components/footer.html",
		"templates/components/gallery.html",
		"templates/components/pagination.html",
		"templates/pages/home.html",
		"templates/pages/cars.html",
		"templates/pages/car-detail.html",
		"templates/pages/contact.html",
		"templates/pages/blog.html",
		"templates/pages/blog-detail.html",
		"templates/admin/whatsapp.html",
	))

	log.Println("Templates parsed successfully!")
	_ = tmpl
}