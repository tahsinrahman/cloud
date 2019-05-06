package vultr

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
	regions, err := client.ListRegions()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(regions)
}

func TestInstance(t *testing.T) {
	client, err := NewClient(getToken())
	if err != nil {
		t.Error(err)
	}
	instances, err := client.ListMachineTypes()
	if err != nil {
		t.Error(err)
	}
	for _, i := range instances {
		fmt.Println(i.Spec.SKU)
	}
	fmt.Println("total:", len(instances))
}

func getToken() credential.Vultr {
	var v credential.Vultr
	v.LoadFromJSON("/home/ac/Downloads/cred/vultr.json")
	return v
}
