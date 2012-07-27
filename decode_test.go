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

func assertString(t *testing.T, name, a, b string) {
	if a != b {
		t.Errorf("%s == %#v != %#v", name, a, b)
	}
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

	if len(list.Servers) != 3 {
		t.Errorf("Server list length is %d, but expected 3.", len(list.Servers))
	}

	assertString(t, "Servers[0].Name", list.Servers[0].Name, "Who")
	assertString(t, "Servers[0].IP", list.Servers[0].IP, "what.invalid")

	assertString(t, "Servers[1].Name", list.Servers[1].Name, "Where")
	assertString(t, "Servers[1].IP", list.Servers[1].IP, "when:12345")

	assertString(t, "Servers[2].Name", list.Servers[2].Name, "☃")
	assertString(t, "Servers[2].IP", list.Servers[2].IP, "snow.man")
}

type EmptyServerList struct {
}

func TestErrMissingField(t *testing.T) {
	f, err := os.Open("testcases/servers.dat")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	var list EmptyServerList

	err = Unmarshal(Uncompressed, f, &list)
	if err == nil {
		t.Error("No error, but one was expected!")
	} else if err.Error() != "nbt: Unhandled TAG_List (0x09)\n\t\tat struct field \"servers\"" {
		t.Error(err)
	}
}

type WronglyTypedServerList struct {
	Servers []WronglyTypedServer `nbt:"servers"`
}

type WronglyTypedServer struct {
	Name string  `nbt:"name"`
	IP   float64 `nbt:"ip"`
}

func TestErrWrongType(t *testing.T) {
	f, err := os.Open("testcases/servers.dat")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	var list WronglyTypedServerList

	err = Unmarshal(Uncompressed, f, &list)
	if err == nil {
		t.Error("No error, but one was expected!")
	} else if err.Error() != "nbt: Tag is TAG_String (0x08), but I don't know how to put that in a float64!\n\t\tat struct field \"ip\"\n\t\tat list index 0\n\t\tat struct field \"servers\"" {
		t.Error(err)
	}
}
