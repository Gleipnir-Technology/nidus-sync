package html

import (
	"bytes"
	"embed"
	//"errors"
	"fmt"
	"html/template"
	//"io"
	"io/fs"
	//"math"
	"net/http"
	//"path"
	//"strconv"
	"strings"
	//"time"

	//"github.com/Gleipnir-Technology/nidus-sync/config"
	//"github.com/aarondl/opt/null"
	//"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

//go:embed template/*
var embeddedFiles embed.FS

type templateSystemEmbed struct {
	allTemplates   *template.Template
	nameToTemplate map[string]*template.Template
	sourceFS       fs.FS
}

func (ts templateSystemEmbed) loadAll() error {
	ts.nameToTemplate = make(map[string]*template.Template, 0)
	err := fs.WalkDir(ts.sourceFS, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}
		short_path := removeLeadingDir(path)

		new_t, err := loadTemplateEmbedded(ts.sourceFS, short_path)
		if err != nil {
			return fmt.Errorf("Failed to add load template '%s': %w", short_path, err)
		}
		_, err = ts.allTemplates.AddParseTree(new_t.Name(), new_t.Tree)
		if err != nil {
			return fmt.Errorf("Failed to add parsed template '%s': %w", path, err)
		}
		ts.nameToTemplate[short_path] = new_t
		log.Debug().Str("path", path).Str("short_path", short_path).Msg("Loaded template")
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load embeded templates: %w", err)
	}
	return nil
}
func (ts templateSystemEmbed) renderOrError(w http.ResponseWriter, template_name string, context interface{}) {
	buf := &bytes.Buffer{}
	t, err := loadTemplateEmbedded(ts.sourceFS, template_name)
	if err != nil {
		log.Error().Err(err).Str("template_name", template_name).Msg("Failed to load embedded template")
		RespondError(w, "Failed to load template", err, http.StatusInternalServerError)
		return
	}
	for name, templ := range ts.nameToTemplate {
		_, err := t.AddParseTree(name, templ.Tree)
		if err != nil {
			log.Error().Err(err).Str("name", name).Msg("Failed to add parse tree")
			RespondError(w, "Failed to add template", err, http.StatusInternalServerError)
			return
		}
	}
	err = t.ExecuteTemplate(buf, template_name, context)
	if err != nil {
		log.Error().Err(err).Str("template_name", template_name).Msg("Failed to render embedded template")
		RespondError(w, "Failed to render template", err, http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}
func loadTemplateEmbedded(sourceFS fs.FS, path string) (*template.Template, error) {
	content, err := fs.ReadFile(sourceFS, "template/"+path)
	if err != nil {
		return nil, fmt.Errorf("error reading template template/%s: %w", path, err)
	}

	new_t := template.New(path)
	addFuncMap(new_t)
	_, err = new_t.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("error parsing '%s': %w", path, err)
	}
	return new_t, nil
}
func removeLeadingDir(path string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return path
}
