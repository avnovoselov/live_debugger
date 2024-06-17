package configuration_test

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/avnovoselov/live_debugger/internal/configuration"
)

func TestParser_Parse(t *testing.T) {
	serverIP := "127.0.0.1"
	serverPort := "8080"
	serverInLocation := "/in"
	serverOutLocation := "/out"
	queueSize := 1000
	loggerLevel := "DEBUG"

	server := fmt.Sprintf(`
ip="%s"
port="%s"
in_location="%s"
out_location="%s"
`, serverIP, serverPort, serverInLocation, serverOutLocation)
	queue := fmt.Sprintf(`size=%d`, queueSize)
	logger := fmt.Sprintf(`level="%s"`, loggerLevel)

	path := "path/to/file.toml"
	data := []byte(fmt.Sprintf(`
[server]
%s
[queue]
%s
[logger]
%s`, server, queue, logger))

	fs := fstest.MapFS{
		path: &fstest.MapFile{Data: data},
	}

	p := configuration.NewParser(fs)
	cfg, err := p.Parse(path)
	require.NoError(t, err)
	fmt.Println(cfg)
	require.Equal(t, serverIP, cfg.Server.IP)
	require.Equal(t, serverPort, cfg.Server.Port)
}
