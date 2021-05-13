generate-templates:
	templ generate

run: generate-templates
ifdef port
	go run . -port $(port)
else
	go run .
endif
