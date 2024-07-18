package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	filter "github.com/dlion/faceit_challenge/internal"
)

func (u *UserHandler) GetUsersHandler(w http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()

	filter := NewUserFilterFromQuery(queryParams)

	log.Printf("Getting Users with limit %d and offset %d", *filter.Limit, *filter.Offset)

	paginatedUsers, err := u.UserService.GetUsers(req.Context(), filter)
	if err != nil {
		log.Print("Can't get paginated users")
		http.Error(w, "Can't get users", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(paginatedUsers); err != nil {
		log.Print(err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func NewUserFilterFromQuery(query url.Values) *filter.UserFilter {
	fbuilder := filter.NewFilterBuilder()

	firstName := query.Get("first_name")
	if firstName != "" {
		fbuilder.ByFirstName(&firstName)
	}

	lastName := query.Get("last_name")
	if lastName != "" {
		fbuilder = fbuilder.ByLastName(&lastName)
	}

	nickname := query.Get("nickname")
	if nickname != "" {
		fbuilder = fbuilder.ByNickname(&nickname)
	}

	country := query.Get("country")
	if country != "" {
		fbuilder = fbuilder.ByCountry(&country)
	}

	email := query.Get("email")
	if email != "" {
		fbuilder = fbuilder.ByEmail(&email)
	}

	limitStr := query.Get("limit")
	if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
		fbuilder.WithLimit(intToint64(limit))
	} else {
		fbuilder.WithLimit(intToint64(10))
	}

	offsetStr := query.Get("offset")
	if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
		fbuilder.WithOffset(intToint64(offset))
	} else {
		fbuilder.WithOffset(intToint64(0))
	}

	return fbuilder.Build()
}

func intToint64(value int) *int64 {
	int64value := int64(value)
	return &int64value
}
