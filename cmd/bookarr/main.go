/*
  Copyright (C) 2017 Sinuhé Téllez Rivera

  dir2opds is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  dir2opds is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with dir2opds.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/signal"
	"syscall"
	"time"

	"bookarr/api/opds1"
	"bookarr/storage/dir"

	"github.com/gin-gonic/gin"
)

var (
	dirRoot = flag.String("dir", "./books", "A directory with books.")
	port    = flag.Int("p", 8080, "port to listen on")
)

func main() {

	flag.Parse()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	router := gin.Default()

	// Create a new instance of the OPDS struct.
	storage := dir.NewFileStore(*dirRoot)
	opdsv1Prefix, _ := url.Parse("/opds/v1")
	s := opds1.New(opdsv1Prefix.String(), storage)

	router.GET(opdsv1Prefix.JoinPath("*path").String(), s.Handler)
	router.GET(opdsv1Prefix.String(), s.Handler)
	router.HEAD("/opds/v1/*all", func(c *gin.Context) {
		c.Header("Content-Type", "application/atom+xml;profile=opds-catalog;kind=navigation")
		c.Status(http.StatusOK)
	})
	router.HEAD("/opds/v1", func(c *gin.Context) {
		c.Header("Content-Type", "application/atom+xml;profile=opds-catalog;kind=navigation")
		c.Status(http.StatusOK)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 30 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
