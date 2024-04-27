package web

import (
	"embed"
	"html/template"
	"io"
	"path/filepath"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

//go:embed *
var htmlFiles embed.FS
var (
	registerAdmin = parse("registerAdmin.html")
	login         = parse("login.html")
	lists         = parse("lists.html")
)

func minifyTemplates(filenames ...string) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		}

		b, err := htmlFiles.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		tmpl.Parse(string(mb))
	}
	return tmpl, nil
}

func parse(file string) *template.Template {
	return template.Must(minifyTemplates("root.html", file))
}

type RegisterAdminParams struct {
	AdminExists bool
}

func RegisterAdmin(w io.Writer, params RegisterAdminParams) error {
	return registerAdmin.Execute(w, params)
}

func LoginPage(w io.Writer) error {
	return login.Execute(w, nil)
}

type ListsPageParams struct {
	Lists []constants.List
}

func ListsPage(w io.Writer, params ListsPageParams) error {
	return lists.Execute(w, params)
}
