package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubearmor/koach/koach/config"
	"github.com/kubearmor/koach/koach/database"
	kg "github.com/kubearmor/koach/koach/log"
	"github.com/kubearmor/koach/koach/server"
)

// StopChan Channel
var StopChan chan struct{}

// init Function
func init() {
	StopChan = make(chan struct{})
}

// ==================== //
// == Signal Handler == //
// ==================== //

// GetOSSigChannel Function
func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}

// ========== //
// == Main == //
// ========== //

func main() {
	// init config
	config.InitConfig()

	// get arguments
	gRPCPortPtr := flag.String("gRPCPort", "3001", "gRPC port")
	migrateDBPtr := flag.Bool("migrate", false, "migrate database")
	periodicDataDeletionAge := flag.String("periodic-data-deletion-age", "30d", "periodic data deletion age")
	flag.Parse()

	// init database
	err := database.InitDatabase(config.C.Database)
	if err != nil {
		kg.Errf("Failed to intizialize database")
		return
	}

	if *migrateDBPtr {
		err := database.MigrateDatabase()
		if err != nil {
			kg.Errf("Failed to migrate database")
			return
		}
	}

	// create koach server
	koachServer := server.NewKoachServer(*gRPCPortPtr, database.DB)
	if koachServer == nil {
		kg.Warnf("Failed to create a koach server (:%s)", *gRPCPortPtr)
		return
	}
	kg.Printf("Created a koach server (:%s)", *gRPCPortPtr)

	go koachServer.GetFeedsFromRelay(config.C.Relay)

	go koachServer.PeriodicDataDeletion(*periodicDataDeletionAge)

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	close(StopChan)

	// destroy the koach server
	if err := koachServer.DestroyKoachServer(); err != nil {
		kg.Warnf("Failed to destroy the koach server (%s)", err.Error())
		return
	}
	kg.Print("Destroyed the koach server")
}
