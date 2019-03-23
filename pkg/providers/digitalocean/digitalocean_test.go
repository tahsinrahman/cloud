package digitalocean

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pharmer/cloud/pkg/util"
)

func TestRegion(t *testing.T) {
	client, err := NewClient(tgetToken())
	if err != nil {
		t.Error(err)
	}
	regions, err := client.GetRegions()
	if err != nil {
		t.Error(err)
	}
	if len(regions) == 0 {
		t.Error("Expected non-empty list of regions")
	}
}

func TestInstance(t *testing.T) {
	client, err := NewClient(tgetToken())
	if err != nil {
		t.Error(err)
	}
	ins, err := client.GetMachineTypes()
	if err != nil {
		t.Error(err)
	}
	if len(ins) == 0 {
		t.Error("Expected non-empty list of intances")
	}
}

func tgetToken() string {
	b, _ := util.ReadFile("/home/ac/Downloads/cred/digitalocean.json")
	v := struct {
		Token string `json:"token"`
	}{}
	fmt.Println(json.Unmarshal(b, &v))
	//fmt.Println(v)
	return v.Token
}
