package config

import "fmt"

type HTTPServerConfig struct {
	host         string
	port         int
	CORS         string
	baseURL      string
	Enpoint      string
	returnFormat string
}

func (h *HTTPServerConfig) GetListenAddr() string {
	return fmt.Sprintf("%s:%d", h.host, h.port)
}

func (h *HTTPServerConfig) GetBaseURL() string {
	if h.Enpoint == "/" {
		switch h.returnFormat {
		case "json":
			return fmt.Sprintf("%s/%%s", h.baseURL)
		}
		return fmt.Sprintf("%s/r/%%s", h.baseURL)
	}
	switch h.returnFormat {
	case "json":
		return fmt.Sprintf("%s%s/%%s", h.baseURL, h.Enpoint)
	}
	return fmt.Sprintf("%s%s/r/%%s", h.baseURL, h.Enpoint)
}
