package web

import (
	"embed"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

//go:embed dist/*
var staticFiles embed.FS
var (
	registerAdmin = parseSingleLayout("dist/registerAdmin.html")
	login         = parseSplitLayout("dist/login.html")
	lists         = parseSplitLayout("dist/lists.html")
	createList    = parseSplitLayout("dist/listCreate.html")
	editList      = parseSplitLayout("dist/listEdit.html")
	listItems     = parseSplitLayout("dist/listItems.html")
	createItem    = parseSplitLayout("dist/itemCreate.html")
)

func init() {
	constants.ROUTER.HandleFunc("/static/{pathname...}", middleware.DefaultPublicMiddlewareChain(staticHandler)).Methods("GET")
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filePath := mux.Vars(r)["pathname..."]
	// Open the file from the embedded file system
	file, err := staticFiles.Open("dist/" + filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get the file extension
	ext := filepath.Ext(filePath)
	// Set the content type based on the file extension
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// If the content type is unknown, default to "application/octet-stream"
		contentType = "application/octet-stream"
	}

	// Set the content type header
	w.Header().Set("Content-Type", contentType)

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}

func minifyTemplates(filenames ...string) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		}

		b, err := staticFiles.ReadFile(filename)
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
	return template.Must(minifyTemplates("dist/root.html", "dist/splitLayout.html", file))
}

func parseSingleLayout(file string) *template.Template {
	return template.Must(minifyTemplates("dist/root.html", "dist/singleLayout.html", file))
}

type globalWebParams struct {
	ShowNavbar bool
	JsFile     string
}

// Register Admin page

type registerAdminParams struct {
	globalWebParams
	AdminExists bool
}

func RegisterAdminParams(adminExists bool) registerAdminParams {
	return registerAdminParams{
		globalWebParams: globalWebParams{
			ShowNavbar: false,
			JsFile:     "registerAdmin",
		},
		AdminExists: adminExists,
	}
}

func RegisterAdmin(w io.Writer, params registerAdminParams) {
	if err := registerAdmin.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Login page

type loginPageParams struct {
	globalWebParams
}

func LoginPage(w io.Writer) {
	if err := login.Execute(w, loginPageParams{globalWebParams{ShowNavbar: false, JsFile: "login"}}); err != nil {
		log.Print(err)
	}
}

// Lists page

type listsPageParams struct {
	Lists []constants.List
	globalWebParams
}

func ListsPageParams(lists []constants.List) listsPageParams {
	return listsPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar: true,
			JsFile:     "lists",
		},
		Lists: lists,
	}
}

func ListsPage(w io.Writer, params listsPageParams) {
	if err := lists.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create List page

type createListParams struct {
	globalWebParams
}

func CreateListPage(w io.Writer) {
	if err := createList.Execute(w, createListParams{globalWebParams{ShowNavbar: true, JsFile: "listCreate"}}); err != nil {
		log.Print(err)
	}
}

// Edit List page

type editListParams struct {
	List constants.List
	globalWebParams
}

func EditListParams(list constants.List) editListParams {
	return editListParams{
		globalWebParams: globalWebParams{
			ShowNavbar: true,
			JsFile:     "listEdit",
		},
		List: list,
	}
}

func EditListPage(w io.Writer, params editListParams) {
	if err := editList.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// List Items page

type listItemsPageParams struct {
	List  constants.List
	Items []constants.Item
	globalWebParams
}

func ListItemsPageParams(list constants.List, items []constants.Item) listItemsPageParams {
	return listItemsPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar: true,
			JsFile:     "listItems",
		},
		List:  list,
		Items: items,
	}
}

func ListItemsPage(w io.Writer, params listItemsPageParams) {
	if err := listItems.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create/Edit Item page

type createEditItemParams struct {
	globalWebParams
	List     constants.List
	Item     constants.Item
	EditMode bool
}

func CreateItemParams(list constants.List) createEditItemParams {
	return createEditItemParams{
		globalWebParams: globalWebParams{
			ShowNavbar: true,
			JsFile:     "itemCreate",
		},
		List:     list,
		Item:     constants.Item{},
		EditMode: false,
	}
}

func EditItemParams(list constants.List, item constants.Item) createEditItemParams {
	return createEditItemParams{
		globalWebParams: globalWebParams{
			ShowNavbar: true,
			JsFile:     "itemCreate",
		},
		List:     list,
		Item:     item,
		EditMode: true,
	}
}

func CreateEditItemPage(w io.Writer, params createEditItemParams) {
	if err := createItem.Execute(w, params); err != nil {
		log.Print(err)
	}
}
