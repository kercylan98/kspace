package distributed

import (
	"testing"
	"time"
)

func TestServer_Release(t *testing.T) {
	zookeeper := Zookeeper{}
	server := Server{zookeeper.InitUse("127.0.0.1:2181")}
	if err := server.Release(Node{
		Name:    "test-server",
		Address: "localhost:1991",
	}); err != nil {
		t.Fatal(err)
	}

	select {
	case <-time.After(3 * time.Second):
		server.Close()
		return
	}
}
