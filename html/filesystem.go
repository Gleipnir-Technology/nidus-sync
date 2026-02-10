package html

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

// The filesystem being used
var templates templateSystem

type templateSystem interface {
	loadAll() error
	renderOrError(http.ResponseWriter, string, interface{})
}

type templateSystemDisk struct {
	sourceFS fs.FS
}

func LoadTemplates() error {
	_, err := os.Stat("html/template")
	if err == nil && !config.IsProductionEnvironment() {
		templates = templateSystemDisk{
			sourceFS: os.DirFS("./html/template"),
		}
	} else {
		templates, err = newTemplateSystemEmbed()
		if err != nil {
			return fmt.Errorf("Failed to load embedded templates: %w", err)
		}
	}
	return nil
}

func (ts templateSystemDisk) loadAll() error {
	return nil
}
func (ts templateSystemDisk) renderOrError(w http.ResponseWriter, template_name string, context interface{}) {
	t, err := parseTemplate(ts.sourceFS, template_name)
	if err != nil {
		log.Error().Err(err).Str("template_name", template_name).Msg("Failed to parse template")
		RespondError(w, "Failed to parse template", err, http.StatusInternalServerError)
		return
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, context)
	if err != nil {
		log.Error().Err(err).Msg("Failed to render template")
		RespondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

/*
	func addSVGTemplates(fsys fs.FS, templ *template.Template) error {
		svgs, err := fs.ReadDir(fsys, ".")
		if err != nil {
			log.Warn().Msg("Failed to read svg directory")
			return nil
		}
		for _, svg := range svgs {
			content, err := fs.ReadFile(fsys, svg.Name())
			if err != nil {
				return fmt.Errorf("Failed to read svg '%s' from embedded filesystem: %w", svg, err)
			}
			svg_name := svg.Name()
			svg_template := fmt.Sprintf("{{define \"%s\"}}%s{{end}}", svg_name, string(content))
			svg_t, err := template.New(svg_name).Parse(svg_template)
			if err != nil {
				return fmt.Errorf("Failed to parse svg '%s' from embedded filesystem: %v", svg, err)
			}
			_, err = templ.AddParseTree(svg_t.Name(), svg_t.Tree)
			if err != nil {
				return fmt.Errorf("Failed to add svg '%s' to embedded template: %v", svg, err)
			}
			//log.Debug().Str("name", svg_name).Msg("add svg template")
		}
		return nil
	}

	func executeTemplate(w io.Writer, data any) error {
		if bt.template == nil {
			name := path.Base(bt.files[0])
			templ, err := parseFromDisk(bt.subdir, bt.files)
			if err != nil {
				return fmt.Errorf("Failed to parse template file: %w", err)
			}
			if templ == nil {
				w.Write([]byte("Failed to read from disk: "))
				return errors.New("Template parsing failed")
			}
			//log.Debug().Str("name", templ.Name()).Msg("Parsed template")
			return templ.ExecuteTemplate(w, name, data)
		} else {
			name := path.Base(bt.files[0])
			return bt.template.ExecuteTemplate(w, name, data)
		}
	}

	func parseEmbedded(embeddedFiles embed.FS, subdir string, files []string) *template.Template {
		funcMap := makeFuncMap()
		// Remap the file names to embedded paths
		embeddedFilePaths := make([]string, 0)
		for _, f := range files {
			embeddedFilePaths = append(embeddedFilePaths, strings.TrimPrefix(f, subdir))
		}
		name := path.Base(embeddedFilePaths[0])
		log.Debug().Str("name", name).Strs("paths", embeddedFilePaths).Msg("Parsing embedded template")
		t, err := template.New(name).Funcs(funcMap).ParseFS(embeddedFiles, embeddedFilePaths...)
		if err != nil {
			panic(fmt.Sprintf("Failed to parse embedded template %s: %v", name, err))
		}
		svg_fs, err := fs.Sub(embeddedFiles, "template/svg")
		if err != nil {
			panic(fmt.Sprintf("Failed to read static/svg: %v", err))
		}
		err = addSVGTemplates(svg_fs, t)
		if err != nil {
			panic(fmt.Sprintf("Failed to add SVG templates: %v", err))
		}
		return t
	}

	func parseFromDisk(subdir string, files []string) (*template.Template, error) {
		funcMap := makeFuncMap()
		name := path.Base(files[0])
		//log.Debug().Str("name", name).Strs("files", files).Msg("parsing from disk")
		templ, err := template.New(name).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s: %w", files, err)
		}
		fsys := os.DirFS(subdir + "/template/svg")
		err = addSVGTemplates(fsys, templ)
		if err != nil {
			return nil, fmt.Errorf("Failed to add SVGs from disk: %w", err)
		}
		return templ, nil
	}
*/
