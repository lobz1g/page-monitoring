package monitor

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type (
	Sender interface {
		SendMsg(text string) error
	}

	page struct {
		Page         string
		LastStateErr bool
	}
)

var (
	Tg          Sender
	Timeout     *time.Ticker
	errLoadPage = errors.New("cannot load page content")
)

func NewPage(p string) *page {
	return &page{Page: p}
}

func (p *page) sendFailInfo() {
	text := fmt.Sprintf("Something went wrong in page %s: %s", p.Page, errLoadPage)
	p.send(text)
}

func (p *page) sendStatusCodeInfo(statusCode int) {
	text := fmt.Sprintf("Status code on page %s is: %d", p.Page, statusCode)
	p.send(text)
}

func (p *page) sendChangeInfo() {
	text := fmt.Sprintf("Something was changed on page %s", p.Page)
	p.send(text)
}

func (page) send(text string) {
	err := Tg.SendMsg(text)
	if err != nil {
		log.Error(err)
	}
}

func (p *page) StartMonitoring(wg *sync.WaitGroup, done chan bool) {
	defer wg.Done()
	var oldData string

	for {
		select {
		case <-done:
			return
		case <-Timeout.C:
			log.Info("Start getting info from site")
			r, err := resty.New().R().Get(p.Page)
			if err != nil {
				if !p.LastStateErr {
					p.sendFailInfo()
					p.LastStateErr = true
				}
				log.Error(err)
				continue
			}
			pageData := string(r.Body())
			log.Info("Getting info from site finished")

			log.Info("Start checking info")
			if r.StatusCode() != 200 {
				p.sendStatusCodeInfo(r.StatusCode())
			} else if oldData != pageData {
				if oldData != "" {
					p.sendChangeInfo()
				}
				oldData = pageData
			}
		}
		log.Info("Checking info finished")
	}
}
