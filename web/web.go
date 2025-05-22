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
	lists         = parseSingleLayout("dist/lists.html")
	createList    = parseSingleLayout("dist/listCreate.html")
	editList      = parseSingleLayout("dist/listEdit.html")
	listItems     = parseSingleLayout("dist/listItems.html")
	createItem    = parseSingleLayout("dist/itemCreate.html")
	sharedList    = parseSingleLayout("dist/sharedList.html")
	sharedList404 = parseSingleLayout("dist/sharedList404.html")
	userAdmin     = parseSingleLayout("dist/userAdmin.html")
	allUsers      = parseSingleLayout("dist/allUsers.html")
	userCreate    = parseSingleLayout("dist/userCreate.html")
	resetRequest  = parseSplitLayout("dist/resetRequest.html")
	resetForm     = parseSplitLayout("dist/resetForm.html")
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
	ShowNavbar        bool
	ShowAdmin         bool
	ShowInstanceAdmin bool
	ChunkName         string
}

// Register Admin page

type registerAdminParams struct {
	globalWebParams
	AdminExists bool
}

func RegisterAdminParams(adminExists bool) registerAdminParams {
	return registerAdminParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        false,
			ShowAdmin:         false,
			ShowInstanceAdmin: false,
			ChunkName:         "registerAdmin",
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

func LoginPage(w io.Writer) {
	if err := login.Execute(w, globalWebParams{ShowNavbar: false, ChunkName: "login"}); err != nil {
		log.Print(err)
	}
}

// Password Reset pages

type resetRequestPageParams struct {
	globalWebParams
}

type resetFormPageParams struct {
	globalWebParams
	TokenValid bool
}

func ResetRequestPage(w io.Writer) {
	if err := resetRequest.Execute(w, resetRequestPageParams{
		globalWebParams: globalWebParams{ShowNavbar: false, ChunkName: "resetRequest"},
	}); err != nil {
		log.Print(err)
	}
}

func ResetFormPage(w io.Writer, tokenValid bool) {
	if err := resetForm.Execute(w, resetFormPageParams{
		globalWebParams: globalWebParams{ShowNavbar: false, ChunkName: "resetForm"},
		TokenValid:      tokenValid,
	}); err != nil {
		log.Print(err)
	}
}

// Lists page

type listsPageParams struct {
	Lists          []constants.List
	SharedListPath string
	globalWebParams
}

func ListsPageParams(lists []constants.List, showAdmin bool, showInstanceAdmin bool) listsPageParams {
	return listsPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "lists",
		},
		Lists:          lists,
		SharedListPath: constants.SHARED_LIST_PATH,
	}
}

func ListsPage(w io.Writer, params listsPageParams) {
	if err := lists.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create List page

func CreateListParams(showAdmin bool, showInstanceAdmin bool) globalWebParams {
	return globalWebParams{
		ShowNavbar:        true,
		ShowAdmin:         showAdmin,
		ShowInstanceAdmin: showInstanceAdmin,
		ChunkName:         "listCreate",
	}
}

func CreateListPage(w io.Writer, params globalWebParams) {
	if err := createList.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Edit List page

type editListParams struct {
	List           constants.List
	SharedListPath string
	globalWebParams
}

func EditListParams(list constants.List, showAdmin bool, showInstanceAdmin bool) editListParams {
	return editListParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "listEdit",
		},
		List:           list,
		SharedListPath: constants.SHARED_LIST_PATH,
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

func ListItemsPageParams(list constants.List, items []constants.Item, showAdmin bool, showInstanceAdmin bool) listItemsPageParams {
	return listItemsPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "listItems",
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

func CreateItemParams(list constants.List, showAdmin bool, showInstanceAdmin bool) createEditItemParams {
	return createEditItemParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "itemCreate",
		},
		List:     list,
		Item:     constants.Item{},
		EditMode: false,
	}
}

func EditItemParams(list constants.List, item constants.Item, showAdmin bool, showInstanceAdmin bool) createEditItemParams {
	return createEditItemParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "itemCreate",
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

// Shared List Items page

type sharedListItemsPageParams struct {
	List      constants.List
	Items     []constants.Item
	ShareCode string
	globalWebParams
}

func SharedListItemsPageParams(shareCode string, list constants.List, items []constants.Item, showAdmin bool, showInstanceAdmin bool) sharedListItemsPageParams {
	return sharedListItemsPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "sharedList",
		},
		List:      list,
		Items:     items,
		ShareCode: shareCode,
	}
}

func SharedListItemsPage(w io.Writer, params sharedListItemsPageParams) {
	if err := sharedList.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Shared List 404 page

type sharedList404PageParams struct {
	ShareCode string
	globalWebParams
}

func SharedList404PageParams(shareCode string) sharedList404PageParams {
	return sharedList404PageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        false,
			ShowAdmin:         false,
			ShowInstanceAdmin: false,
			ChunkName:         "sharedList404",
		},
		ShareCode: shareCode,
	}
}

func SharedList404Page(w io.Writer, params sharedList404PageParams) {
	if err := sharedList404.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// User Admin page

type userAdminPageParams struct {
	Users  []constants.UserRead
	SelfId int
	globalWebParams
}

func UserAdminPageParams(users []constants.UserRead, selfId int, showAdmin bool, showInstanceAdmin bool) userAdminPageParams {
	return userAdminPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "userAdmin",
		},
		Users:  users,
		SelfId: selfId,
	}
}

func UserAdminPage(w io.Writer, params userAdminPageParams) {
	if err := userAdmin.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// All Users page (Instance Admin)

type allUsersPageParams struct {
	Users  []constants.UserRead
	SelfId int
	globalWebParams
}

func AllUsersPageParams(users []constants.UserRead, selfId int, showInstanceAdmin bool) allUsersPageParams {
	return allUsersPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         true,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "allUsers",
		},
		Users:  users,
		SelfId: selfId,
	}
}

func AllUsersPage(w io.Writer, params allUsersPageParams) {
	if err := allUsers.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create User page

type createUserPageParams struct {
	globalWebParams
	GroupAdmins []constants.UserRead
}

func CreateUserParams(showAdmin bool, showInstanceAdmin bool, groupAdmins []constants.UserRead) createUserPageParams {
	return createUserPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        true,
			ShowAdmin:         showAdmin,
			ShowInstanceAdmin: showInstanceAdmin,
			ChunkName:         "userCreate",
		},
		GroupAdmins: groupAdmins,
	}
}

func CreateUserPage(w io.Writer, params createUserPageParams) {
	if err := userCreate.Execute(w, params); err != nil {
		log.Print(err)
	}
}
