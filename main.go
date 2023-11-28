package main

import (
	"os"
	"os/signal"
	"syscall"

	C "github.com/matu666/clash-backup/constant"
	"github.com/matu666/clash-backup/hub"
	"github.com/matu666/clash-backup/proxy/http"
	"github.com/matu666/clash-backup/proxy/socks"
	"github.com/matu666/clash-backup/tunnel"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := C.GetConfig()
	if err != nil {
		log.Fatalf("Read config error: %s", err.Error())
	}

	port, socksPort := C.DefalutHTTPPort, C.DefalutSOCKSPort
	section := cfg.Section("General")
	if key, err := section.GetKey("port"); err == nil {
		port = key.Value()
	}

	if key, err := section.GetKey("socks-port"); err == nil {
		socksPort = key.Value()
	}

	err = tunnel.GetInstance().UpdateConfig()
	if err != nil {
		log.Fatalf("Parse config error: %s", err.Error())
	}

	go http.NewHttpProxy(port)
	go socks.NewSocksProxy(socksPort)

	// Hub
	if key, err := section.GetKey("external-controller"); err == nil {
		go hub.NewHub(key.Value())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
