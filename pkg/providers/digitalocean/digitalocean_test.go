package digitalocean

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pharmer/cloud/pkg/util"
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

func getToken() Options {
	b, _ := util.ReadFile("/home/ac/Downloads/cred/digitalocean.json")
	var v Options
	fmt.Println(json.Unmarshal(b, &v))
	return v
}
