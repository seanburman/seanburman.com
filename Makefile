.PHONY: templ

make:
	go run .

assets: templ
	tsc -p "static/scripts"
	./es-build
	./tailwindcss -i static/stylesheets/tailwind.css -o static/stylesheets/tailwind.min.css --minify
	sass static/sass:static/stylesheets

templ:
	/Users/seanburman/go/bin/templ generate

minify:
	./es-build

tsc:
	tsc -p "static/scripts" --watch

tailwind:
	./tailwindcss -i static/stylesheets/tailwind.css -o static/stylesheets/tailwind.min.css --watch --minify

sass:
	sass --watch static/sass:static/stylesheets