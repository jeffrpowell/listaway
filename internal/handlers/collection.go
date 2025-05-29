package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	// Collection routes
	constants.ROUTER.HandleFunc("/collections", middleware.DefaultMiddlewareChain(collectionsHandler))
	constants.ROUTER.HandleFunc("/collections/create", middleware.DefaultMiddlewareChain(createCollectionGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/namecheck", middleware.DefaultMiddlewareChain(collectionNameCheckGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}", middleware.Chain(collectionHandler, append(middleware.DefaultMiddlewareSlice, middleware.CollectionIdOwner("collectionId"))...))
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/edit", middleware.Chain(editCollectionGET, append(middleware.DefaultMiddlewareSlice, middleware.CollectionIdOwner("collectionId"))...)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/lists/{listId:[0-9]+}", middleware.Chain(collectionListHandler, append(middleware.DefaultMiddlewareSlice, middleware.CollectionIdOwner("collectionId"), middleware.ListIdOwner("listId"))...))

	// Collection sharing routes
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/share", middleware.Chain(collectionShareHandler, append(middleware.DefaultMiddlewareSlice, middleware.CollectionIdOwner("collectionId"))...))
	constants.ROUTER.HandleFunc("/"+constants.SHARED_COLLECTION_PATH+"/{shareCode}", middleware.DefaultPublicMiddlewareChain(sharedCollectionGET)).Methods("GET")
}

// collectionsHandler handles requests for /collections
func collectionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		collectionsGET(w, r)
	case "POST":
		collectionsPOST(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// collectionsGET handles GET requests for /collections
func collectionsGET(w http.ResponseWriter, r *http.Request) {
	// Check if the Accept header indicates JSON is wanted
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := acceptHeader == "application/json"

	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	collections, err := database.GetCollections(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Return JSON if requested, otherwise render the HTML page
	if wantsJSON {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(collections); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			log.Print(err)
		}
	} else {
		admin := helper.IsUserAdmin(r)
		instanceAdmin := helper.IsUserInstanceAdmin(r)
		
		// Render the collections page
		collectionsPage := web.CollectionsPageParams(collections, admin, instanceAdmin)
		web.CollectionsPage(w, collectionsPage)
	}
}

// collectionsPOST handles POST requests for /collections
func collectionsPOST(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var params constants.CollectionPostParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	taken, err := database.CollectionNameTaken(userId, params.Name)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "You already have a collection with that name", http.StatusConflict)
		return
	}

	newId, err := database.CreateCollection(userId, params.Name, params.Description)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strconv.Itoa(newId)))
}

// collectionHandler handles requests for /collections/{collectionId}
func collectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		collectionGET(w, r)
	case "PUT":
		collectionPUT(w, r)
	case "DELETE":
		collectionDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// collectionGET handles GET requests for /collections/{collectionId}
func collectionGET(w http.ResponseWriter, r *http.Request) {
	// Check if the Accept header indicates JSON is wanted
	acceptHeader := r.Header.Get("Accept")
	wantsJSON := acceptHeader == "application/json"

	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware

	collection, err := database.GetCollection(collectionId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	lists, err := database.GetCollectionLists(collectionId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Return JSON if requested, otherwise render the HTML page
	if wantsJSON {
		// Construct a response with both collection details and lists
		type CollectionResponse struct {
			Collection constants.Collection       `json:"collection"`
			Lists      []constants.CollectionList `json:"lists"`
		}

		response := CollectionResponse{
			Collection: collection,
			Lists:      lists,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			log.Print(err)
		}
	} else {
		admin := helper.IsUserAdmin(r)
		instanceAdmin := helper.IsUserInstanceAdmin(r)
		
		// Render the collection detail page
		collectionDetailPage := web.CollectionDetailPageParams(collection, lists, admin, instanceAdmin)
		web.CollectionDetailPage(w, collectionDetailPage)
	}
}

// collectionPUT handles PUT requests for /collections/{collectionId}
func collectionPUT(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var params constants.CollectionPostParams
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Check if name is already taken by another collection (excluding current collection)
	rows, err := database.GetCollections(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	for _, c := range rows {
		if c.Name == params.Name && int(c.Id) != collectionId {
			http.Error(w, "You already have a collection with that name", http.StatusConflict)
			return
		}
	}

	err = database.UpdateCollection(collectionId, params)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// collectionDELETE handles DELETE requests for /collections/{collectionId}
func collectionDELETE(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware

	confirmationName := r.URL.Query().Get("name")
	if confirmationName == "" {
		http.Error(w, "Confirmation name is required", http.StatusBadRequest)
		return
	}

	success, err := database.DeleteCollection(collectionId, confirmationName)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if !success {
		http.Error(w, "Confirmation name doesn't match", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// collectionListHandler handles requests for /collections/{collectionId}/lists/{listId}
func collectionListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		collectionListPUT(w, r)
	case "DELETE":
		collectionListDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// collectionListPUT handles PUT requests for /collections/{collectionId}/lists/{listId}
func collectionListPUT(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware
	listId, _ := helper.GetPathVarInt(r, "listId")             // Error already checked in middleware

	// Check if list is already in collection
	inCollection, err := database.ListInCollection(collectionId, listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Parse request for display order
	type DisplayOrderParams struct {
		DisplayOrder int `json:"displayOrder"`
	}

	var params DisplayOrderParams
	params.DisplayOrder = 0 // Default value

	// Try to decode the request body into the struct
	if r.ContentLength > 0 {
		err = json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			// If we can't decode, just use the default display order
			log.Printf("Could not decode display order params: %v", err)
		}
	}

	// Add or update list in collection
	if inCollection {
		err = database.UpdateListDisplayOrder(collectionId, listId, params.DisplayOrder)
	} else {
		err = database.AddListToCollection(collectionId, listId, params.DisplayOrder)
	}

	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// collectionListDELETE handles DELETE requests for /collections/{collectionId}/lists/{listId}
func collectionListDELETE(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware
	listId, _ := helper.GetPathVarInt(r, "listId")             // Error already checked in middleware

	err := database.RemoveListFromCollection(collectionId, listId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// createCollectionGET handles GET requests for /collections/create
func createCollectionGET(w http.ResponseWriter, r *http.Request) {
	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	
	// Render the collection creation page
	params := web.CreateCollectionParams(admin, instanceAdmin)
	web.CreateCollectionPage(w, params)
}

// editCollectionGET handles GET requests for /collections/{collectionId}/edit
func editCollectionGET(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware

	collection, err := database.GetCollection(collectionId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Collection not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)
	
	// Render the collection edit page
	editParams := web.EditCollectionParams(collection, admin, instanceAdmin)
	web.EditCollectionPage(w, editParams)
}

// collectionShareHandler handles requests for /collections/{collectionId}/share
func collectionShareHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		collectionSharePUT(w, r)
	case "DELETE":
		collectionShareDELETE(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// collectionSharePUT handles PUT requests for /collections/{collectionId}/share
// Creates a share link for a collection
func collectionSharePUT(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware

	code, err := database.GenerateCollectionShareCode(collectionId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Write([]byte(code))
}

// collectionShareDELETE handles DELETE requests for /collections/{collectionId}/share
// Unpublishes a collection's share link
func collectionShareDELETE(w http.ResponseWriter, r *http.Request) {
	collectionId, _ := helper.GetPathVarInt(r, "collectionId") // Error already checked in middleware

	err := database.UnpublishCollectionShareCode(collectionId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// sharedCollectionGET handles GET requests for /sharedcollection/{shareCode}
// Displays a shared collection view
func sharedCollectionGET(w http.ResponseWriter, r *http.Request) {
	shareCode := mux.Vars(r)["shareCode"]

	collection, err := database.GetCollectionFromShareCode(shareCode)
	if err != nil {
		if err == sql.ErrNoRows {
			sharedCollection404Page := web.SharedCollection404PageParams(shareCode)
			web.SharedCollection404Page(w, sharedCollection404Page)
			return
		}
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	lists, err := database.GetCollectionLists(int(collection.Id))
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	admin := helper.IsUserAdmin(r)
	instanceAdmin := helper.IsUserInstanceAdmin(r)

	// Call the web package function to render the shared collection page
	sharedCollectionPage := web.SharedCollectionPageParams(shareCode, collection, lists, admin, instanceAdmin)
	web.SharedCollectionPage(w, sharedCollectionPage)
}

// collectionNameCheckGET handles GET requests for /collections/namecheck
// Checks if a collection name is already in use
func collectionNameCheckGET(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	taken, err := database.CollectionNameTaken(userId, name)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if taken {
		http.Error(w, "Collection name is already in use", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
