package linode

import (
	"fmt"
	"testing"

	"github.com/pharmer/cloud/pkg/credential"
)

func TestRegion(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	rList, err := client.ListRegions()
	if err != nil {
		t.Error(err)
	}
	for _, r := range rList {
		fmt.Println(r.Location)
	}
}

func TestInstance(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	iList, err := client.ListMachineTypes()
	if err != nil {
		t.Error(err)
	}
	for _, i := range iList {
		fmt.Println(i.Spec.Description)
	}
}

func getToken() credential.Linode {
	var v credential.Linode
	v.LoadFromJSON("/home/ac/Downloads/cred/linode.json")
	return v
}
