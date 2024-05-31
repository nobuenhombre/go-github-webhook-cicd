package router

import (
	"github.com/gin-gonic/gin"
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/router/handlers"
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/router/middlewares"
	domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"
	"io"
	"os"
)

type HTTPRouter struct {
	Router      *gin.Engine
	Handlers    *handlers.HttpHandler
	Middlewares *middlewares.HttpMiddleware
}

func NewHTTPRouter(logFile *os.File, dom domainapp.IDomainApp) (router *HTTPRouter) {
	router = new(HTTPRouter)
	router.WriteToLog(logFile)

	router.Handlers = handlers.NewHttpHandler(dom)
	router.Middlewares = middlewares.NewHttpMiddleware(dom)
	router.Router = gin.Default()
	router.Router.Use(router.Middlewares.CORSMiddleware())
	router.AddRoutes()

	return
}

func (r *HTTPRouter) WriteToLog(logFile *os.File) {
	if logFile != nil {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}
}

func (r *HTTPRouter) AddRoutes() {
	r.Router.GET("/", r.Handlers.DefaultHandler)
}
