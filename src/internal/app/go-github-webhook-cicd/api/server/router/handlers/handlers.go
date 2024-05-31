package handlers

import (
	domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"
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
