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
	nameToTemplate map[string]*template.Template
	sourceFS       fs.FS
}

func newTemplateSystemEmbed() (templateSystemEmbed, error) {
	ts := templateSystemEmbed{
		sourceFS:       embeddedFiles,
		nameToTemplate: make(map[string]*template.Template),
	}

	// Load all templates
	if err := ts.loadAll(); err != nil {
		return ts, err
	}

	return ts, nil
}
func (ts templateSystemEmbed) loadAll() error {
	// Then, parse all remaining templates into their named slots, adding the shared stuff
	page_subdirs := []string{"rmo", "sync"}
	for _, subdir := range page_subdirs {
		err := ts.loadTemplateSubdir(subdir)
		if err != nil {
			return fmt.Errorf("Failed to load subdir '%s': %w", subdir, err)
		}
	}
	return nil
}

func (ts templateSystemEmbed) loadTemplateSubdir(subdir string) error {
	template_fs, err := fs.Sub(ts.sourceFS, "template")
	if err != nil {
		return fmt.Errorf("Failed to create template sub-fs: %w", err)
	}
	err = fs.WalkDir(template_fs, subdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		new_t, err := parseTemplate(template_fs, path)
		if err != nil {
			return fmt.Errorf("error parsing '%s': %w", path, err)
		}
		ts.nameToTemplate[path] = new_t
		log.Debug().Str("path", path).Msg("Loaded page template")
		return nil
	})
	return err
}

func (ts templateSystemEmbed) addSubdirTemplates(t *template.Template, subdir string) error {
	var err error
	//log.Debug().Msg("Adding subdir templates")
	err = fs.WalkDir(ts.sourceFS, subdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := fs.ReadFile(ts.sourceFS, path)
		if err != nil {
			return fmt.Errorf("error reading template %s: %w", path, err)
		}

		new_t := template.New(path)
		addFuncMap(new_t)
		_, err = new_t.Parse(string(content))
		if err != nil {
			return fmt.Errorf("error parsing '%s': %w", path, err)
		}
		short_path := removeLeadingDir(path)
		_, err = t.AddParseTree(short_path, new_t.Tree)
		if err != nil {
			return fmt.Errorf("error adding parse tree '%s': %w", path, err)
		}
		log.Debug().Str("path", short_path).Msg("Loaded shared component template")
		return nil
	})
	return err
}

func (ts templateSystemEmbed) renderOrError(w http.ResponseWriter, template_name string, context interface{}) {
	buf := &bytes.Buffer{}

	// Execute the template directly from the pre-parsed set
	templ, ok := ts.nameToTemplate[template_name]
	if !ok {
		log.Error().Str("template_name", template_name).Msg("Can't find template")
		RespondError(w, "Failed to find template", nil, http.StatusInternalServerError)
		return
	}
	err := templ.Execute(buf, context)
	if err != nil {
		log.Error().Err(err).Str("template_name", template_name).Msg("Failed to render template")
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
