package web

import (
	"embed"
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"slices"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

//go:embed dist/*
var staticFiles embed.FS
var (
	registerAdmin       = parseSingleLayout("dist/registerAdmin.html")
	login               = parseSingleLayout("dist/login.html")
	lists               = parseSingleLayout("dist/lists.html")
	createList          = parseSingleLayout("dist/listCreate.html")
	editList            = parseSingleLayout("dist/listEdit.html")
	listItems           = parseSingleLayout("dist/listItems.html")
	createItem          = parseSingleLayout("dist/itemCreate.html")
	sharedList          = parseSingleLayout("dist/sharedList.html")
	sharedList404       = parseSingleLayout("dist/sharedList404.html")
	resetForm           = parseSingleLayout("dist/resetForm.html")
	createCollection    = parseSingleLayout("dist/collectionCreate.html")
	editCollection      = parseSingleLayout("dist/collectionEdit.html")
	collectionDetail    = parseSingleLayout("dist/collectionDetail.html")
	sharedCollection    = parseSingleLayout("dist/sharedCollection.html")
	sharedCollection404 = parseSingleLayout("dist/sharedCollection404.html")
	userAdmin           = parseSingleLayout("dist/userAdmin.html")
	allUsers            = parseSingleLayout("dist/allUsers.html")
	userCreate          = parseSingleLayout("dist/userCreate.html")
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
		tmpl = tmpl.Funcs(template.FuncMap{
			"containsUint64": slices.Contains[[]uint64],
		})

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

func parseSingleLayout(file string) *template.Template {
	return template.Must(minifyTemplates("dist/root.html", file))
}

func isAuthenticated(r *http.Request) bool {
	if r == nil {
		return false
	}
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func newGlobalWebParams(r *http.Request, showNavbar, showAdmin, showInstanceAdmin bool, chunkName string) globalWebParams {
	return globalWebParams{
		ShowNavbar:        showNavbar,
		ShowAdmin:         showAdmin,
		ShowInstanceAdmin: showInstanceAdmin,
		ChunkName:         chunkName,
		IsAuthenticated:   isAuthenticated(r),
	}
}

type globalWebParams struct {
	ShowNavbar        bool
	ShowAdmin         bool
	ShowInstanceAdmin bool
	ChunkName         string
	IsAuthenticated   bool
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

// Password Reset page
type resetFormPageParams struct {
	globalWebParams
	TokenValid bool
}

func ResetFormPage(w io.Writer, tokenValid bool) {
	if err := resetForm.Execute(w, resetFormPageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        false,
			ShowAdmin:         false,
			ShowInstanceAdmin: false,
			ChunkName:         "resetForm",
		},
		TokenValid: tokenValid,
	}); err != nil {
		log.Print(err)
	}
}

// Lists page

type listsPageParams struct {
	Lists                []constants.List
	Collections          []constants.Collection
	GroupSharedLists     []constants.ListSharedWithGroup
	GroupSharingEnabled  bool
	SharedListPath       string
	SharedCollectionPath string
	globalWebParams
}

func ListsPageParams(r *http.Request, lists []constants.List, collections []constants.Collection, groupSharedLists []constants.ListSharedWithGroup, groupSharingEnabled bool, showAdmin bool, showInstanceAdmin bool) listsPageParams {
	return listsPageParams{
		globalWebParams:      newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "lists"),
		Lists:                lists,
		Collections:          collections,
		GroupSharedLists:     groupSharedLists,
		GroupSharingEnabled:  groupSharingEnabled,
		SharedListPath:       constants.SHARED_LIST_PATH,
		SharedCollectionPath: constants.SHARED_COLLECTION_PATH,
	}
}

func ListsPage(w io.Writer, params listsPageParams) {
	if err := lists.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create List page

func CreateListParams(r *http.Request, showAdmin bool, showInstanceAdmin bool) globalWebParams {
	return newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "listCreate")
}

func CreateListPage(w io.Writer, params globalWebParams) {
	if err := createList.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Edit List page

type editListParams struct {
	List                constants.List
	IsOwner             bool
	GroupSharingEnabled bool
	SharedListPath      string
	globalWebParams
}

func EditListParams(r *http.Request, list constants.List, isOwner bool, groupSharingEnabled bool, showAdmin bool, showInstanceAdmin bool) editListParams {
	return editListParams{
		globalWebParams:     newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "listEdit"),
		List:                list,
		IsOwner:             isOwner,
		GroupSharingEnabled: groupSharingEnabled,
		SharedListPath:      constants.SHARED_LIST_PATH,
	}
}

func EditListPage(w io.Writer, params editListParams) {
	if err := editList.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// List Items page

type listItemsPageParams struct {
	List    constants.List
	Items   []constants.Item
	CanEdit bool
	globalWebParams
}

func ListItemsPageParams(r *http.Request, list constants.List, items []constants.Item, canEdit bool, showAdmin bool, showInstanceAdmin bool) listItemsPageParams {
	return listItemsPageParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "listItems"),
		List:            list,
		Items:           items,
		CanEdit:         canEdit,
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

func CreateItemParams(r *http.Request, list constants.List, showAdmin bool, showInstanceAdmin bool) createEditItemParams {
	return createEditItemParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "itemCreate"),
		List:            list,
		Item:            constants.Item{},
		EditMode:        false,
	}
}

func EditItemParams(r *http.Request, list constants.List, item constants.Item, showAdmin bool, showInstanceAdmin bool) createEditItemParams {
	return createEditItemParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "itemCreate"),
		List:            list,
		Item:            item,
		EditMode:        true,
	}
}

func CreateEditItemPage(w io.Writer, params createEditItemParams) {
	if err := createItem.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Shared List Items page

type sharedListItemsPageParams struct {
	List                constants.List
	Items               []constants.Item
	ShareCode           string
	CollectionShareCode string
	HasParentCollection bool
	globalWebParams
}

func SharedListItemsPageParams(r *http.Request, shareCode string, list constants.List, items []constants.Item, showAdmin bool, showInstanceAdmin bool) sharedListItemsPageParams {
	return sharedListItemsPageParams{
		globalWebParams:     newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "sharedList"),
		List:                list,
		Items:               items,
		ShareCode:           shareCode,
		HasParentCollection: false,
		CollectionShareCode: "",
	}
}

// NestedSharedListItemsPageParams creates parameters for a shared list that's being viewed from a parent collection
func NestedSharedListItemsPageParams(r *http.Request, shareCode string, collectionShareCode string, list constants.List, items []constants.Item, showAdmin bool, showInstanceAdmin bool) sharedListItemsPageParams {
	return sharedListItemsPageParams{
		globalWebParams:     newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "sharedList"),
		List:                list,
		Items:               items,
		ShareCode:           shareCode,
		HasParentCollection: true,
		CollectionShareCode: collectionShareCode,
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

func UserAdminPageParams(r *http.Request, users []constants.UserRead, selfId int, showAdmin bool, showInstanceAdmin bool) userAdminPageParams {
	return userAdminPageParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "userAdmin"),
		Users:           users,
		SelfId:          selfId,
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

func AllUsersPageParams(r *http.Request, users []constants.UserRead, selfId int, showInstanceAdmin bool) allUsersPageParams {
	return allUsersPageParams{
		globalWebParams: newGlobalWebParams(r, true, true, showInstanceAdmin, "allUsers"),
		Users:           users,
		SelfId:          selfId,
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

func CreateUserParams(r *http.Request, showAdmin bool, showInstanceAdmin bool, groupAdmins []constants.UserRead) createUserPageParams {
	return createUserPageParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "userCreate"),
		GroupAdmins:     groupAdmins,
	}
}

func CreateUserPage(w io.Writer, params createUserPageParams) {
	if err := userCreate.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Create Collection page
func CreateCollectionParams(r *http.Request, showAdmin bool, showInstanceAdmin bool) globalWebParams {
	return newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "collectionCreate")
}

func CreateCollectionPage(w io.Writer, params globalWebParams) {
	if err := createCollection.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Edit Collection page
type editCollectionParams struct {
	Collection           constants.Collection
	SharedCollectionPath string
	globalWebParams
}

func EditCollectionParams(r *http.Request, collection constants.Collection, showAdmin bool, showInstanceAdmin bool) editCollectionParams {
	return editCollectionParams{
		globalWebParams:      newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "collectionEdit"),
		Collection:           collection,
		SharedCollectionPath: constants.SHARED_COLLECTION_PATH,
	}
}

func EditCollectionPage(w io.Writer, params editCollectionParams) {
	if err := editCollection.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Collection Detail page
type collectionDetailPageParams struct {
	Collection           constants.Collection
	ListIdsInCollection  []uint64
	ListIdsWithShareCode []uint64
	AllLists             []constants.List
	SharedListPath       string
	globalWebParams
}

func CollectionDetailPageParams(r *http.Request, collection constants.Collection, listIdsInCollection []uint64, listIdsWithShareCode []uint64, allLists []constants.List, showAdmin bool, showInstanceAdmin bool) collectionDetailPageParams {
	return collectionDetailPageParams{
		globalWebParams:      newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "collectionDetail"),
		Collection:           collection,
		ListIdsInCollection:  listIdsInCollection,
		ListIdsWithShareCode: listIdsWithShareCode,
		AllLists:             allLists,
		SharedListPath:       constants.SHARED_LIST_PATH,
	}
}

func CollectionDetailPage(w io.Writer, params collectionDetailPageParams) {
	if err := collectionDetail.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Shared Collection page
type sharedCollectionPageParams struct {
	Collection constants.Collection
	Lists      []constants.List
	ShareCode  string
	globalWebParams
}

func SharedCollectionPageParams(r *http.Request, shareCode string, collection constants.Collection, lists []constants.List, showAdmin bool, showInstanceAdmin bool) sharedCollectionPageParams {
	return sharedCollectionPageParams{
		globalWebParams: newGlobalWebParams(r, true, showAdmin, showInstanceAdmin, "sharedCollection"),
		Collection:      collection,
		Lists:           lists,
		ShareCode:       shareCode,
	}
}

func SharedCollectionPage(w io.Writer, params sharedCollectionPageParams) {
	if err := sharedCollection.Execute(w, params); err != nil {
		log.Print(err)
	}
}

// Shared Collection 404 page
type sharedCollection404PageParams struct {
	ShareCode string
	globalWebParams
}

func SharedCollection404PageParams(shareCode string) sharedCollection404PageParams {
	return sharedCollection404PageParams{
		globalWebParams: globalWebParams{
			ShowNavbar:        false,
			ShowAdmin:         false,
			ShowInstanceAdmin: false,
			ChunkName:         "sharedCollection404",
		},
		ShareCode: shareCode,
	}
}

func SharedCollection404Page(w io.Writer, params sharedCollection404PageParams) {
	if err := sharedCollection404.Execute(w, params); err != nil {
		log.Print(err)
	}
}
