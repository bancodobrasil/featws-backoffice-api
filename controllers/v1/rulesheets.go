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
// @Summary 			Criação de Folha de Regra
// @Description         Nessa operação cria uma folha de regra no banco de dados do FeatWS. Para realizar a criação é necessário completar a folha de regra, com no mínimo:
// @Description  		- **nome** no parâmetro *name*;
// @Description   		- **slug** no parâmetro *slug*;
// @Description  		- **descrição** no parâmetro *description*.
// @Description
// @Description  		```
// @Description  		{
// @Description  			"description": "teste no Swagger da API do FeatWS",
// @Description  			"name": "teste Swagger API",
// @Description  			"slug": "teste_Swagger_API"
// @Description  		}
// @Description  		```
// @Description 		Ambos esses parâmetros devem ser uma string, ou seja, deve estar entre "aspas". Não é possível ter uma folha de regra com o mesmo nome de outra.
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				Rulesheet body payloads.Rulesheet true "Rulesheet body"
// @Success 			200 {object} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			Authentication Api Key
// @Security 			Authentication Bearer Token
// @Router 				/rulesheets [post]
func (rc *rulesheets) CreateRulesheet() gin.HandlerFunc {

	return func(c *gin.Context) {

		// Pass the context of gin Request
		ctx, cancel := context.WithTimeout(c.Request.Context(), 1000*time.Second)
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
		//id := c.Query("id")
		c.JSON(http.StatusCreated, response)
	}

}

// GetRulesheets 		godoc
// @Summary 			Listar as Folhas de Regra
// @Description			É possível listar as folhas de regra de algumas maneiras como veremos a seguir:
// @Description
// @Description			- **Sem nenhum parâmetro:** Ao realizar a chamada do endpoint sem a passagem de parâmetros, todas as folhas de regra existentes serão retornadas, contendo informações como nome, ID, e caso estejam disponíveis, descrição e slug.
// @Description			- **Usando o *count*:** Ao habilitar o *count* para *True* será retornado do endpoint o número de Folhas de Regras existentes.
// @Description			- **Usando o *limit*:** Ao utilizar o parâmetro *limit* deve-se especificar o número máximo de respostas desejadas que serão retornadas pela array.
// @Description			- **Usando o *page*:** Ao utilizar o parâmetro *page*, serão retornadas as folhas de regra correspondentes a essa página, onde as folhas são ordenadas em ordem crescente pelo seu ID.
// @Tags 				Rulesheet
// @Accept  			json
// @Produce  			json
// @Param				count query boolean false "Total of results"
// @Param				limit query integer false "Max length of the array returned"
// @Param				page query integer false "Page number that is multiplied by 'limit' to calculate the offset"
// @Success 			200 {array} payloads.Rulesheet
// @Header 				200 {string} Authorization "token access"
// @Failure 			400 {object} responses.Error "Bad Format"
// @Failure 			500 {object} responses.Error "Internal Server Error"
// @Failure 			default {object} responses.Error
// @Response 			404 "Not Found"
// @Security 			Authentication Api Key
// @Security 			Authentication Bearer Token
// @Router 				/rulesheets/ [get]
func (rc *rulesheets) GetRulesheets() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		query := c.Request.URL.Query()
		filter := make(map[string]interface{})
		// TODO: Implement filters correctly

		for param, value := range query {
			if len(value) == 1 {
				filter[param] = value[0]
				continue
			}
			filter[param] = value
		}

		opts := &services.FindOptions{}

		limit, ok := query["limit"]
		if ok {
			limitInt, err := strconv.Atoi(limit[0])
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Error{
					Error: err.Error(),
				})
				log.Errorf("Error on fetch more than one rulesheet: %v", err)
				return
			}
			opts.Limit = limitInt
		}

		page, ok := query["page"]
		if ok {
			pageInt, err := strconv.Atoi(page[0])
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Error{
					Error: err.Error(),
				})
				log.Errorf("Error on fetch more than one rulesheet: %v", err)
				return
			}
			opts.Page = pageInt
		}

		_, isCount := query["count"]

		if !isCount {
			dtos, err := rc.service.Find(ctx, filter, opts)
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
// @Security 			Authentication Api Key
// @Security 			Authentication Bearer Token
// @Router 				/rulesheets/{id} [get]
func (rc *rulesheets) GetRulesheet() gin.HandlerFunc {

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
// @Security 			Authentication Api Key
// @Security 			Authentication Bearer Token
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

		foudedEntity, err := rc.service.Get(ctx, id)
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

		if payload.Slug != "" {
			c.JSON(http.StatusBadRequest, responses.Error{
				Error: "You can't update a slug already defined",
			})
			log.Errorf("You can't update a slug already defined: %v", err)
			return
		}
		payload.Slug = foudedEntity.Slug

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
// @Security 			Authentication Api Key
// @Security 			Authentication Bearer Token
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
