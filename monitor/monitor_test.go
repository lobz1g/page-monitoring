package monitor

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	telegramMock struct {
		*msg
	}

	msg struct {
		Text string
	}
)

func (t *telegramMock) SendMsg(text string) error {
	t.Text = text
	return nil
}

func Test_page_sendChangeInfo(t *testing.T) {
	oldTg := Tg
	defer func() { Tg = oldTg }()

	tests := []struct {
		name string
		page *page
		want string
	}{
		{
			name: "check msg",
			page: &page{Page: "test"},
			want: "Something was changed on page test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &msg{}
			Tg = &telegramMock{got}
			tt.page.sendChangeInfo()
			if !reflect.DeepEqual(got.Text, tt.want) {
				t.Errorf("sendChangeInfo() got = %v, want %v", got.Text, tt.want)
			}
		})
	}
}

func Test_page_sendStatusCodeInfo(t *testing.T) {
	oldTg := Tg
	defer func() { Tg = oldTg }()

	tests := []struct {
		name       string
		page       *page
		statusCode int
		want       string
	}{
		{
			name:       "check msg",
			page:       &page{Page: "test"},
			statusCode: 200,
			want:       "Status code on page test is: 200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &msg{}
			Tg = &telegramMock{got}
			tt.page.sendStatusCodeInfo(tt.statusCode)
			if !reflect.DeepEqual(got.Text, tt.want) {
				t.Errorf("sendChangeInfo() got = %v, want %v", got.Text, tt.want)
			}
		})
	}
}

func Test_page_sendFailInfo(t *testing.T) {
	oldTg := Tg
	defer func() { Tg = oldTg }()

	tests := []struct {
		name string
		page *page
		want string
	}{
		{
			name: "check msg",
			page: &page{Page: "test"},
			want: "Something went wrong in page test: " + errLoadPage.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &msg{}
			Tg = &telegramMock{got}
			tt.page.sendFailInfo()
			if !reflect.DeepEqual(got.Text, tt.want) {
				t.Errorf("sendChangeInfo() got = %v, want %v", got.Text, tt.want)
			}
		})
	}
}

func Test_page_StartMonitoring1(t *testing.T) {
	oldLogLevel := log.GetLevel()
	oldTimeout := Timeout
	oldTg := Tg
	defer func() {
		Tg = oldTg
		Timeout = oldTimeout
		log.SetLevel(oldLogLevel)
	}()
	log.SetLevel(log.FatalLevel)

	var statusCode int
	startServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
			io.WriteString(w, time.Now().String())
		}))
	}
	server := startServer()
	defer server.Close()

	type args struct {
		wg   *sync.WaitGroup
		done chan bool
	}
	tests := []struct {
		name       string
		page       *page
		args       args
		wantPage   *page
		wantMsg    string
		isErrConn  bool
		statusCode int
	}{
		{
			name:       "status code is not 200",
			page:       &page{Page: server.URL},
			wantPage:   &page{Page: server.URL},
			statusCode: 400,
			wantMsg:    "Status code on page",
		},
		{
			name:       "changed page content",
			page:       &page{Page: server.URL},
			wantPage:   &page{Page: server.URL},
			statusCode: 200,
			wantMsg:    "Something was changed on page",
		},
		{
			name:      "fail while getting page content",
			page:      &page{Page: "localhost"},
			wantPage:  &page{Page: "localhost", LastStateErr: true},
			isErrConn: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeCh := make(chan time.Time)
			defer close(timeCh)

			wg := new(sync.WaitGroup)
			Timeout = &time.Ticker{C: timeCh}
			statusCode = tt.statusCode
			done := make(chan bool)
			gotMsg := &msg{}
			Tg = &telegramMock{gotMsg}

			wg.Add(1)
			go tt.page.StartMonitoring(wg, done)
			timeCh <- time.Now()
			time.Sleep(time.Second)
			timeCh <- time.Now()
			done <- true
			wg.Wait()

			if !reflect.DeepEqual(tt.page, tt.wantPage) {
				t.Errorf("monitor() page = %v, wantPage %v", tt.page, tt.wantPage)
			}
			if !strings.Contains(gotMsg.Text, tt.wantMsg) {
				t.Errorf("Sender = %s does not contain message %s", gotMsg.Text, tt.wantMsg)
			}
		})
	}
}
