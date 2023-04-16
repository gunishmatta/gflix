package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"
)

// Video struct represents a video object
type Video struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Genre       string             `bson:"genre,omitempty"`
	Description string             `bson:"description"`
	AgeRating   int                `bson:"ageRating"`
	CreatedAt   time.Time          `bson:"createdAt"`
	Url         string             `bson:"url"`
}

type CreateVideoRequest struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description"`
	Genre       string                `json:"genre" validate:"required"`
	AgeRating   int                   `json:"age_rating" validate:"required"`
	VideoFile   *multipart.FileHeader `json:"-" validate:"required"`
}
