package bootstrap_test

import (
	"bytes"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-log/log/provider/stream"
)

func TestBootstrap_Flush(t *testing.T) {
	log.Message("pre-configuration message")
	log.Message("another initial message")

	var buf bytes.Buffer
	p := stream.StreamProvider{
		Writer: &buf,
	}

	err := log.DefaultFromProvider(t.Context(), p)
	assert.Nil(t, err)

	err = log.OnClose()
	assert.Nil(t, err)

	output := buf.String()
	assert.Contains(t, output, "pre-configuration")
}
