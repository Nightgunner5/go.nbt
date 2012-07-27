package nbt

import (
	"os"
	"testing"
)

type ServerList struct {
	Servers []Server `nbt:"servers"`
}

type Server struct {
	Name string `nbt:"name"`
	IP   string `nbt:"ip"`
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
