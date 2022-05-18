package main

import (
	"html/template"
	"path/filepath"
)

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	// Get all the template names that live in the pages/ directory
	pages, err := filepath.Glob(filepath.Join(templateRoot, "pages/*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the name from the file
		name := filepath.Base(page)

		// Parse the base template into a template set
		ts, err := template.ParseFiles(filepath.Join(templateRoot, "base.tmpl"))

		if err != nil {
			return nil, err
		}

		// call ParseGlob from the current template set ts to parse
		// partials templates
		// ts, err = ts.ParseGlob(filepath.Join(templateRoot, "partials/*.tmpl"))
		// if err != nil {
		//  return nil, err
		// }

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}