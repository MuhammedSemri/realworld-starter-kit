package main

import (
	hfn "chilledoj/myreal/handlerfn"
	"chilledoj/myreal/models"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func buildRouter(db *models.AppDB, logger *log.Logger) http.Handler {

	env := hfn.AppEnvironment{DB: db, Logger: logger}

	na := notImplemented{}
	r := mux.NewRouter()

	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api := r.PathPrefix("/api").Subrouter()
	// TODO
	// Unprotected Routes
	api.Handle("/users", hfn.Register(&env)).Methods("POST").Name("Registration")
	api.Handle("/users/login", hfn.Login(&env)).Methods("POST").Name("Authentication")
	api.Handle("/tags", na).Methods("GET").Name("Get Tags")
	api.Handle("/articles/{slug}", na).Methods("GET").Name("Get Article")

	// TODO
	// OPTIONAL Auth - What does that mean how to implement?
	// First thought is to write 2 sets of middleware
	//   JWT2Ctx - will validate JWT if present and will add user from valid JWT into request context
	//   MustHaveUser - will check the request context for the presence of a user and will redirect if not present
	//
	api.Handle("/articles/{slug}/comments", na).Methods("GET").Name("Get Comments from an Article")
	api.Handle("/articles", na).Methods("GET").Name("List Articles")
	api.Handle("/profiles/{username}", na).Methods("GET").Name("Get Profile")

	// TODO
	// Protected Routes
	api.Handle("/user", hfn.MustHaveUser(hfn.GetCurrentUser(&env).(hfn.AppHandler))).Methods("GET").Name("Get Current User")
	api.Handle("/user", hfn.MustHaveUser(hfn.UpdateUser(&env).(hfn.AppHandler))).Methods("PUT").Name("Update User")
	api.Handle("/profiles/{username}/follow", na).Methods("POST").Name("Follow User")
	api.Handle("/profiles/{username}/follow", na).Methods("Delete").Name("Delete User")
	api.Handle("/articles/feed", na).Methods("GET").Name("Feed Articles")
	api.Handle("/articles", na).Methods("POST").Name("Create Article")
	api.Handle("/articles/{slug}", na).Methods("PUT").Name("Update Article")
	api.Handle("/articles/{slug}", na).Methods("DELETE").Name("Delete Article")
	api.Handle("/articles/{slug}/favourite", na).Methods("POST").Name("Favourite Article")
	api.Handle("/articles/{slug}/favourite", na).Methods("DELETE").Name("Unavourite Article")
	api.Handle("/articles/{slug}/comments", na).Methods("POST").Name("Add Comments to an Article")
	api.Handle("/articles/{slug}/comments/{id}", na).Methods("DELETE").Name("Delete Comment")

	jwtRouter := hfn.Jwt2Ctx{Env: &env, Fn: r}
	return handlers.RecoveryHandler()(handlers.CombinedLoggingHandler(os.Stdout, jwtRouter))
}

type notImplemented struct{}

func (notImplemented) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "NOT IMPLEMENTED")
}
