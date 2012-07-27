package nbt

import (
	"testing"
	"os"
)

type ServerList struct {
	
}

func TestServerList(t *testing.T) {
	f, err := os.Open("testcases/servers.dat")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

        var list ServerList

	err = Unmarshal(Uncompressed, f, &list)
	if err != nil {
		t.Error(err)
	}
	t.Error(list)
}
