package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sid-sun/ioctl-api/config"
	"github.com/sid-sun/ioctl-api/src/service"
	"github.com/sid-sun/ioctl-api/src/utils"
	"github.com/sid-sun/ioctl-api/src/view"
	"github.com/sid-sun/ioctl-api/src/view/http/router"
)

type webView struct {
	httpServer *http.Server
}

func (w *webView) Serve() {
	utils.Logger.Info(fmt.Sprintf("[Views] [Web] [Start] Going to listening on http://%s", w.httpServer.Addr))
	go func() {
		if err := w.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal(fmt.Sprintf("[Views] [Web] [Start] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()
}

func (w *webView) Shutdown(ctx context.Context) {
	if err := w.httpServer.Shutdown(ctx); err != nil {
		utils.Logger.Error(fmt.Sprintf("[Views] [Web] [Shutdown]: %v", err))
		panic(err)
	}
}

func Init(svc service.Service, cfg *config.HTTPServerConfig) view.View {
	rtr := router.New(svc, cfg)
	srv := &http.Server{
		Addr:    cfg.GetListenAddr(),
		Handler: rtr,
	}

	return &webView{
		httpServer: srv,
	}
}
