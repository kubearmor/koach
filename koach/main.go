package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

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
	// get arguments
	gRPCPortPtr := flag.String("gRPCPort", "32767", "gRPC port")
	flag.Parse()

	// create koach server
	koachServer := server.NewKoachServer(*gRPCPortPtr)
	if koachServer == nil {
		kg.Warnf("Failed to create a koach server (:%s)", *gRPCPortPtr)
		return
	}
	kg.Printf("Created a koach server (:%s)", *gRPCPortPtr)

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
