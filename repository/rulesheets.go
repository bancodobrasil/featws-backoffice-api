package repository

import (
	"context"

	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Rulesheets ...
type Rulesheets struct {
	collection *mongo.Collection
}

var instanceRulesheets = Rulesheets{}

// GetRulesheetsRepository ...
func GetRulesheetsRepository() Rulesheets {
	if instanceRulesheets.collection == nil {
		instanceRulesheets.collection = database.GetCollection("rulesheets")
	}

	return instanceRulesheets
}

// Create ...
func (r Rulesheets) Create(ctx context.Context, rulesheet *models.Rulesheet) error {

	result, err := r.collection.InsertOne(ctx, rulesheet)
	if err != nil {
		return err
	}

	rulesheet.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// Find ...
func (r Rulesheets) Find(ctx context.Context, filter interface{}) (list []*models.Rulesheet, err error) {

	if filter == nil {
		filter = bson.M{}
	}

	results, err := r.collection.Find(ctx, filter)
	if err != nil {
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var rulesheet *models.Rulesheet
		if err = results.Decode(&rulesheet); err != nil {
			return
		}

		list = append(list, rulesheet)
	}

	return
}

func buildFilter(id string) interface{} {
	oid, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		return bson.M{"_id": oid}
	}
	return bson.M{"name": id}
}

// Get ...
func (r Rulesheets) Get(ctx context.Context, id string) (rulesheet *models.Rulesheet, err error) {

	result := r.collection.FindOne(ctx, buildFilter(id))

	err = result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return
	}

	result.Decode(&rulesheet)

	return
}

// Update ...
func (r Rulesheets) Update(ctx context.Context, entity models.Rulesheet) (updated *models.Rulesheet, err error) {

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
func (r Rulesheets) Delete(ctx context.Context, id string) (deleted bool, err error) {

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
