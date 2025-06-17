package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/helper"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/web"
)

func init() {
	// Collection routes
	constants.ROUTER.HandleFunc("/collections", middleware.DefaultMiddlewareChain(collectionsPOST)).Methods("POST")
	constants.ROUTER.HandleFunc("/collections/create", middleware.DefaultMiddlewareChain(createCollectionGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/namecheck", middleware.DefaultMiddlewareChain(collectionNameCheckGET)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}", middleware.Chain(collectionHandler, append([]middleware.Middleware{middleware.CollectionIdOwner("collectionId")}, middleware.DefaultMiddlewareSlice...)...))
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/edit", middleware.Chain(editCollectionGET, append([]middleware.Middleware{middleware.CollectionIdOwner("collectionId")}, middleware.DefaultMiddlewareSlice...)...)).Methods("GET")
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/lists/{listId:[0-9]+}", middleware.Chain(collectionListHandler, append([]middleware.Middleware{middleware.ListIdOwner("listId"), middleware.CollectionIdOwner("collectionId")}, middleware.DefaultMiddlewareSlice...)...))

	// Collection sharing routes
	constants.ROUTER.HandleFunc("/collections/{collectionId:[0-9]+}/share", middleware.Chain(collectionShareHandler, append(middleware.DefaultMiddlewareSlice, middleware.CollectionIdOwner("collectionId"))...))
	constants.ROUTER.HandleFunc("/"+constants.SHARED_COLLECTION_PATH+"/{shareCode}", middleware.DefaultPublicMiddlewareChain(sharedCollectionGET)).Methods("GET")
}

// collectionsPOST handles POST requests for /collections
func collectionsPOST(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	var collectionName string = r.FormValue("name")
	//Don't trust client input

	taken, err := database.CollectionNameTaken(userId, collectionName)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	if taken {
		http.Error(w, "You already have a collection with that name", http.StatusConflict)
		return
	}

	var description string = r.FormValue("description")
	newId, err := database.CreateCollection(userId, collectionName, description)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/collections/%d", newId))
	w.WriteHeader(http.StatusOK)
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
	userId, err := helper.GetUserId(r)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}
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

	listIdsInCollection, err := database.GetCollectionListIds(collectionId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	allLists, err := database.GetLists(userId)
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Return JSON if requested, otherwise render the HTML page
	if wantsJSON {
		// Construct a response with both collection details and lists
		type CollectionResponse struct {
			Collection constants.Collection `json:"collection"`
			Lists      []uint64             `json:"lists"`
		}

		response := CollectionResponse{
			Collection: collection,
			Lists:      listIdsInCollection,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			log.Print(err)
		}
	} else {
		// Get list IDs that have share codes for this user
		listIdsWithShareCode, err := database.GetListIdsWithShareCode(userId)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}

		admin := helper.IsUserAdmin(r)
		instanceAdmin := helper.IsUserInstanceAdmin(r)

		// Render the collection detail page
		collectionDetailPage := web.CollectionDetailPageParams(
			r,
			collection,
			listIdsInCollection,
			listIdsWithShareCode,
			allLists,
			admin,
			instanceAdmin,
		)
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

	// If not in collection, add it
	if !inCollection {
		// First check if the collection has a share code
		collection, err := database.GetCollection(collectionId)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}

		// If collection has a share code, check if list already has a share code
		if collection.ShareCode.Valid && len(collection.ShareCode.String) > 0 {
			list, err := database.GetList(listId)
			if err != nil {
				http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
				log.Print(err)
				return
			}

			// If list doesn't have a share code, generate one
			if !list.ShareCode.Valid || len(list.ShareCode.String) == 0 {
				_, err = database.GenerateShareCode(listId)
				if err != nil {
					http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
					log.Print(err)
					return
				}
			}
		}

		// Now add the list to the collection
		err = database.AddListToCollection(collectionId, listId)
		if err != nil {
			http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
			log.Print(err)
			return
		}
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
	params := web.CreateCollectionParams(r, admin, instanceAdmin)
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
	editParams := web.EditCollectionParams(r, collection, admin, instanceAdmin)
	web.EditCollectionPage(w, editParams)
}

// collectionShareHandler handles requests for /collections/{collectionId}/share
func collectionShareHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT", "POST": // Accept both PUT (API) and POST (form submission)
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
	sharedCollectionPage := web.SharedCollectionPageParams(r, shareCode, collection, lists, admin, instanceAdmin)
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
