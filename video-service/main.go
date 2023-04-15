package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	_ "strconv"
	"time"

	"github.com/gorilla/mux"
)

// Video struct represents a video object
type Video struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string             `bson:"title,omitempty"`
	Genre string             `bson:"genre,omitempty"`
	Age   int                `bson:"age,omitempty"`
}

// MongoDB database and collection names
const (
	dbName         = "videos"
	collectionName = "videos"
)

// MongoDB client and collection instances
var (
	client    *mongo.Client
	videoColl *mongo.Collection
)

func main() {
	// Initialize the MongoDB client and collection instances
	connectMongoDB()

	// Create a new router
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/videos", getVideos).Methods("GET")
	router.HandleFunc("/videos/{id}", getVideo).Methods("GET")
	router.HandleFunc("/videos", createVideo).Methods("POST")
	//router.HandleFunc("/videos/{id}", updateVideo).Methods("PUT")
	//router.HandleFunc("/videos/{id}", deleteVideo).Methods("DELETE")
	//router.HandleFunc("/videos/search", searchVideos).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// connectMongoDB initializes the MongoDB client and collection instances
func connectMongoDB() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the video collection instance
	videoColl = client.Database(dbName).Collection(collectionName)
}

// getVideos returns all the videos
func getVideos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := videoColl.Find(ctx, bson.M{})
	if err != nil {
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get videos from database")
		}
	}(cursor, ctx)
	var videos []Video
	if err = cursor.All(ctx, &videos); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode videos from database")
		return
	}
	respondWithJSON(w, http.StatusOK, videos)
}

func response(w http.ResponseWriter, statusCode int, payload interface{}, contentType string) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(statusCode)
	switch payload.(type) {
	case string:
		_, err := fmt.Fprint(w, payload)
		if err != nil {
			return
		}
	default:
		jsonResponse, err := json.Marshal(payload)
		if err != nil {
			response(w, http.StatusInternalServerError, "Failed to serialize response", "")
			return
		}
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response(w, statusCode, payload, "application/json")
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response(w, statusCode, message, "")
}

// getVideo returns a specific video by ID
func getVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	videoID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	var video Video
	err = videoColl.FindOne(context.Background(), bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(video)
	if err != nil {
		return
	}
}

// createVideo creates a new video
func createVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var video Video
	err := json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		log.Fatal(err)
	}

	result, err := videoColl.InsertOne(context.Background(), video)
	if err != nil {
		log.Fatal(err)
	}

	newID := result.InsertedID.(primitive.ObjectID)
	video.ID = newID

	err = json.NewEncoder(w).Encode(video)
	if err != nil {
		return
	}
}
