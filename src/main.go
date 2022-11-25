package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/sid-sun/ioctl-api/config"
	"github.com/sid-sun/ioctl-api/src/model"
	"github.com/sid-sun/ioctl-api/src/service"
	"github.com/sid-sun/ioctl-api/src/storageprovider"
	"github.com/sid-sun/ioctl-api/src/utils"
	"github.com/sid-sun/ioctl-api/src/view"
	"github.com/sid-sun/ioctl-api/src/view/http"
)

func main() {
	config.Load()
	utils.InitLogger(config.Cfg)

	sp := storageprovider.InitS3StorageProvider()

	sc := model.NewS3SnippetController(sp)
	svc := service.NewSnippetService(sc, config.Cfg.Svc)

	// Initialise and start serving webview
	httpView := http.Init(svc, &config.Cfg.Http)
	httpView.Serve()

	gracefulShutdown([]view.View{httpView})
}

func gracefulShutdown(views []view.View) {
	// Listen for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Once interrupted - shutdown all views
	utils.Logger.Info("[Main] [gracefulShutdown]: Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, v := range views {
		go v.Shutdown(ctx)
	}
}
