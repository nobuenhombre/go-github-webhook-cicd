package middlewares

import (
	"github.com/gin-gonic/gin"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"net/http"
)

const (
	Project = "project"
)

func (mid *HttpMiddleware) GithubMiddleware(project *configgithub.GitHubProjectConfig) gin.HandlerFunc {
	return func(c *gin.Context) {

		githubService := mid.Domain.GetGithubService()

		err := githubService.ValidatePushEventRequest(c.Request, project)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		c.Set(Project, project)

		c.Next()
	}
}
