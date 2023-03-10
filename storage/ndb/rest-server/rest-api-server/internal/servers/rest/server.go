package rest

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"hopsworks.ai/rdrs/internal/config"
	"hopsworks.ai/rdrs/internal/dal/heap"
	"hopsworks.ai/rdrs/internal/handlers/batchpkread"
	"hopsworks.ai/rdrs/internal/handlers/pkread"
	"hopsworks.ai/rdrs/internal/handlers/stat"
)

type RonDBRestServer struct {
	log    *zap.Logger
	server *http.Server
}

func New(log *zap.Logger, host string, port uint16, tlsConfig *tls.Config, heap *heap.Heap) *RonDBRestServer {
	restApiAddress := fmt.Sprintf("%s:%d", host, port)
	log.Sugar().Infof("Initialising REST API server with network address: '%s'", restApiAddress)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New() // gin.Default() for better logging
	registerHandlers(log, router, heap)
	return &RonDBRestServer{
		log: log,
		server: &http.Server{
			Addr:      restApiAddress,
			Handler:   router,
			TLSConfig: tlsConfig,
		},
	}
}

func (s *RonDBRestServer) Start(quit chan os.Signal) (cleanupFunc func()) {
	go func() {
		var err error
		conf := config.GetAll()
		if conf.Security.EnableTLS {
			err = s.server.ListenAndServeTLS(
				conf.Security.CertificateFile,
				conf.Security.PrivateKeyFile,
			)
		} else {
			err = s.server.ListenAndServe()
		}
		if errors.Is(err, http.ErrServerClosed) {
			s.log.Info("REST server closed")
		} else if err != nil {
			s.log.Sugar().Errorf("REST server failed; error: %v", err)
			quit <- syscall.SIGINT
		}
	}()
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := s.server.Shutdown(ctx)
		if err != nil {
			s.log.Sugar().Errorf("failed shutting down REST API server; error: %v", err)
		}
	}
}

type RouteHandler struct {
	log                *zap.Logger
	statsHandler       stat.Handler
	pkReadHandler      pkread.Handler
	batchPkReadHandler batchpkread.Handler
}

func registerHandlers(log *zap.Logger, router *gin.Engine, heap *heap.Heap) {
	router.Use(ErrorHandler(log))

	versionGroup := router.Group(config.VERSION_GROUP)

	routeHandler := &RouteHandler{
		log:                log,
		statsHandler:       stat.New(heap),
		pkReadHandler:      pkread.New(heap),
		batchPkReadHandler: batchpkread.New(heap),
	}

	// ping
	versionGroup.GET("/"+config.PING_OPERATION, routeHandler.Ping)

	// stat
	versionGroup.GET("/"+config.STAT_OPERATION, routeHandler.Stat)

	// batch
	versionGroup.POST("/"+config.BATCH_OPERATION, routeHandler.BatchPkRead)

	// pk read
	tableSpecificGroup := versionGroup.Group(config.DB_TABLE_PP)
	tableSpecificGroup.POST(config.PK_DB_OPERATION, routeHandler.PkRead)
}

// Inspired from https://stackoverflow.com/a/69948929/9068781
func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for i, ginErr := range c.Errors {
			log.Sugar().Errorf("GIN error nr %d: %s", i, ginErr.Error())
		}

		if len(c.Errors) > 0 {
			// Just get the last error to the client
			// status -1 doesn't overwrite existing status code
			c.JSON(-1, c.Errors[len(c.Errors)-1].Error())
		}
	}
}
