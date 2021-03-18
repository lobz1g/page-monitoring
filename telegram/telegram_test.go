package telegram

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

type receivedMsg struct {
	Url    string
	Body   string
	Method string
}

func TestTelegram_SendMsg(t *testing.T) {
	oldBaseUrl := baseUrl
	defer func() {
		baseUrl = oldBaseUrl
	}()

	var rm receivedMsg
	startServer := func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, _ := ioutil.ReadAll(r.Body)

			rm = receivedMsg{
				Body:   string(body),
				Method: r.Method,
				Url:    r.RequestURI,
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, time.Now().String())
		}))
	}
	server := startServer()
	defer server.Close()

	wg := new(sync.WaitGroup)

	baseUrl = server.URL + "/"
	tests := []struct {
		name    string
		text    string
		channel string
		want    receivedMsg
		wantErr bool
	}{
		{
			name:    "correct",
			wantErr: false,
			text:    "text",
			channel: "channel",
			want: receivedMsg{
				Method: "POST",
				Url:    "/sendMessage",
				Body:   "{\"chat_id\":\"channel\",\"disable_web_page_preview\":true,\"text\":\"text\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := Telegram{channel: tt.channel}
			err := tg.SendMsg(tt.text)
			got := rm
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendMsg() got = %v, want %v", got, tt.want)
			}
			wg.Wait()
		})
	}
}

func TestNew(t *testing.T) {
	oldBaseUrl := baseUrl
	defer func() {
		baseUrl = oldBaseUrl
	}()

	tests := []struct {
		name        string
		token       string
		channel     string
		want        *Telegram
		wantBaseUrl string
	}{
		{
			name:        "correct",
			channel:     "channel",
			token:       "test",
			want:        &Telegram{channel: "@channel"},
			wantBaseUrl: fmt.Sprintf("%s/bot%s/", defaultTelegramUrl, "test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.token, tt.channel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(baseUrl, tt.wantBaseUrl) {
				t.Errorf("baseUrl = %v, want %v", baseUrl, tt.wantBaseUrl)
			}
		})
	}
}
