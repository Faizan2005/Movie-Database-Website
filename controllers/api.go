package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Faizan2005/Movie-Database/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

//type APIFunc func(w http.ResponseWriter, r *http.Request) error

/*func makeHTTPHandlerFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			panic(err)
		}
	}
}*/

func makeHTTPHandlerFunc(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling %s %s", r.Method, r.URL.Path)
		if err := fn(w, r); err != nil {
			log.Printf("Error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)

}

type APIServer struct {
	listenAddr string
	client     *mongo.Client
	movies     *mongo.Collection
	user       *mongo.Collection
	// embeddedmovies *mongo.Collection
}

func NewAPIServer(listenAddr string, client *mongo.Client) *APIServer {

	db1 := client.Database("sample_mflix")

	db2 := client.Database("credDB")

	return &APIServer{
		listenAddr: listenAddr,
		client:     client,
		movies:     db1.Collection("movies"),
		user:       db2.Collection("user"),
		//	embeddedmovies: db.Collection("embedded_movies"),
	}
}

func (s *APIServer) Run() {

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://0.0.0.0:8001"}, // Allow your frontend origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Allow specific methods
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allow specific headers
	})

	// Use the CORS handler
	handler := c.Handler(router)

	router.HandleFunc("/user/signup", makeHTTPHandlerFunc(s.handleUserSignup)).Methods("POST")
	router.HandleFunc("/user/login", makeHTTPHandlerFunc(s.handleUserLogin)).Methods("POST")

	router.Handle("/admin/movie", s.jwtMiddleware(s.adminMiddleware((makeHTTPHandlerFunc(s.handleMovie))))).Methods("GET", "POST")
	router.Handle("/admin/movie/{id}", s.jwtMiddleware(s.adminMiddleware((makeHTTPHandlerFunc(s.handleMovieByID))))).Methods("GET", "PUT", "DELETE")

	router.Handle("/search", s.jwtMiddleware(makeHTTPHandlerFunc(s.handleSearchMovies))).Methods("GET")

	router.Handle("/movie/{id}/playback", s.jwtMiddleware(makeHTTPHandlerFunc(s.handleMoviePlayback))).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running"))
	}).Methods("GET")

	log.Println("Starting server on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, handler)
}

func (s *APIServer) handleMovie(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetMovies(w, r)
	}

	if r.Method == "POST" {
		return s.handleAddMovie(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleMovieByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetMovie(w, r)
	}

	if r.Method == "PUT" {
		return s.handleUpdateMovie(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteMovie(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetMovie(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	var movie bson.M

	err = s.movies.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		return fmt.Errorf("movie not found %v", err)
	}

	return WriteJSON(w, http.StatusOK, movie)
}

func (s *APIServer) handleAddMovie(w http.ResponseWriter, r *http.Request) error {

	var movie utils.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		return fmt.Errorf("invalid input %v", err)
	}

	result, err := s.movies.InsertOne(r.Context(), movie)
	if err != nil {
		return fmt.Errorf("invalid input %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "movie added successfully",
		"movie_id": result.InsertedID,
	})
}

func (s *APIServer) handleUpdateMovie(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	var updates bson.M
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		return fmt.Errorf("failed to decode request body: %v", err)
	}

	_, err = s.movies.UpdateOne(r.Context(), bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return fmt.Errorf("failed to update movie: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Movie updated successfully",
	})
}

func (s *APIServer) handleDeleteMovie(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	_, err = s.movies.DeleteOne(r.Context(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}

	return WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Movie deleted successfully",
	})
}

func (s *APIServer) handleGetMovies(w http.ResponseWriter, r *http.Request) error {

	cursor, err := s.movies.Find(r.Context(), bson.M{})

	if err != nil {
		return fmt.Errorf("failed to retrieve movies: %v", err)
	}

	defer cursor.Close(r.Context())

	var movies []utils.Movie
	if err := cursor.All(r.Context(), &movies); err != nil {
		return fmt.Errorf("failed to parse movies: %v", err)
	}

	return WriteJSON(w, http.StatusOK, movies)
}

func (s *APIServer) handleSearchMovies(w http.ResponseWriter, r *http.Request) error {
	// Get the query parameter 'q' from the URL
	query := r.URL.Query().Get("q")
	if query == "" {
		return fmt.Errorf("query parameter 'q' is required")
	}

	// Create the search query using Atlas Search
	filter := bson.D{
		{
			Key: "$search", Value: bson.M{
				"index": "movie-search-index", // Replace with your index name
				"autocomplete": bson.M{
					"query":      query,
					"path":       "title",      // Field to search in
					"tokenOrder": "sequential", // To match words in order
					"fuzzy": bson.M{
						"maxEdits":      2,   // Allow 2 character differences
						"prefixLength":  1,   // First character cannot change
						"maxExpansions": 256, // Max expansions for fuzzy terms
					},
				},
			},
		},
	}

	// Perform the aggregation query
	cursor, err := s.movies.Aggregate(r.Context(), mongo.Pipeline{filter})
	if err != nil {
		return fmt.Errorf("failed to search movies: %v", err)
	}
	defer cursor.Close(r.Context())

	// Decode the results into a slice of movies
	var movies []bson.M
	if err := cursor.All(r.Context(), &movies); err != nil {
		return fmt.Errorf("failed to decode movies: %v", err)
	}

	// Return the list of movies as a JSON response
	return WriteJSON(w, http.StatusOK, movies)
}

const TMDB_API_KEY = "19a62d9356a6e640e067a6c79b07dade"

func (s *APIServer) handleMoviePlayback(w http.ResponseWriter, r *http.Request) error {
	movieID := mux.Vars(r)["id"]
	log.Printf("Received request for movie ID: %s", movieID) // Log movie ID

	// Convert movieID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(movieID)
	if err != nil {
		log.Printf("Invalid movie ID: %v", err)
		return fmt.Errorf("invalid movie ID: %v", err)
	}

	// Fetch the movie from MongoDB
	var movie bson.M
	err = s.movies.FindOne(r.Context(), bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		log.Printf("Movie not found: %v", err)
		return fmt.Errorf("movie not found: %v", err)
	}

	// Extract IMDb ID
	imdbMap, ok := movie["imdb"].(primitive.M)
	if !ok {
		return fmt.Errorf("IMDb field is not a valid map")
	}

	imdbID, ok := imdbMap["id"].(int32)
	if !ok || imdbID == 0 {
		return fmt.Errorf("IMDb ID not found or is invalid")
	}

	imdbIDStr := fmt.Sprintf("tt%d", int(imdbID))
	log.Printf("Fetching trailer from TMDb for IMDb ID: %s", imdbIDStr)

	// Step 1: Get TMDb movie ID using IMDb ID
	tmdbFindURL := fmt.Sprintf("https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=imdb_id", imdbIDStr, TMDB_API_KEY)
	resp, err := http.Get(tmdbFindURL)
	if err != nil {
		log.Printf("Failed to fetch TMDb movie details: %v", err)
		return fmt.Errorf("failed to fetch TMDb movie details: %v", err)
	}
	defer resp.Body.Close()

	var tmdbFindResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tmdbFindResponse); err != nil {
		log.Printf("Failed to parse TMDb response: %v", err)
		return fmt.Errorf("failed to parse TMDb response: %v", err)
	}

	// Log the entire TMDb Find Response for debugging
	log.Printf("TMDb Find Response: %+v", tmdbFindResponse)

	// Extract TMDb movie ID from the response
	movieResults, ok := tmdbFindResponse["movie_results"].([]interface{})
	if !ok || len(movieResults) == 0 {
		return fmt.Errorf("no movie found for IMDb ID: %s", imdbIDStr)
	}
	tmdbMovieID := int(movieResults[0].(map[string]interface{})["id"].(float64))

	// Step 2: Get trailer using TMDb movie ID
	tmdbVideosURL := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/videos?api_key=%s&language=en-US", tmdbMovieID, TMDB_API_KEY)
	videoResp, err := http.Get(tmdbVideosURL)
	if err != nil {
		log.Printf("Failed to fetch TMDb trailers: %v", err)
		return fmt.Errorf("failed to fetch TMDb trailers: %v", err)
	}
	defer videoResp.Body.Close()

	var tmdbVideoResponse map[string]interface{}
	if err := json.NewDecoder(videoResp.Body).Decode(&tmdbVideoResponse); err != nil {
		log.Printf("Failed to parse TMDb video response: %v", err)
		return fmt.Errorf("failed to parse TMDb video response: %v", err)
	}

	// Log the TMDb video response to see the structure
	log.Printf("TMDb Video Response: %+v", tmdbVideoResponse)

	// Extract "results" from TMDb response
	results, ok := tmdbVideoResponse["results"].([]interface{})
	if !ok || len(results) == 0 {
		return fmt.Errorf("no videos found for movie ID: %d", tmdbMovieID)
	}

	var trailerURL string
	// Iterate over the results to find the "Trailer" type video
	for _, result := range results {
		video := result.(map[string]interface{})
		if video["type"] == "Trailer" && video["site"] == "YouTube" {
			trailerURL = fmt.Sprintf("https://www.youtube.com/embed/%s", video["key"])
			break
		}
	}

	if trailerURL == "" {
		return fmt.Errorf("no YouTube trailer found for movie ID: %d", tmdbMovieID)
	}

	// Return the trailer URL
	response := map[string]string{
		"trailer_url": trailerURL,
	}
	return WriteJSON(w, http.StatusOK, response)
}

var SECRET_KEY = []byte("1234567890")

func GenerateJWT(email string, isAdmin bool) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	// Create a token with claims
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["isAdmin"] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Set expiry time

	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation:", err)
		return "", err
	}
	return tokenString, nil
}

func (s *APIServer) handleUserSignup(w http.ResponseWriter, r *http.Request) error {
	var user utils.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user signup request: %v", err)
		return fmt.Errorf("error decoding: %v", err)
	}

	// Check if user already exists
	var existingUser utils.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.user.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		log.Printf("User  already exists with email: %s", user.Email)
		return fmt.Errorf("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return fmt.Errorf("error hashing password: %v", err)
	}

	user.Password = string(hashedPassword)

	user.IsAdmin = false

	result, err := s.user.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error while adding user credentials: %v", err)
		return fmt.Errorf("error while adding user credentials: %v", err)
	}

	log.Printf("User  signed up successfully: %v", result.InsertedID)
	return WriteJSON(w, http.StatusOK, result)
}

func (s *APIServer) handleUserLogin(w http.ResponseWriter, r *http.Request) error {
	var user utils.User
	var userDB utils.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding user login request: %v", err)
		return fmt.Errorf("error decoding: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.user.FindOne(ctx, bson.M{"email": user.Email}).Decode(&userDB)
	if err != nil {
		log.Printf("User  not found: %v", err)
		return fmt.Errorf("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
	if err != nil {
		log.Printf("Password mismatch for user: %s", user.Email)
		return fmt.Errorf("invalid password")
	}

	jwtToken, err := GenerateJWT(userDB.Email, userDB.IsAdmin)
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		return fmt.Errorf("error generating token: %v", err)
	}

	log.Printf("User  logged in successfully: %s", user.Email)
	return WriteJSON(w, http.StatusOK, map[string]string{"token": jwtToken})
}

func (s *APIServer) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "authorization header is missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "authorization header is missing", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return SECRET_KEY, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !claims["isAdmin"].(bool) {
			http.Error(w, "Access denied: admin only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
