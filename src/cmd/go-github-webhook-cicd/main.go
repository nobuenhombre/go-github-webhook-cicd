package main

import (
	"github.com/getsentry/sentry-go"
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/api/server"
	"go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/cli"
	configapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/config"
	domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"
	"log"
	"os"
	"time"
)

func main() {
	// Reading the values in the command line keys
	cliCfg, err := cli.NewConfig()
	if err != nil {
		log.Fatalf(" -[exit]- cli NewConfig() error [%v]\n", err)
	}

	// We write the log to a file, not to the screen
	var logFile *os.File

	if len(cliCfg.LogFile) > 0 {
		logFile, err = os.OpenFile(cliCfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf(" -[exit]- error OpenFile log file (%v): %v", cliCfg.LogFile, err)
		}

		log.SetOutput(logFile)

		defer func() {
			err = logFile.Close()
			if err != nil {
				log.Fatalf(" -[exit]- error Closing log file (%v): %v", cliCfg.LogFile, err)
			}
		}()
	}

	appCfg := configapp.NewConfig()
	err = appCfg.Load(cliCfg.ConfigFile)
	if err != nil {
		log.Fatalf(" -[exit]- appCfg.Load() error [%v]\n", err)
	}

	//----------------------------
	// Sentry.io - to catch PANICS
	//----------------------------
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              appCfg.Sentry.DSN,
		Environment:      appCfg.Sentry.Environment,
		Debug:            true,
		AttachStacktrace: true,
		Release:          "go-github-webhook-cicd:v0.0.1",
	})
	if err != nil {
		log.Fatalf(" -[exit]- cli sentry.Init() error [%v]\n", err)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(5 * time.Second)
	//---------------------------------

	func() {
		defer sentry.Recover()
		//--------------------------------
		// do all of the scary things here
		//--------------------------------

		sentry.CaptureMessage("App Run")

		dom, err := domainapp.NewAppDomain(appCfg)
		if err != nil {
			log.Fatalf(" -[exit]- domainapp.NewAppDomain() error [%v]\n", err)
		}

		srv, err := server.NewHTTPServer(&appCfg.Hosts.API, logFile, dom)
		if err != nil {
			log.Fatalf(" -[exit]- server.NewHTTPServer() error [%v]\n", err)
		}

		srv.Run() // wait terminal signal
	}()
}
