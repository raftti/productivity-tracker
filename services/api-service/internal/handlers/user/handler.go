package user

import (
	pb "api-service/pkg/proto"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	dbClient pb.DBServiceClient
}

const (
	PATH_USER_BY_ID = "/user/{id}"
	PATH_USER = "/user"
)

func NewUserHandler(dbClient pb.DBServiceClient) *UserHandler {
	return &UserHandler{dbClient: dbClient}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Get(PATH_USER_BY_ID, func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
	
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
	
		user, err := h.dbClient.GetUser(r.Context(), &pb.UserRequest{Id: int32(idInt)})
		if err != nil {
			log.Printf("Error calling GetUser: %v", err)
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
			return
		}
	
		marshaledUser, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user: %v", err)
			http.Error(w, "Error marshaling user", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaledUser)
	})

	r.Post(PATH_USER, func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Name string `json:"name"`
			Email string `json:"email"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := h.dbClient.CreateUser(r.Context(), &pb.CreateUserRequest{Name: reqBody.Name, Email: reqBody.Email})
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		marshaledUser, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user: %v", err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(marshaledUser)
	})
}
