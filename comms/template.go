package comms

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	templatehtml "html/template"
	"io"
	"os"
	"path"
	"strings"
	templatetxt "text/template"

	"github.com/Gleipnir-Technology/nidus-sync/config"
	"github.com/rs/zerolog/log"
)

type builtTemplate struct {
	name         string
	templateHTML *templatehtml.Template
	templateTXT  *templatetxt.Template
}

func (bt *builtTemplate) executeTemplateHTML(w io.Writer, content any) error {
	if bt.templateHTML == nil {
		file := templateFileHTML(bt.name)
		templ, err := parseFromDiskHTML(file)
		if err != nil {
			return fmt.Errorf("Failed to parse template file: %w", err)
		}
		if templ == nil {
			w.Write([]byte("Failed to read from disk: "))
			return errors.New("Template parsing failed")
		}
		//log.Debug().Str("name", templ.Name()).Msg("Parsed template")
		return templ.ExecuteTemplate(w, bt.name+".html", content)
	} else {
		return bt.templateHTML.ExecuteTemplate(w, bt.name+".html", content)
	}
}
func (bt *builtTemplate) executeTemplateTXT(w io.Writer, content any) error {
	if bt.templateTXT == nil {
		file := templateFileTXT(bt.name)
		templ, err := parseFromDiskTXT(file)
		if err != nil {
			return fmt.Errorf("Failed to parse template file: %w", err)
		}
		if templ == nil {
			w.Write([]byte("Failed to read from disk: "))
			return errors.New("Template parsing failed")
		}
		//log.Debug().Str("name", templ.Name()).Msg("Parsed template")
		return templ.ExecuteTemplate(w, bt.name+".txt", content)
	} else {
		return bt.templateTXT.ExecuteTemplate(w, bt.name+".txt", content)
	}
}
func templateFileHTML(name string) string {
	return fmt.Sprintf("comms/template/%s.html", name)
}
func templateFileTXT(name string) string {
	return fmt.Sprintf("comms/template/%s.txt", name)
}

func buildTemplate(name string) *builtTemplate {
	files_on_disk := true
	file_html := templateFileHTML(name)
	file_txt := templateFileTXT(name)
	for _, f := range []string{file_html, file_txt} {
		_, err := os.Stat(f)
		if err != nil {
			files_on_disk = false
			if !config.IsProductionEnvironment() {
				log.Warn().Str("file", f).Msg("email template file is not on disk")
			}
			break
		}
	}
	var result builtTemplate
	if files_on_disk {
		result = builtTemplate{
			name:         name,
			templateHTML: nil,
			templateTXT:  nil,
		}
	} else {
		result = builtTemplate{
			name:         name,
			templateHTML: parseEmbeddedHTML(embeddedFiles, "comms", file_html),
			templateTXT:  parseEmbeddedTXT(embeddedFiles, "comms", file_txt),
		}
	}
	return &result
}

func parseEmbeddedHTML(embeddedFiles embed.FS, subdir string, file string) *templatehtml.Template {
	// Remap the file names to embedded paths
	embeddedFilePaths := []string{strings.TrimPrefix(file, subdir)}
	name := path.Base(embeddedFilePaths[0])
	log.Debug().Str("name", name).Strs("paths", embeddedFilePaths).Msg("Parsing embedded template")
	return templatehtml.Must(
		templatehtml.New(name).ParseFS(embeddedFiles, embeddedFilePaths...))
}
func parseEmbeddedTXT(embeddedFiles embed.FS, subdir string, file string) *templatetxt.Template {
	// Remap the file names to embedded paths
	embeddedFilePaths := []string{strings.TrimPrefix(file, subdir)}
	name := path.Base(embeddedFilePaths[0])
	log.Debug().Str("name", name).Strs("paths", embeddedFilePaths).Msg("Parsing embedded template")
	return templatetxt.Must(
		templatetxt.New(name).ParseFS(embeddedFiles, embeddedFilePaths...))
}

func parseFromDiskHTML(file string) (*templatehtml.Template, error) {
	name := path.Base(file)
	//log.Debug().Str("name", name).Strs("files", files).Msg("parsing from disk")
	templ, err := templatehtml.New(name).ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %w", file, err)
	}
	return templ, nil
}

func parseFromDiskTXT(file string) (*templatetxt.Template, error) {
	name := path.Base(file)
	//log.Debug().Str("name", name).Strs("files", files).Msg("parsing from disk")
	templ, err := templatetxt.New(name).ParseFiles(file)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s: %w", file, err)
	}
	return templ, nil
}

func renderEmailTemplates(t *builtTemplate, content interface{}) (text string, html string, err error) {
	buf_txt := &bytes.Buffer{}
	err = t.executeTemplateTXT(buf_txt, content)
	if err != nil {
		return "", "", fmt.Errorf("Failed to render TXT template: %w", err)
	}
	buf_html := &bytes.Buffer{}
	err = t.executeTemplateHTML(buf_html, content)
	if err != nil {
		return "", "", fmt.Errorf("Failed to render HTML template: %w", err)
	}
	return buf_txt.String(), buf_html.String(), nil
}
