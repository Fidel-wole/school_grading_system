package database

import (
	"context"
	"fmt"
	//"sort"
	"time"

	"cloudnotte_practice/graph/model"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func (db *DB) AddSubject(name string) (*model.Subject, error) {
	if name == "" {
		return nil, fmt.Errorf("term name cannot be empty")
	}
	subjectCollec := db.client.Database("school_management_system").Collection("subjects")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	subjectDetails := model.Subject{
		Name: name,
	}

	res, err := subjectCollec.InsertOne(ctx, subjectDetails)
	if err != nil {
		return nil, fmt.Errorf("error adding term: %v", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		subjectDetails.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}
	return &subjectDetails, nil
}


func (db *DB) GetSubjectByID(subjectId string) (*model.Subject, error) {
	// Convert string ID to MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {
		return nil, fmt.Errorf("invalid subject ID: %v", err)
	}

	// Temporary struct to map MongoDB `_id` to Go's ObjectID type
	var mongoSubject struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}

	// Retrieve the subject from MongoDB
	err = db.client.Database("school_management_system").
		Collection("subjects").
		FindOne(context.Background(), bson.M{"_id": objectId}).
		Decode(&mongoSubject)

	if err != nil {
		return nil, fmt.Errorf("error finding subject: %v", err)
	}

	// Convert MongoDB ObjectID to string for the gqlgen Subject struct
	subject := &model.Subject{
		ID:   mongoSubject.ID.Hex(),
		Name: mongoSubject.Name,
	}

	return subject, nil
}