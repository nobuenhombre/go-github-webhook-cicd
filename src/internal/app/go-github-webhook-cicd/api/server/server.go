package server

import (
	"errors"
	"fmt"
	configserver "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/config"
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server/router"
	domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"
	"log"
	"net/http"
	"os"
)

type HTTPServer struct {
	Router *router.HTTPRouter
	Server *http.Server
}

// appDomain *domain.AppDomain
func NewHTTPServer(config *configserver.HTTPServerConfig, logFile *os.File, dom domainapp.IDomainApp) (srv *HTTPServer, err error) {
	srv = new(HTTPServer)

	srv.Router = router.NewHTTPRouter(logFile, dom)

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: srv.Router.Router,
	}

	return srv, nil
}

func (srv *HTTPServer) Run() {
	go func() {
		err := srv.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	srv.gracefulShutDown()
}
