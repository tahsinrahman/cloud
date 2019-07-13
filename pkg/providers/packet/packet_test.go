package packet

import (
	"testing"

	"pharmer.dev/cloud/pkg/credential"
)

func TestRegion(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	regions, err := client.ListRegions()
	if err != nil {
		t.Error(err)
	}
	if len(regions) == 0 {
		t.Error("Expected non-empty list of regions")
	}
}

func TestInstance(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	ins, err := client.ListMachineTypes()
	if err != nil {
		t.Error(err)
	}
	if len(ins) == 0 {
		t.Error("Expected non-empty list of intances")
	}
}

func getToken() credential.Packet {
	var v credential.Packet
	v.LoadFromJSON("/home/ac/Downloads/cred/packet.json")
	return v
}
