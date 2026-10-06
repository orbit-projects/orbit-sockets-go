package sockets_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	sockets "github.com/orbit-projects/orbit-sockets-go"
)

func TestMessageStringOmitsPayload(t *testing.T) {
	message, err := sockets.NewMessage(sockets.TextMessage, []byte("private-value"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(message.String(), "private-value") {
		t.Fatal("message string leaked payload")
	}
}

func TestCloseValidation(t *testing.T) {
	for _, item := range []struct {
		code   websocket.StatusCode
		reason string
	}{
		{code: 1005},
		{code: 1016},
		{code: 3000, reason: strings.Repeat("é", 62)},
	} {
		if _, err := sockets.NewClose(item.code, item.reason); err == nil {
			t.Fatalf("expected invalid close value: code=%d reason_bytes=%d", item.code, len(item.reason))
		}
	}
	close, err := sockets.NewClose(1000, strings.Repeat("é", 61))
	if err != nil {
		t.Fatalf("expected 122-byte reason to pass: %v", err)
	}
	if strings.Contains(close.String(), close.Reason()) {
		t.Fatal("close string leaked reason")
	}
}

func TestSettingsRejectUnboundedAndMalformedValues(t *testing.T) {
	settings := sockets.NewSettings()
	settings.MaxConnections = 0
	if _, err := sockets.NewServer(settings, nil); err == nil {
		t.Fatal("expected invalid connection limit")
	}
	settings = sockets.NewSettings()
	settings.AllowedOrigins = []string{"https://example.com/path"}
	if _, err := sockets.NewServer(settings, nil); err == nil {
		t.Fatal("expected origin path to be rejected")
	}
}

func TestHandlerRoundTripAndPathValidation(t *testing.T) {
	server, err := sockets.NewServer(sockets.NewSettings(), func(ctx context.Context, session *sockets.Session) error {
		message, err := session.Receive(ctx)
		if err != nil {
			return err
		}
		return session.Send(ctx, message)
	})
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(server.Handler())
	defer httpServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	address := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, address, nil)
	if err != nil {
		t.Fatal("dial failed")
	}
	defer conn.CloseNow()
	if err := conn.Write(ctx, websocket.MessageText, []byte("echo")); err != nil {
		t.Fatal("write failed")
	}
	typ, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal("read failed")
	}
	if typ != websocket.MessageText || string(data) != "echo" {
		t.Fatalf("unexpected echo: type=%v size=%d", typ, len(data))
	}

	response, err := http.Get(httpServer.URL + "/wrong")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status %d", response.StatusCode)
	}
}
