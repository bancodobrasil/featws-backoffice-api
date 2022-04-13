package v1

// ResolveHandler ...
// func ResolveHandler(c *gin.Context) {

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	var input payloads.Resolve
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, responses.Error{
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	resolveContext := dtos.NewResolveV1(input)

// 	err := services.Resolve(ctx, input.Resolver, &resolveContext)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, responses.Error{
// 			Error: err.Error(),
// 		})
// 		return
// 	}

// 	resolverOutput := responses.NewResolve(resolveContext)

// 	c.JSON(http.StatusOK, resolverOutput)
// }
