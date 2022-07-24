package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/model"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/storageprovider"
	"github.com/fitant/xbin-api/src/utils"
	"github.com/fitant/xbin-api/src/view"
	"github.com/fitant/xbin-api/src/view/http"
)

func main() {
	config.Load()
	utils.InitLogger(config.Cfg)


	sp := storageprovider.InitS3StorageProvider()
	
	sc := model.NewMongoSnippetController(sp)
	svc := service.NewSnippetService(sc, config.Cfg.Svc)

	// Initialise and start serving webview
	httpView := http.Init(svc, &config.Cfg.Http)
	httpView.Serve()

	gracefulShutdown([]view.View{httpView})
}

func gracefulShutdown(views []view.View) {
	// Listen for interrrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Once interrupted - shutdown all views
	utils.Logger.Info("[Main] [gracefulShutdown]: Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, view := range views {
		go view.Shutdown(ctx)
	}
}
