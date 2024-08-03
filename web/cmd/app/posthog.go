//go:build js
// +build js

package main

import (
	"encoding/json"
	"net/url"
	"strings"
	"syscall/js"

	"github.com/ethereum/go-ethereum/log"
)

func isLocalhost() bool {
	window := js.Global()
	location := window.Get("location")
	hostname := location.Get("hostname").String()

	// Check if the hostname is 'localhost' or a localhost IPv4 (127.0.0.1) or IPv6 ([::1]) address
	return hostname == "localhost" || strings.HasPrefix(hostname, "127.") || hostname == "[::1]"
}

func getHref() string {
	window := js.Global()
	location := window.Get("location")
	return location.Get("href").String()
}

type PostHog interface {
	Start()
	Stop()
	Capture(event string, properties map[string]interface{}) error
}

func NewPostHog(urlStr, apiKey, distId string) PostHog {
	if isLocalhost() {
		return &localPostHog{}
	}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return &remotePostHog{
		Url:      parsedUrl,
		ApiKey:   apiKey,
		distId:   distId,
		href:     getHref(),
		postChan: make(chan []byte, 100),
		stopChan: make(chan struct{}),
	}
}

type remotePostHog struct {
	Url      *url.URL
	ApiKey   string
	distId   string
	href     string
	postChan chan []byte
	stopChan chan struct{}
}

func (p *remotePostHog) Start() {
	go func() {
		for {
			select {
			case <-p.stopChan:
				return
			case jsonData := <-p.postChan:
				// Prepare the JavaScript Fetch API options
				opts := js.Global().Get("Object").New()
				opts.Set("method", "POST")
				opts.Set("body", string(jsonData))

				// Properly initialize the headers using JavaScript's Headers object
				headers := js.Global().Get("Headers").New()
				headers.Call("append", "Content-Type", "application/json")
				opts.Set("headers", headers)

				captureUrl := p.Url.ResolveReference(&url.URL{Path: "/capture/"}).String()

				// Make the fetch call
				js.Global().Call("fetch", captureUrl, opts)
			}
		}
	}()
}

func (p *remotePostHog) Stop() {
	p.stopChan <- struct{}{}
}

func (p *remotePostHog) Capture(event string, properties map[string]interface{}) error {
	// log.Debug("Capturing event", "event", event, "properties", properties)
	properties["distinct_id"] = p.distId
	properties["_href"] = p.href
	data := map[string]interface{}{
		"event":      event,
		"properties": properties,
		"api_key":    p.ApiKey,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	p.postChan <- jsonData
	return nil
}

type localPostHog struct{}

func (p *localPostHog) Start() {
	// This is a no-op implementation for the local development environment
}

func (p *localPostHog) Stop() {
	// This is a no-op implementation for the local development environment
}

func (p *localPostHog) Capture(event string, properties map[string]interface{}) error {
	log.Debug("Capturing event", "event", event, "properties", properties)
	return nil
}
