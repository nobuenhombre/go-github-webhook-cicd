package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"io"
	"net/http"
	"strings"
)

type PushEventRequestHeaders struct {
	Event                      string
	Delivery                   string
	HookID                     string
	HookInstallationTargetID   string
	HookInstallationTargetType string
	Signature                  string
	Signature256               string
}

func NewPushEventRequestHeaders(r *http.Request) *PushEventRequestHeaders {
	headers := &PushEventRequestHeaders{}
	headers.Event = r.Header.Get("X-GitHub-Event")
	headers.Delivery = r.Header.Get("X-GitHub-Delivery")
	headers.HookID = r.Header.Get("X-GitHub-Hook-ID")
	headers.HookInstallationTargetID = r.Header.Get("X-GitHub-Hook-Installation-Target-ID")
	headers.HookInstallationTargetType = r.Header.Get("X-GitHub-Hook-Installation-Target-Type")
	headers.Signature = r.Header.Get("X-Hub-Signature")
	headers.Signature256 = r.Header.Get("X-Hub-Signature-256")

	return headers
}

func (prh *PushEventRequestHeaders) Validate() error {
	if prh.Event != "push" {
		return ge.Pin(ge.New("Event header is not [push]"))
	}

	if prh.Delivery == "" {
		return ge.Pin(ge.New("Delivery header is empty"))
	}

	if prh.HookID == "" {
		return ge.Pin(ge.New("HookID header is empty"))
	}

	if prh.HookInstallationTargetID == "" {
		return ge.Pin(ge.New("HookInstallationTargetID header is empty"))
	}

	if prh.HookInstallationTargetType == "" {
		return ge.Pin(ge.New("HookInstallationTargetType header is empty"))
	}

	if prh.Signature == "" {
		return ge.Pin(ge.New("Signature header is empty"))
	}

	if prh.Signature256 == "" {
		return ge.Pin(ge.New("Signature256 header is empty"))
	}

	return nil
}

type PushEventRequest struct {
	Headers *PushEventRequestHeaders
	Body    []byte
}

func NewPushEventRequest(r *http.Request) (*PushEventRequest, error) {
	request := &PushEventRequest{}

	request.Headers = NewPushEventRequestHeaders(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ge.Pin(err)
	}

	request.Body = body

	return request, nil
}

func (pr *PushEventRequest) Validate(secret string) error {
	err := pr.Headers.Validate()
	if err != nil {
		return ge.Pin(err)
	}

	signatureParts := strings.SplitN(pr.Headers.Signature, "=", 2)
	if len(signatureParts) != 2 {
		return ge.Pin(ge.New("Signature header does not contain (hash type and hash)", ge.Params{"prh.Signature": pr.Headers.Signature}))
	}

	if signatureParts[0] != "sha1" {
		return ge.Pin(ge.New("Signature type should be (sha1)", ge.Params{"signatureType": signatureParts[0]}))
	}

	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(pr.Body)

	hash := fmt.Sprintf("%x", hm.Sum(nil))

	if !hmac.Equal([]byte(hash), []byte(signatureParts[1])) {
		return ge.Pin(ge.New("Signature is invalid"))
	}

	return nil
}
