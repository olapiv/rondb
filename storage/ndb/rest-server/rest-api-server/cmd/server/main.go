/*
 * This file is part of the RonDB REST API Server
 * Copyright (c) 2022 Hopsworks AB
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"
	"hopsworks.ai/rdrs/internal/config"
	"hopsworks.ai/rdrs/internal/dal/heap"
	"hopsworks.ai/rdrs/internal/log"
	"hopsworks.ai/rdrs/internal/servers"
	"hopsworks.ai/rdrs/version"
)

func main() {
	configFileArg := flag.String("config", "", "Configuration file path")
	versionArg := flag.Bool("version", false, "Print API and application version")
	flag.Parse()

	if *versionArg == true {
		fmt.Printf("App version %s, API version %s\n", version.VERSION, version.API_VERSION)
		return
	}

	err := config.SetFromFileIfExists(*configFileArg)
	if err != nil {
		panic(err)
	}
	conf := config.GetAll()

	logger, err := log.SetupLogger(conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	replace := zap.ReplaceGlobals(logger)
	defer logger.Sync()
	defer replace()

	logger.Sugar().Infof("Current configuration: %s", conf)
	logger.Sugar().Infof("Starting Version : %s, Git Branch: %s (%s). Built on %s at %s",
		version.VERSION, version.BRANCH, version.GITCOMMIT, version.BUILDTIME, version.HOSTNAME)
	logger.Sugar().Infof("Starting API Version : %s", version.API_VERSION)

	runtime.GOMAXPROCS(conf.Internal.GOMAXPROCS)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal)

	newHeap, releaseBuffers, err := heap.New()
	if err != nil {
		panic(err)
	}
	defer releaseBuffers()

	err, cleanupServers := servers.CreateAndStartDefaultServers(logger, newHeap, quit)
	if err != nil {
		panic(err)
	}
	defer cleanupServers()

	/*
	 kill (no param) default send syscall.SIGTERM
	 kill -2 is syscall.SIGINT
	 kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	*/
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Sugar().Infof("Shutting down server...")
}
