package distributed

import (
	"github.com/kercylan98/kspace/src/pkg/orm"
	"testing"
	"time"
)

func TestServer_Release(t *testing.T) {
	zookeeper := orm.Zookeeper{}
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
