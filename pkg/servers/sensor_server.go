package servers

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pingcap/log"
	"go.uber.org/zap"
	"github.com/zhangjinpeng87/openbms/pkg/config"
)



// SensorServer is a TCP server that listens for sensor data.
type SensorServer struct {
	l *net.TCPListener
	c *config.SensorServerConfig
}

// NewSensorServer creates a new SensorServer.
func NewSensorServer(cfg *config.SensorServerConfig) *SensorServer {
	return &SensorServer{c: cfg}
}

// Start starts the server.
func (s *SensorServer) Start() error {
	// Create a TCP listener on localhost:8080.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(s.c.Host),
		Port: s.c.Port,
	})
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.l = l


}

// prepare prepares the server for serve.
func (s *SensorServer) prepare() error {
	// Create a TCP listener on localhost:8080.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(s.c.Host),
		Port: s.c.Port,
	})
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.l = l

	// create router
	r := mux.NewRouter()
	r.HandleFunc("/sensor", s.handleSensor).Methods(http.MethodPost)
}

func (s *SensorServer) handleSensor(w http.ResponseWriter, r *http.Request) {
	var state BatteryState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		log.Errorf("failed to decode sensor: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Debug("received sensor: %+v", sensor)
	w.WriteHeader(http.StatusOK)
}

func (s *SensorServer) updateState(state BatteryState) {
	log.Debug("received sensor: %+v", sensor)
}


	