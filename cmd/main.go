// Package api re_web_page_analyzer API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/api"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/config"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/log"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
)

func main() {
	cf := config.InitConfig()
	// Init logger
	initLogger(cf.LogLevel)

	// injected client for services
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	staff := service.NewStaffService()
	parser := service.NewParserService(httpClient, cf.WorkerCount)
	r := api.NewHandler(staff, parser)
	srv := &http.Server{Addr: cf.ApiListener, Handler: r}
	log.Infof("Start service on http://%s", cf.ApiListener)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Infof("Service - listen: %s", err)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var once sync.Once
		for range signalChan {
			once.Do(func() {
				log.Warningf("Service %s received a shutdown signal...", cf.ServiceName)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(ctx); err != nil {
					log.Errorln(err)
				}
				log.Infoln("Service - Received an interrupt closing connection...")
				log.Warningf("Service %s stopped successfully", cf.ServiceName)
				cleanupDone <- true
			})
		}
	}()
	log.Warningf("Service %s started successfully", cf.ServiceName)
	<-cleanupDone
}

func initLogger(loggerLevel log.Level) {
	logger, err := log.NewServiceLogger(log.Fields{
		log.FieldApplication: "web_analyzer",
		log.FieldService:     "re_web_page_analyzer",
		log.FieldCategory:    "web",
	}, loggerLevel)

	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(-1)
	}
	log.SetLogger(logger)
	log.Info("Logger was successfully initialized")
}
