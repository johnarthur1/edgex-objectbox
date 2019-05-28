//
// Copyright (c) 2017
// Mainflux
// Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/objectbox/edgex-objectbox"
	"github.com/objectbox/edgex-objectbox/internal"
	"github.com/objectbox/edgex-objectbox/internal/export/client"
	"github.com/objectbox/edgex-objectbox/internal/pkg/correlation"
	"github.com/objectbox/edgex-objectbox/internal/pkg/startup"
	"github.com/objectbox/edgex-objectbox/internal/pkg/usage"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

func main() {
	start := time.Now()
	var (
		useRegistry bool
		useProfile  string
	)

	flag.BoolVar(&useRegistry, "registry", false, "Indicates the service should use Registry.")
	flag.BoolVar(&useRegistry, "r", false, "Indicates the service should use Registry.")
	flag.StringVar(&useProfile, "profile", "", "Specify a profile other than default.")
	flag.StringVar(&useProfile, "p", "", "Specify a profile other than default.")
	flag.Usage = usage.HelpCallback
	flag.Parse()

	params := startup.BootParams{UseRegistry: useRegistry, UseProfile: useProfile, BootTimeout: internal.BootTimeoutDefault}
	startup.Bootstrap(params, client.Retry, logBeforeInit)

	ok := client.Init(useRegistry)
	if !ok {
		logBeforeInit(fmt.Errorf("%s: Service bootstrap failed", clients.ExportClientServiceKey))
		os.Exit(1)
	}

	client.LoggingClient.Info("Service dependencies resolved...")
	client.LoggingClient.Info(fmt.Sprintf("Starting %s %s ", clients.ExportClientServiceKey, edgex.Version))

	http.TimeoutHandler(nil, time.Millisecond*time.Duration(client.Configuration.Service.Timeout), "Request timed out")
	client.LoggingClient.Info(client.Configuration.Service.StartupMsg)

	errs := make(chan error, 2)
	listenForInterrupt(errs)
	client.StartHTTPServer(errs)

	// Time it took to start service
	client.LoggingClient.Info("Service started in: " + time.Since(start).String())
	client.LoggingClient.Info("Listening on port: " + strconv.Itoa(client.Configuration.Service.Port))
	c := <-errs
	client.Destruct()
	client.LoggingClient.Warn(fmt.Sprintf("terminating: %v", c))

	os.Exit(0)
}

func logBeforeInit(err error) {
	l := logger.NewClient(clients.ExportClientServiceKey, false, "", models.InfoLog)
	l.Error(err.Error())
}

func listenForInterrupt(errChan chan error) {
	go func() {
		correlation.LoggingClient = client.LoggingClient
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}
