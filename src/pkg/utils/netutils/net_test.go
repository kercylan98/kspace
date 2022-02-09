package netutils

import "testing"

func TestGetOutBoundIP(t *testing.T) {
	t.Log(GetOutBoundIP())
}
