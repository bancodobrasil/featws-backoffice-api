package middlewares

type Middleware interface {
	Run()
}

func InitializeMiddlewares() {
	NewVerifyAuthTokenMiddleware()
}
