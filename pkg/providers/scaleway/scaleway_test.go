package scaleway

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pharmer/cloud/pkg/util"
)

func TestInstance(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	insList, err := client.ListMachineTypes()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(insList)
}

func getToken() Options {
	b, _ := util.ReadFile("/home/ac/Downloads/cred/scaleway.json")
	var v Options
	fmt.Println(json.Unmarshal(b, &v))
	//fmt.Println(v)
	return v
}
