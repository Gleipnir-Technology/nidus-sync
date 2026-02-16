package email

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"errors"
	"fmt"
	templatehtml "html/template"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	templatetxt "text/template"
	"time"

	"github.com/Gleipnir-Technology/bob"
	"github.com/Gleipnir-Technology/bob/dialect/psql"
	"github.com/Gleipnir-Technology/bob/dialect/psql/um"
	"github.com/Gleipnir-Technology/bob/types/pgtypes"
	"github.com/Gleipnir-Technology/nidus-sync/db"
	"github.com/Gleipnir-Technology/nidus-sync/db/enums"
	"github.com/Gleipnir-Technology/nidus-sync/db/models"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/rs/zerolog/log"
)

//go:embed template/*
var embeddedFiles embed.FS

var (
	templateByID                             map[int32]*builtTemplate
	templateInitialID                        int32
	templateReportNotificationConfirmationID int32
)

type contentEmailBase struct {
	URLLogo          string
	URLUnsubscribe   string
	URLViewInBrowser string
}

type ContentEmailRender struct {
	IsBrowser bool
	C         any
}

type templatePair struct {
	baseName    string
	messageType enums.CommsMessagetypeemail
	htmlContent string
	txtContent  string
	htmlHash    string
	txtHash     string
}

func LoadTemplates() error {
	all_templates, err := readTemplates(embeddedFiles)
	if err != nil {
		return fmt.Errorf("Failed to read templates: %w", err)
	}
	ctx := context.TODO()
	tx, err := db.PGInstance.BobDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	templateByID = make(map[int32]*builtTemplate, 0)
	for name, p := range all_templates {
		template_id, err := templateDBID(tx, name, p)
		if err != nil {
			return fmt.Errorf("Failed to add '%s' to DB: %w", name, err)
		}
		template_html, err := templatehtml.New(name).Parse(p.htmlContent)
		if err != nil {
			return fmt.Errorf("Failed to parse HTML portion of '%s': %w", name, err)
		}
		template_txt, err := templatetxt.New(name).Parse(p.txtContent)
		if err != nil {
			return fmt.Errorf("Failed to parse HTML portion of '%s': %w", name, err)
		}
		built := builtTemplate{
			name:         name,
			templateHTML: template_html,
			templateTXT:  template_txt,
		}
		templateByID[template_id] = &built
		//log.Debug().Int32("id", template_id).Str("name", name).Msg("Added template to cache")
	}
	templateInitialID, err = loadTemplateID(ctx, tx, enums.CommsMessagetypeemailInitialContact)
	if err != nil {
		return fmt.Errorf("Failed to load template ID: %s", err)
	}
	templateReportNotificationConfirmationID, err = loadTemplateID(ctx, tx, enums.CommsMessagetypeemailReportNotificationConfirmation)
	if err != nil {
		return fmt.Errorf("Failed to load report-notification-confirmation template ID: %s", err)
	}
	tx.Commit(ctx)
	return nil
}

func RenderHTML(template_id int32, s pgtypes.HStore) (html []byte, err error) {
	data := db.ConvertFromPGData(s)
	t, ok := templateByID[template_id]
	if !ok {
		return []byte{}, fmt.Errorf("Failed to lookup template %d", template_id)
	}
	buf_html := &bytes.Buffer{}
	content := ContentEmailRender{
		C:         data,
		IsBrowser: true,
	}
	err = t.executeTemplateHTML(buf_html, content)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to render HTML template: %w", err)
	}
	return buf_html.Bytes(), nil
}

func loadTemplateID(ctx context.Context, tx bob.Tx, t enums.CommsMessagetypeemail) (int32, error) {
	templates, err := models.CommsEmailTemplates.Query(
		models.SelectWhere.CommsEmailTemplates.MessageType.EQ(t),
		models.SelectWhere.CommsEmailTemplates.Superceded.IsNull(),
	).All(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("Failed to query template '%s': %w", t, err)
	}
	switch len(templates) {
	case 0:
		return 0, fmt.Errorf("No matching templates for '%s", t)
	case 1:
		return templates[0].ID, nil
	default:
		return 0, fmt.Errorf("Found %d templates for '%s', should only have 1", len(templates), t)
	}
}

func readTemplates(filesystem embed.FS) (results map[string]*templatePair, err error) {
	// First pass: read files and organize by base name
	results = make(map[string]*templatePair)

	err = fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Read file content
		content, err := filesystem.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading template %s: %w", path, err)
		}

		// Calculate hash
		hash := fmt.Sprintf("%x", sha256.Sum256(content))

		// Extract base name and extension
		ext := strings.ToLower(filepath.Ext(path))
		baseName := strings.TrimSuffix(filepath.Base(path), ext)

		// Store in map by base name
		if _, exists := results[baseName]; !exists {
			t, err := messageTypeFromName(baseName)
			if err != nil {
				return fmt.Errorf("Cannot parse email templates: %w", err)
			}
			results[baseName] = &templatePair{
				baseName:    baseName,
				messageType: *t,
			}
		}

		// Add content based on extension
		switch ext {
		case ".html", ".htm":
			results[baseName].htmlContent = string(content)
			results[baseName].htmlHash = hash
		case ".txt":
			results[baseName].txtContent = string(content)
			results[baseName].txtHash = hash
		}

		return nil
	})

	if err != nil {
		return results, fmt.Errorf("error walking template directory: %w", err)
	}

	return results, nil
}

func templateDBID(tx bob.Tx, name string, pair *templatePair) (int32, error) {
	ctx := context.Background()

	// Skip incomplete pairs
	if pair.htmlContent == "" {
		return 0, fmt.Errorf("Bad template pair '%s': no html content")
	}
	if pair.txtContent == "" {
		return 0, fmt.Errorf("Bad template pair '%s': no txt content")
	}

	// Check if a template with these hashes already exists
	rows, err := models.CommsEmailTemplates.Query(
		models.SelectWhere.CommsEmailTemplates.ContentHashHTML.EQ(pair.htmlHash),
		models.SelectWhere.CommsEmailTemplates.ContentHashTXT.EQ(pair.txtHash),
		models.SelectWhere.CommsEmailTemplates.MessageType.EQ(pair.messageType),
	).All(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("Failed to query for existing template: %w", err)
	}
	if len(rows) > 1 {
		return 0, fmt.Errorf("Got %d template rows, should only have 1", len(rows))
	} else if len(rows) == 1 {
		return rows[0].ID, nil
	}

	// Supercede previous templates of this type
	_, err = psql.Update(
		um.Table(models.CommsEmailTemplates.Alias()),
		um.SetCol("superceded").ToArg(time.Now()),
		//um.Where(models.CommsEmailTemplates.Columns.MessageType.EQ(psql.Arg(pair.messageType))),
		um.Where(psql.Quote("message_type").EQ(psql.Arg(pair.messageType))),
		//um.Where(models.CommsEmailTemplates.Columns.Superceded.IsNull()),
		um.Where(psql.Quote("superceded").IsNull()),
	).Exec(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("error superceding templates: %w", err)
	}

	new_template, err := models.CommsEmailTemplates.Insert(&models.CommsEmailTemplateSetter{
		ContentHTML:     omit.From(pair.htmlContent),
		ContentTXT:      omit.From(pair.txtContent),
		ContentHashHTML: omit.From(pair.htmlHash),
		ContentHashTXT:  omit.From(pair.txtHash),
		Created:         omit.From(time.Now()),
		Superceded:      omitnull.FromPtr[time.Time](nil),
		MessageType:     omit.From(pair.messageType),
	}).One(ctx, tx)
	if err != nil {
		return 0, fmt.Errorf("Failed to insert new template: %w", err)
	}
	log.Info().Int32("id", new_template.ID).Str("type", string(pair.messageType)).Msg("Added new email template")

	return new_template.ID, nil
}

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
		return templ.ExecuteTemplate(w, bt.name, content)
	} else {
		return bt.templateHTML.ExecuteTemplate(w, bt.name, content)
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
		return templ.ExecuteTemplate(w, bt.name, content)
	} else {
		return bt.templateTXT.ExecuteTemplate(w, bt.name, content)
	}
}
func templateFileHTML(name string) string {
	return fmt.Sprintf("comms/template/%s.html", name)
}
func templateFileTXT(name string) string {
	return fmt.Sprintf("comms/template/%s.txt", name)
}

func messageTypeFromName(n string) (*enums.CommsMessagetypeemail, error) {
	for _, t := range enums.AllCommsMessagetypeemail() {
		if n == string(t) {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("Unrecognized email type '%s'", n)
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

func publicReportID(s string) string {
	if len(s) != 12 {
		return s
	}
	return s[0:4] + "-" + s[4:8] + "-" + s[8:12]
}

func renderEmailTemplates(template_id int32, data map[string]string) (text string, html string, err error) {
	buf_txt := &bytes.Buffer{}
	t, ok := templateByID[template_id]
	if !ok {
		return "", "", fmt.Errorf("Failed to lookup template %d", template_id)
	}
	content := ContentEmailRender{
		C:         data,
		IsBrowser: false,
	}
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
