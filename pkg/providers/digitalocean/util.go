package digitalocean

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
)

func ParseRegion(region *godo.Region) *v1.Region {
	return &v1.Region{
		Spec: v1.RegionSpec{
			Location: region.Name,
			Region:   region.Slug,
			Zones: []string{
				region.Slug,
			},
		},
	}
}

func ParseSizes(size *godo.Size) (*v1.MachineType, error) {
	return &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         size.Slug,
			Description: size.Slug,
			CPU:         resource.NewQuantity(int64(size.Vcpus), resource.DecimalExponent),
			RAM:         resource.NewQuantity(int64(size.Memory), resource.BinarySI),
			Disk:        resource.NewScaledQuantity(int64(size.Disk), 3),
			//Category:    ParseCategoryFromSlug(size.Slug),
			Zones: size.Regions,
		},
	}, nil
}

func ParseCategoryFromSlug(slug string) string {
	if strings.HasPrefix(slug, "m-") {
		return "High Memory"
	} else if strings.HasPrefix(slug, "c-") {
		return "High Cpu"
	} else {
		return "General Purpose"
	}
}
