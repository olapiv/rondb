package servers

import (
	"crypto/tls"
	"fmt"
	"os"

	"go.uber.org/zap"
	"hopsworks.ai/rdrs/internal/config"
	"hopsworks.ai/rdrs/internal/dal"
	"hopsworks.ai/rdrs/internal/dal/heap"

	"hopsworks.ai/rdrs/internal/security/apikey/authcache"
	"hopsworks.ai/rdrs/internal/security/tlsutils"
	"hopsworks.ai/rdrs/internal/servers/grpc"
	"hopsworks.ai/rdrs/internal/servers/rest"
)

func CreateAndStartDefaultServers(log *zap.Logger, heap *heap.Heap, quit chan os.Signal) (err error, cleanup func()) {
	cleanup = func() {}

	// Connect to RonDB
	conf := config.GetAll()
	connectString := config.GenerateMgmdConnectString(conf)
	dalErr := dal.InitRonDBConnection(log, connectString, true)
	if dalErr != nil {
		return fmt.Errorf("failed creating RonDB connection; error: %w", dalErr), cleanup
	}
	cleanupRonDB := func() {
		dalErr = dal.ShutdownConnection(log)
		if dalErr != nil {
			log.Error(dalErr.Error())
		}
	}

	var tlsConfig *tls.Config
	if conf.Security.EnableTLS {
		tlsConfig, err = tlsutils.GenerateTLSConfig(
			conf.Security.RequireAndVerifyClientCert,
			conf.Security.RootCACertFile,
			conf.Security.CertificateFile,
			conf.Security.PrivateKeyFile,
		)
		if err != nil {
			cleanupRonDB()
			return fmt.Errorf("failed generating tls configuration; error: %w", err), cleanup
		}
	}

	grpcServer := grpc.New(log, tlsConfig, heap)
	err, cleanupGrpc := grpc.Start(
		log,
		grpcServer,
		conf.GRPC.ServerIP,
		conf.GRPC.ServerPort,
		quit,
	)
	if err != nil {
		cleanupRonDB()
		return fmt.Errorf("failed starting gRPC server; error: %w", err), cleanup
	}

	restServer := rest.New(
		log,
		conf.REST.ServerIP,
		conf.REST.ServerPort,
		tlsConfig,
		heap,
	)
	cleanupRest := restServer.Start(quit)
	return nil, func() {
		cleanupRonDB()
		cleanupGrpc()
		cleanupRest()

		// Clean API Key Cache
		authcache.Reset()
	}
}
