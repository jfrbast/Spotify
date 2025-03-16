package temp

import (
	"html/template"
	"log"
)

var Temp *template.Template

var funcMap = template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
}

func InitTemplates() {
	var err error
	Temp, err = template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}
