package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rpsl/sentry-mattermost-sidecar/config"
	mCfg "github.com/rpsl/sentry-mattermost-sidecar/middlewares"
	routers "github.com/rpsl/sentry-mattermost-sidecar/routes"
	log "github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
		FullTimestamp:   true,
	})

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf(err.Error())
	}

	// todo: print motd screen on start
	log.Println(cfg.Debug)

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
		log.Infoln("Web server started in Release mode")
	} else {
		log.Infoln("Web server started in DEBUG mode")
	}

	r := gin.New()
	r.Use(ginlogrus.Logger(log.New()), gin.Recovery())
	r.Use(mCfg.AttachConfig(cfg))

	routers.SetRoutes(r)

	_ = r.Run(fmt.Sprintf(
		"%s:%s",
		cfg.ListenHost,
		cfg.ListenPort,
	))
}
