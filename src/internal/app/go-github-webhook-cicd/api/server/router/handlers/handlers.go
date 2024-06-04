package handlers

import (
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/router/middlewares"
	domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	Domain domainapp.IDomainApp
}

func NewHttpHandler(dom domainapp.IDomainApp) (handler *HttpHandler) {
	handler = new(HttpHandler)
	handler.Domain = dom
	return handler
}

func (h *HttpHandler) DefaultHandler(c *gin.Context) {
	c.String(
		http.StatusOK,
		"Welcome API Server",
	)
}

func (h *HttpHandler) GithubWebHookHandler(c *gin.Context) {
	project := c.MustGet(middlewares.Project).(*configgithub.GitHubProjectConfig)

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"project": project,
	})
}
