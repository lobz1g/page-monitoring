package config

import (
	"errors"
	"reflect"
	"testing"
)

const correctData = `
{
  "debug": true,
  "token": "token",
  "channel": "1",
  "timeout": "2m",
  "url": [
    "test"
  ]
}
`

func TestGet(t *testing.T) {
	oldReadFile := readFile
	oldUnmarshal := unmarshal
	defer func() {
		readFile = oldReadFile
		unmarshal = oldUnmarshal
	}()

	tests := []struct {
		name      string
		want      *cfg
		wantErr   bool
		readFile  func(name string) ([]byte, error)
		unmarshal func(data []byte, v interface{}) error
	}{
		{
			name:    "error read file",
			want:    nil,
			wantErr: true,
			readFile: func(string) ([]byte, error) {
				return nil, errors.New("fake error")
			},
			unmarshal: oldUnmarshal,
		},
		{
			name:    "error unmarshal",
			want:    nil,
			wantErr: true,
			readFile: func(string) ([]byte, error) {
				return []byte{}, nil
			},
			unmarshal: func([]byte, interface{}) error {
				return errors.New("fake error")
			},
		},
		{
			name: "correct",
			want: &cfg{
				Channel:   "1",
				Timeout:   "2m",
				Urls:      []string{"test"},
				Token:     "token",
				DebugMode: true,
			},
			wantErr: false,
			readFile: func(string) ([]byte, error) {
				return []byte(correctData), nil
			},
			unmarshal: oldUnmarshal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readFile = tt.readFile
			unmarshal = tt.unmarshal

			got, err := Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
