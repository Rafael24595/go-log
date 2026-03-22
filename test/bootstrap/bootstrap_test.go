package bootstrap_test

import (
	"bytes"
	"strings"
	"testing"

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
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err = log.OnClose()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "pre-configuration") {
		t.Error("Bootstrap logs were lost during Init")
	}
}
