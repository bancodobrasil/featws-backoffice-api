package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bancodobrasil/featws-api/dtos"
	payloads "github.com/bancodobrasil/featws-api/payloads/v1"
	responses "github.com/bancodobrasil/featws-api/responses/v1"
	"github.com/bancodobrasil/featws-api/services"
	"github.com/gin-gonic/gin"
)

// Rulesheets ...
type Rulesheets interface {
	CreateRulesheet() gin.HandlerFunc
	GetRulesheets() gin.HandlerFunc
	GetRulesheet() gin.HandlerFunc
	UpdateRulesheet() gin.HandlerFunc
	DeleteRulesheet() gin.HandlerFunc
}

type rulesheets struct {
	service services.Rulesheets
}

// NewRulesheets ...
func NewRulesheets(service services.Rulesheets) Rulesheets {
	return &rulesheets{
		service: service,
	}
}

// CreateRulesheet 	  	godoc
// @Summary 			Create Rulesheet
// @Description 		Create Rulesheet description
// @Tags 			Rulesheet
// @Accept  		json
// @Produce  			json
// @Param			Rulesheet body payloads.Rulesheet true "Rulesheet body"
// @Success 			200 {object} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 		404 "Not Found"
// @Security 			ApiKeyAuth
// @Router 				/rulesheets [post]
func (rc *rulesheets) CreateRulesheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Pass the context of gin Request
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		var payload payloads.Rulesheet
		defer cancel()

		// validate the request body
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on validate request body: %v", err)

			return
		}

		// use the validator libraty to validate required fields
		if validationErr := validatePayload(&payload); validationErr != nil {
			c.JSON(http.StatusBadRequest, validationErr)
			log.Errorf("Error on validate required fields: %v", validationErr)
			return
		}

		dto, err := dtos.NewRulesheetV1(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on define rulesheet entity: %v", err)
			return
		}

		err = rc.service.Create(ctx, &dto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on create rulesheet: %v", err)
			return
		}

		var response = responses.NewRulesheet(&dto)
		c.JSON(http.StatusCreated, response)
	}
}

// GetRulesheets 		godoc
// @Summary 			List Rulesheets
// @Description 		List Rulesheet description
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				count query boolean false "Total of results"
// @Success 			200 {array} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			ApiKeyAuth
// @Router 				/rulesheets/ [get]
func (rc *rulesheets) GetRulesheets() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		query := c.Request.URL.Query()
		filter := make(map[string]interface{})
		for param, value := range query {
			if len(value) == 1 {
				filter[param] = value[0]
				continue
			}
			filter[param] = value
		}
		_, isCount := query["count"]

		if !isCount {
			dtos, err := rc.service.Find(ctx, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Error{
					Error: err.Error(),
				})
				log.Errorf("Error on fetch more than one rulesheet: %v", err)
				return
			}

			var response = make([]responses.Rulesheet, len(dtos))

			for index, dto := range dtos {
				response[index] = responses.NewRulesheet(dto)
			}

			c.JSON(http.StatusOK, response)
		} else {
			count, err := rc.service.Count(ctx, filter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Error{
					Error: err.Error(),
				})
				log.Errorf("Error on count filterValue: %v", err)
				return
			}

			var response = []responses.Rulesheet{
				{
					FindResult: responses.FindResult{
						Count: count,
					},
				},
			}

			c.JSON(http.StatusOK, response)
		}
	}
}

// GetRulesheet 		godoc
// @Summary 			Get Rulesheet by ID
// @Description 		Get Rulesheet by ID description
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				id path string true "Rulesheet ID"
// @Success 			200 {array} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			ApiKeyAuth
// @Router 				/rulesheets/{id} [get]
func (rc *rulesheets) GetRulesheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Second)
		defer cancel()
		id, exists := c.Params.Get("id")

		if !exists {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: "Required param 'id'",
			})
			log.Error("Error on check if the rulesheet exist")
			return
		}

		entity, err := rc.service.Get(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on fetch unique rulesheet: %v", err)
			return
		}

		if entity != nil {
			var response = responses.NewRulesheet(entity)

			c.JSON(http.StatusOK, response)
			return
		}

		c.String(http.StatusNotFound, "")
	}
}

// UpdateRulesheet 		godoc
// @Summary 			Update Rulesheet by ID
// @Description 		Update Rulesheet by ID description
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				id path string true "Rulesheet ID"
// @Param				rulesheet body payloads.Rulesheet true "Rulesheet body"
// @Success 			200 {array} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			ApiKeyAuth
// @Router 				/rulesheets/{id} [put]
func (rc *rulesheets) UpdateRulesheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		id, exists := c.Params.Get("id")

		if !exists {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: "Required param 'id'",
			})
			log.Error("Error on check if the rulesheet exist")
			return
		}

		_, err := rc.service.Get(ctx, id)
		if err != nil {
			c.String(http.StatusNotFound, "")
			log.Errorf("You are trying to update a non existing record: %v", err)
			return
		}

		var payload payloads.Rulesheet
		// validate the request body
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on binding the payload: %v", err)
		}

		// use the validator libraty to validate required fields
		if validationErr := validatePayload(&payload); validationErr != nil {
			c.JSON(http.StatusBadRequest, validationErr)
			log.Errorf("Error on validate required fields: %v", validationErr)
			return
		}

		iid, _ := strconv.ParseUint(id, 10, 32)
		payload.ID = uint(iid)

		dto, err := dtos.NewRulesheetV1(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on define entity: %v", err)
			return
		}

		updatedEntity, err := rc.service.Update(ctx, dto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on update entity: %v", err)
			return
		}

		if updatedEntity != nil {
			var response = responses.NewRulesheet(updatedEntity)

			c.JSON(http.StatusOK, response)
			return
		}

		c.String(http.StatusNotFound, "")
	}
}

// DeleteRulesheet 		godoc
// @Summary 			Delete Rulesheet by ID
// @Description 		Delete Rulesheet by ID description
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				id path string true "Rulesheet ID"
// @Success 			200 {string} string ""
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			ApiKeyAuth
// @Router 				/rulesheets/{id} [delete]
func (rc *rulesheets) DeleteRulesheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		id, exists := c.Params.Get("id")

		if !exists {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: "Required param 'id'",
			})
			log.Error("Error on check if the rulesheet exist")
			return
		}

		deleted, err := rc.service.Delete(ctx, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Error{
				Error: err.Error(),
			})
			log.Errorf("Error on delete rulesheet: %v", err)
			return
		}

		if !deleted {
			c.String(http.StatusNotFound, "")
			return
		}

		c.String(http.StatusNoContent, "")
	}
}
