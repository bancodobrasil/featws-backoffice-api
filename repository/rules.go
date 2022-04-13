package repository

import (
	"context"

	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Rules ...
type Rules struct {
	collection *mongo.Collection
}

var instance = Rules{}

// GetRulesRepository ...
func GetRulesRepository() Rules {
	if instance.collection == nil {
		instance.collection = database.GetCollection("rules")
	}

	return instance
}

// Create ...
func (r Rules) Create(ctx context.Context, rule *models.Rule) error {

	result, err := r.collection.InsertOne(ctx, rule)
	if err != nil {
		return err
	}

	rule.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// Find ...
func (r Rules) Find(ctx context.Context, filter interface{}) (list []models.Rule, err error) {

	if filter == nil {
		filter = bson.M{}
	}

	results, err := r.collection.Find(ctx, filter)
	if err != nil {
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var rule models.Rule
		if err = results.Decode(&rule); err != nil {
			return
		}

		list = append(list, rule)
	}

	return
}

// Get ...
func (r Rules) Get(ctx context.Context, id string) (rule *models.Rule, err error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	result := r.collection.FindOne(ctx, bson.M{"_id": oid})

	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}

	result.Decode(&rule)

	return
}

// Update ...
func (r Rules) Update(ctx context.Context, entity models.Rule) (updated *models.Rule, err error) {

	//update := bson.M{"name": entity.Name, "type": entity.Type, "opti": entity.Headers}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": entity.ID}, bson.M{"$set": entity})

	if err != nil {
		return
	}

	updated, err = r.Get(ctx, entity.ID.Hex())
	if err != nil {
		return
	}

	return
}

// Delete ...
func (r Rules) Delete(ctx context.Context, id string) (deleted bool, err error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		return
	}

	deleted = result.DeletedCount == 1

	return
}
