package models

import (
	v1 "github.com/bancodobrasil/featws-api/payloads/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Rule ...
type Rule struct {
	ID      primitive.ObjectID     `bson:"_id,omitempty"`
	Name    string                 `bson:"name,omitempty"`
	Type    string                 `bson:"type,omitempty"`
	Options map[string]interface{} `bson:"options,omitempty"`
}

// NewRuleV1 ...
func NewRuleV1(payload v1.Rule) (entity Rule, err error) {

	id := primitive.NilObjectID

	if payload.ID != "" {
		id, err = primitive.ObjectIDFromHex(payload.ID)
		if err != nil {
			return
		}
	}

	entity = Rule{
		ID:      id,
		Name:    payload.Name,
		Type:    payload.Type,
		Options: payload.Options,
	}
	return
}
