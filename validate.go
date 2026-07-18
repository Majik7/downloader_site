package main

import "net/url"

// isValidURL mirrors WTForms' URL() validator: requires an http(s) scheme and a host.
func isValidURL(raw string) bool {
	if raw == "" {
		return false
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return u.Host != ""
}
