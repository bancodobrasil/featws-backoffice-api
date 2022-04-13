package models

import (
	v1 "github.com/bancodobrasil/featws-api/payloads/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rulesheet ...
type Rulesheet struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
}

// NewRulesheetV1 ...
func NewRulesheetV1(payload v1.Rulesheet) (entity Rulesheet, err error) {

	id := primitive.NilObjectID

	if payload.ID != "" {
		id, err = primitive.ObjectIDFromHex(payload.ID)
		if err != nil {
			return
		}
	}

	entity = Rulesheet{
		ID:   id,
		Name: payload.Name,
	}
	return
}
