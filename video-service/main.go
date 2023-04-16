package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	_ "strconv"
	"time"
	"video-service/helpers"

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
	//router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
	//	httpSwagger.URL("localhost:8000/docs/swa"), // URL pointing to the API docs JSON file
	//))
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
func getVideos(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := videoColl.Find(ctx, bson.M{})
	if err != nil {
		return
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to get videos from database")
		}
	}(cursor, ctx)
	var videos []Video
	if err = cursor.All(ctx, &videos); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Failed to decode videos from database")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, videos)
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

// @Summary Create a new video
// @Description Create a new video with the specified title, description, and genre
// @Tags videos
// @Accept json
// @Produce json
// @Param body body Video true "Request body"
// @Success 201 {object} Video
// @Failure 400
// @Failure 500
// @Router /videos [post]
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
