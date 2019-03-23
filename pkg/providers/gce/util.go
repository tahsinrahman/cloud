package gce

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"

	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
)

const (
	CategoryUnknown string = "unknown"
)

func ParseRegion(region *compute.Region) (*v1.Region, error) {
	r := &v1.Region{
		Spec: v1.RegionSpec{
			Region: region.Name,
		},
	}
	r.Spec.Zones = []string{}
	for _, url := range region.Zones {
		zone, err := ParseZoneFromUrl(url)
		if err != nil {
			return nil, err
		}
		r.Spec.Zones = append(r.Spec.Zones, zone)
	}

	return r, nil
}

func ParseZoneFromUrl(url string) (string, error) {
	words := strings.Split(url, "/")
	if len(words) == 0 {
		return "", errors.Errorf("Invaild url: unable to parse zone from url")
	}
	return words[len(words)-1], nil
}

func ParseMachine(machine *compute.MachineType) (*v1.MachineType, error) {
	return &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         machine.Name,
			Description: machine.Description,
			CPU:         resource.NewQuantity(machine.GuestCpus, resource.DecimalExponent),
			RAM:         resource.NewQuantity(machine.MemoryMb, resource.BinarySI),
			Disk:        resource.NewQuantity(machine.MaximumPersistentDisksSizeGb, resource.BinarySI),
			//Category:    ParseCategoryFromSKU(machine.Name),
		},
	}, nil
}

//gce SKU format: [something]-category-[somethin/empty]
func ParseCategoryFromSKU(sku string) string {
	words := strings.Split(sku, "-")
	if len(words) < 2 {
		return CategoryUnknown
	} else {
		return words[1]
	}
}
