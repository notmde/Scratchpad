package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = os.Getenv("SECRET_KEY")

var store, err = newDBStore()

func main() {

	if err != nil {
		log.Println(err)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/api/signup", signUp)
	r.Post("/api/login", login)
	r.Post("/api/logout", verifyAuth(logout))
	r.Patch("/api/update", verifyAuth(userUpdate))
	r.Delete("/api/delete", verifyAuth(Deletion))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8080/frontend/", "http://127.0.0.1:8080", "http://localhost:8080/frontend/", "http://localhost:8080"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		//Debug: true,
	})

	http.ListenAndServe(":8081", c.Handler(r))
}

func signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	user := new(User)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)

	if len(user.ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user id can't be empty"})
		return
	}

	if len(user.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "password can't be empty"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	if err = store.insertDB(user); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user already exists"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "user created"})
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	body := new(User)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &body)

	user := new(User)

	if err := store.getDB(body.ID, user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid password"})
		return
	}

	// jwt token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  body.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	//tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	tokenString, err := token.SignedString([]byte(SECRET_KEY))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "cannot create tokenstring"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "authentication",
		Value:    tokenString,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged in", "canvas_data": user.CanvasData})
}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user := new(User)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)

	err := store.updateCanvasDB(r.Header.Get("id"), user.CanvasData)

	if err != nil {
		log.Println(err)
		log.Println(r.Header.Get("id"))
		log.Println(user.CanvasData)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to update data"})
		return
	}

	// delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "authentication",
		Value:    "Logout",
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "logged out"})
}

func userUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user := new(User)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	err = store.updatePasswordDB(r.Header.Get("id"), user.Password)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to update data"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func Deletion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user := new(User)

	reqBody, _ := io.ReadAll(r.Body)
	json.Unmarshal(reqBody, &user)

	if err := store.deleteDB(r.Header.Get("id")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user_id does not exist"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
