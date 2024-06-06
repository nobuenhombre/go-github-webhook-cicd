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
		"Welcome Github webhook API Server",
	)
}

func (h *HttpHandler) GithubWebHookHandler(c *gin.Context) {
	project := c.MustGet(middlewares.Project).(*configgithub.GitHubProjectConfig)

	err := h.Domain.GetQueueService().Push(project)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
