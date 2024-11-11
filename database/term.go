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

func (db *DB) AddTerm(name string) (*model.Term, error) {
	if name == "" {
		return nil, fmt.Errorf("term name cannot be empty")
	}

	termCollec := db.client.Database("school_management_system").Collection("terms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	termDetails := model.Term{
		Name: name,
	}

	res, err := termCollec.InsertOne(ctx, termDetails)
	if err != nil {
		return nil, fmt.Errorf("error adding term: %v", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		termDetails.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}

	return &termDetails, nil
}

// Fetch Term by ID
func (db *DB) GetTermByID(termID primitive.ObjectID) (*model.Term, error) {
	term := &model.Term{}
	err := db.client.Database("school_management_system").Collection("terms").FindOne(context.Background(), bson.M{"_id": termID}).Decode(term)
	if err != nil {
		return nil, fmt.Errorf("error finding term: %v", err)
	}
	return term, nil
}