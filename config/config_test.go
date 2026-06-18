package config_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/kafka-mas/TelegramSender/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testconf struct {
	token string
	socks socks
}
type socks struct {
	address string
	port    int
	user    string
	pass    string
}

func TestAddRecord(t *testing.T) {

	tests := []struct {
		name          string
		confFactory   func(t *testing.T) string
		expectedValue testconf
		wantErr       bool
	}{
		{
			name: "Succes test",
			confFactory: func(t *testing.T) string {
				conf := []byte(`token: 1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890
socks:
  address: 127.0.0.1
  port: 10808
  user: user
  pass: pass
`)
				return setupConfig(t, conf, 0666)
			},
			expectedValue: testconf{
				token: "1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890",
				socks: socks{
					address: "127.0.0.1",
					port:    10808,
					user:    "user",
					pass:    "pass",
				},
			},
			wantErr: false,
		},
		{
			name: "Succes test",
			confFactory: func(t *testing.T) string {
				conf := []byte(`token: 1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890
socks:
  address: 127.0.0.1
  port: 10808
`)
				return setupConfig(t, conf, 0666)
			},
			expectedValue: testconf{
				token: "1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890",
				socks: socks{
					address: "127.0.0.1",
					port:    10808,
					user:    "",
					pass:    "",
				},
			},
			wantErr: false,
		},
		{
			name: "Succes test",
			confFactory: func(t *testing.T) string {
				conf := []byte(`token: 1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890`)
				return setupConfig(t, conf, 0666)
			},
			expectedValue: testconf{
				token: "1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890",
				socks: socks{
					address: "",
					port:    0,
					user:    "",
					pass:    "",
				},
			},
			wantErr: false,
		},
		{
			name: "Error test",
			confFactory: func(t *testing.T) string {
				conf := []byte(`token: 1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890`)
				return setupConfig(t, conf, 0444)
			},
			expectedValue: testconf{
				token: "1234567890:ABCDEFGHJKLMNOPQRSTUVWXYZ1234567890",
				socks: socks{
					address: "",
					port:    0,
					user:    "",
					pass:    "",
				},
			},
			wantErr: false,
		},
		{
			name: "Error test",
			confFactory: func(t *testing.T) string {
				conf := []byte(`phone:123`)
				return setupConfig(t, conf, 0666)
			},
			expectedValue: testconf{
				token: "",
				socks: socks{
					address: "",
					port:    0,
					user:    "",
					pass:    "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confPath := tt.confFactory(t)
			c := config.Yaml(confPath)
			cc, err := c.Read()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				tc := testconf{
					token: cc.Token,
					socks: socks{
						address: cc.Socks.Address,
						port:    cc.Socks.Port,
						user:    cc.Socks.User,
						pass:    cc.Socks.Pass,
					},
				}
				assert.Equal(t, tt.expectedValue, tc)
				assert.NoError(t, err)
			}
		})
	}
}

func setupConfig(t *testing.T, config []byte, perm os.FileMode) string {
	tmpDir := t.TempDir()
	confPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(confPath, config, perm)
	if err != nil {
		log.Println("Error create test config", err)
		return confPath
	}

	require.NoError(t, err)
	return confPath
}
