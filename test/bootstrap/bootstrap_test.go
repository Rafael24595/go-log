package bootstrap_test

import (
	"bytes"
	"testing"

	"github.com/Rafael24595/go-assert/assert"
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
	assert.Inside(t, "pre-configuration", output)
}
