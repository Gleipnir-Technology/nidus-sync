package html

import (
	"fmt"
	"html/template"
	"io/fs"
	//"github.com/rs/zerolog/log"
)

func addSupportingTemplates(sourceFS fs.FS, t *template.Template) error {
	err := fs.WalkDir(sourceFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(sourceFS, path)
		if err != nil {
			return fmt.Errorf("error reading template %s: %w", path, err)
		}

		new_t := template.New(path)
		addFuncMap(new_t)
		_, err = new_t.Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parsing '%s': %w", path, err)
		}
		_, err = t.AddParseTree(new_t.Name(), new_t.Tree)
		if err != nil {
			return fmt.Errorf("error adding parse tree '%s': %w", path, err)
		}
		//log.Debug().Str("path", path).Msg("Read template")
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking template directory: %w", err)
	}
	return nil
}
func addSVGTemplates(filesystem fs.FS, t *template.Template) error {
	err := fs.WalkDir(filesystem, "svg", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := fs.ReadFile(filesystem, path)
		if err != nil {
			return fmt.Errorf("Failed to read svg '%s' from filesystem: %w", path, err)
		}
		svg_name := removeLeadingDir(path)
		svg_template := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", svg_name, string(content))
		svg_t, err := template.New(svg_name).Parse(svg_template)
		if err != nil {
			return fmt.Errorf("Failed to parse svg '%s' from embedded filesystem: %v", path, err)
		}
		_, err = t.AddParseTree(svg_t.Name(), svg_t.Tree)
		if err != nil {
			return fmt.Errorf("Failed to add svg '%s' to embedded template: %v", path, err)
		}
		//log.Debug().Str("name", svg_name).Msg("add svg template")
		return nil
	})
	//log.Debug().Msg("Done adding SVG templates")
	return err
}
func parseTemplate(sourceFS fs.FS, template_name string) (*template.Template, error) {
	t, err := parseTemplateFile(sourceFS, template_name)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse template: %w", err)
	}
	err = addSupportingTemplates(sourceFS, t)
	if err != nil {
		return nil, fmt.Errorf("Failed to add supporting templates: %w", err)
	}
	err = addSVGTemplates(sourceFS, t)
	if err != nil {
		return nil, fmt.Errorf("Failed to add supporting svg templates: %w", err)
	}
	return t, nil
}
func parseTemplateFile(sourceFS fs.FS, filename string) (*template.Template, error) {
	t := template.New(filename)
	//log.Debug().Str("filename", filename).Msg("parsing template")
	addFuncMap(t)
	content, err := fs.ReadFile(sourceFS, filename)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %w", filename, err)
	}
	_, err = t.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s': %w", filename, err)
	}
	return t, nil
}
