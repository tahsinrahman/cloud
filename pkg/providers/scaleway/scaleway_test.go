package scaleway

import (
	"fmt"
	"testing"

	"github.com/pharmer/cloud/pkg/credential"
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

func getToken() credential.Scaleway {
	var v credential.Scaleway
	v.LoadFromJSON("/home/ac/Downloads/cred/scaleway.json")
	return v
}
