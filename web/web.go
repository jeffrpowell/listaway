package web

import (
	"embed"
	"html/template"
	"io"
	"log"
	"path/filepath"

	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

//go:embed *
var htmlFiles embed.FS
var (
	registerAdmin = parseSingleLayout("registerAdmin.html")
	login         = parseSplitLayout("login.html")
	lists         = parseSplitLayout("lists.html")
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

		mb, err := m.Bytes("text/html", b) //BUG: lower-cases go interpolation tags
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.Parse(string(mb))
		if err != nil {
			return nil, err
		}
	}
	return tmpl, nil
}

func parseSplitLayout(file string) *template.Template {
	return template.Must(minifyTemplates("root.html", "splitLayout.html", file))
}

func parseSingleLayout(file string) *template.Template {
	return template.Must(minifyTemplates("root.html", "singleLayout.html", file))
}

// Register Admin page

type registerAdminParams struct {
	AdminExists bool
	ShowNavbar  bool
}

func RegisterAdminParams(adminExists bool) registerAdminParams {
	return registerAdminParams{
		AdminExists: adminExists,
		ShowNavbar:  false,
	}
}

func RegisterAdmin(w io.Writer, params registerAdminParams) {
	if err := registerAdmin.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Login page

type loginPageParams struct {
	ShowNavbar bool
}

func LoginPage(w io.Writer) {
	if err := login.Execute(w, loginPageParams{ShowNavbar: false}); err != nil {
		log.Print(err)
	}
}

// Lists page

type listsPageParams struct {
	Lists      []constants.List
	ShowNavbar bool
}

func ListsPageParams(lists []constants.List) listsPageParams {
	return listsPageParams{
		Lists:      lists,
		ShowNavbar: true,
	}
}

func ListsPage(w io.Writer, params listsPageParams) {
	if err := lists.Execute(w, params); err != nil {
		log.Print(err)
	}
}
