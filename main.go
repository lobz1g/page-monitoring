package main

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/lobz1g/page-monitoring/config"
	"github.com/lobz1g/page-monitoring/monitor"
	"github.com/lobz1g/page-monitoring/telegram"
)

const defaultTimeout = time.Minute * 30

func main() {
	log.SetFormatter(&log.JSONFormatter{DisableTimestamp: true})
	log.SetReportCaller(true)

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.DebugMode {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	parsedDuration, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		monitor.Timeout = time.NewTicker(defaultTimeout)
	} else {
		monitor.Timeout = time.NewTicker(parsedDuration)
	}
	defer monitor.Timeout.Stop()

	monitor.Tg = telegram.New(cfg.Token, cfg.Channel)

	wg := new(sync.WaitGroup)
	done := make(chan bool)
	defer close(done)
	for _, v := range cfg.Urls {
		p := monitor.NewPage(v)
		wg.Add(1)
		go p.StartMonitoring(wg, done)
	}
	wg.Wait()
}
