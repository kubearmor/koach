package server

import (
	"net"
	"time"

	kg "github.com/kubearmor/koach/koach/log"
)

// ================== //
// == Koach Server == //
// ================== //

// KoachServer Structure
type KoachServer struct {
	// port
	Port string

	// gRPC listener
	Listener net.Listener
}

// NewKoachServer Function
func NewKoachServer(port string) *KoachServer {
	ks := &KoachServer{}

	ks.Port = port

	// listen to gRPC port
	listener, err := net.Listen("tcp", ":"+ks.Port)
	if err != nil {
		kg.Errf("Failed to listen a port (%s)\n", ks.Port)
		return nil
	}
	ks.Listener = listener

	return ks
}

// DestroyKoachServer Function
func (ks *KoachServer) DestroyKoachServer() error {
	// wait for a while
	time.Sleep(time.Second * 1)

	// close listener
	if ks.Listener != nil {
		if err := ks.Listener.Close(); err != nil {
			kg.Err(err.Error())
		}
		ks.Listener = nil
	}

	return nil
}
