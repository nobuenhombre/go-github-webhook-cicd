package middlewares

import domainapp "go-github-webhook-cicd/src/internal/app/go-github-webhook-cicd/domain"

type HttpMiddleware struct {
	Domain domainapp.IDomainApp
}

func NewHttpMiddleware(dom domainapp.IDomainApp) (mid *HttpMiddleware) {
	mid = new(HttpMiddleware)
	mid.Domain = dom
	return mid
}
