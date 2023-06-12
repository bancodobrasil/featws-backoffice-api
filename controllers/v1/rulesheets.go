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

// Rulesheets defines a set of methods for handling CRUD operations on rulesheets in a Gin web framework.
//
// Property:
//   - CreateRulesheet: A function that handles the creation of a new rulesheet. It should receive data from the client and store it in a db or other storage system.
//   - GetRulesheets: is a function that handles the HTTP GET request to retrieve a list of all rulesheets. It returns a gin.HandlerFunc which is a function that handles the request and sends the response.
//   - GetRulesheet: is a function that handles the HTTP GET request to retrieve a specific rulesheet from a database or other data source. It takes in a gin.Context object as a parameter and returns a gin.HandlerFunc that can be used as a middleware to handle the request. The function should typically extract the ID
//   - UpdateRulesheet: is a function that handles the updating of a specific rulesheet. It takes in a gin context object and returns a gin handler function. This handler function should retrieve the updated rulesheet data from the request body, validate it, and update the corresponding rulesheet in the database. The handler
//   - DeleteRulesheet: is a function that handles the deletion of a specific
//
// rulesheet. It is a gin.HandlerFunc, which means it is a function that takes in a gin.Context object
// as its parameter and returns nothing. The function should retrieve the ID of the rulesheet to be
// deleted from the request parameters or
type Rulesheets interface {
	CreateRulesheet() gin.HandlerFunc
	GetRulesheets() gin.HandlerFunc
	GetRulesheet() gin.HandlerFunc
	UpdateRulesheet() gin.HandlerFunc
	DeleteRulesheet() gin.HandlerFunc
}

// The type "rulesheets" contains a service called "services.Rulesheets". The "service" property is a variable of type "services.Rulesheets". It is likely
// that this variable is used to access or manipulate data related to rulesheets in some way within the code.
type rulesheets struct {
	service services.Rulesheets
}

// NewRulesheets creates a new instance of the Rulesheets struct with a given service.
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
// @Description			Para criar uma folha de regra basta clicar em **Try it out** , complete a folha de regra com os dados desejados, em seguida, clique em **Execute**.
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
// The above code is defining a function that creates a new rulesheet. It receives a gin context and
// returns a gin handler function. The function first validates the request body and required fields
// using the validator library. It then creates a new rulesheet DTO using the payload received in the
// request. Finally, it calls the service to create the rulesheet and returns a JSON response with the
// created rulesheet data.
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
// @Description
// @Description			Para listar as folhas de regra basta clicar em **Try it out** , complete com o formado desejados, em seguida, clique em **Execute**.
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
// GetRulesheets returns a `gin.HandlerFunc`, this function handles HTTP requests to retrieve rulesheets
// from a database. It first extracts any query parameters from the request URL and uses them to set options
// for the database query. If the `count` parameter is present, it returns the count of rulesheets that match
// the query. Otherwise, it retrieves the rulesheets from the database and returns them as a JSON response.
// The response is formatted using a `responses.Rulesheet` struct.
func (rc *rulesheets) GetRulesheets() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		query := c.Request.URL.Query()
		filter := make(map[string]interface{})
		// TODO: Implement filters correctly

		// for param, value := range query {
		// 	if len(value) == 1 {
		// 		filter[param] = value[0]
		// 		continue
		// 	}
		// 	filter[param] = value
		// }

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
// @Summary 			Obter Folha de Regra por ID
// @Description 		Para se obter a folha de regra por ID, basta clicar em **Try it out** e colocar o ID desejado em *id*. Em seguida, clique em **Execute** e caso o ID exista retornará a folha de regra com o número de ID desejado.
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
// GetRulesheet is defining a function that returns a Gin middleware function
// that retrieves a rulesheet entity from a service based on the ID parameter passed in the
// request. If the entity is found, it is returned as a JSON response with a 200 status code. If the
// entity is not found, a 404 status code is returned. If there is an error during the retrieval
// process, a 500 status code is returned with an error message. The function uses a context with a
// timeout of 10 seconds to ensure that the request does not take too long to complete.
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
// @Summary 			Atualizar Folha de Regra por ID
// @Description			Para atualizar ou editar uma folha de regra, é necessário enviar o ID da folha desejada no campo *id*, juntamente com os parâmetros da regra no corpo da solicitação no parâmetro *rulesheet*. Para realizar essa atualização clique no botão **Try it out** e preencher os campos com os dados desejados, em seguida, clicar em **Execute** para enviar a solicitação de atualização.
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
// UpdateRulesheet is defining a function that handles the update of a rulesheet entity. It receives a
// request with a JSON payload containing the updated information for the entity, validates the
// payload, and updates the entity in the database using a service. If the update is successful, it
// returns a JSON response with the updated entity information. If the entity is not found, it returns
// a 404 status code.
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
// @Summary 			Deletar Folha de Regra por ID
// @Description 		Para excluir uma folha de regra, é necessário clicar no botão **Try it out** e preencher o campo *id* com o ID da folha de regra que se deseja excluir. Em seguida, clique em **Execute** para enviar a solicitação de exclusão.
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
// DeleteRulesheet is defining a function that handles HTTP requests to delete a rulesheet. It first
// extracts the ID of the rulesheet to be deleted from the request parameters. It then uses a service
// to attempt to delete the rulesheet with the given ID. If the deletion is successful, it returns a
// 204 No Content response. If the rulesheet is not found, it returns a 404 Not Found response. If
// there is an error during the deletion process, it returns a 500 Internal Server Error response with
// an error message.
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
