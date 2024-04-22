package web

import (
	"embed"
	"html/template"
	"io"

	"github.com/jeffrpowell/listaway/internal/constants"
)

//go:embed *
var htmlFiles embed.FS
var (
	login = parse("login.html")
	lists = parse("lists.html")
)

func parse(file string) *template.Template {
	return template.Must(template.New("root.html").ParseFS(htmlFiles, "root.html", file))
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
