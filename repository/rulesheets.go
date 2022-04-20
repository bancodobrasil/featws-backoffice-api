package repository

import (
	"context"

	"github.com/bancodobrasil/featws-api/database"
	"github.com/bancodobrasil/featws-api/models"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("Error on insert the result into collection: %v", err)
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
		log.Errorf("Error on find the results from collection: %v", err)
		return
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var rulesheet *models.Rulesheet
		if err = results.Decode(&rulesheet); err != nil {
			log.Errorf("Error on decode the result into rulesheet: %v", err)
			return
		}

		list = append(list, rulesheet)
	}

	return
}

//TODO PERGUNTAR PRO RAPHA
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
		log.Errorf("Error on find one result into collection: %v", err)
		//TODO PERGUNTAR PRO RAPHA
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
		log.Errorf("Error on update one into collection: %v", err)
		return
	}

	updated, err = r.Get(ctx, entity.ID.Hex())
	if err != nil {
		log.Errorf("Error on get the updated result into collection: %v", err)
		return
	}

	return
}

// Delete ...
func (r Rulesheets) Delete(ctx context.Context, id string) (deleted bool, err error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Error on transform the oid in Object ID: %v", err)
		return
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": oid})

	if err != nil {
		log.Errorf("Error on delete one result from colletion: %v", err)
		return
	}

	deleted = result.DeletedCount == 1

	return
}
