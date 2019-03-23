package linode

import (
	"github.com/linode/linodego"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ParseRegion(in *linodego.Region) *v1.Region {
	return &v1.Region{
		Spec: v1.RegionSpec{
			Location: in.Country,
			Region:   in.ID,
			Zones: []string{
				in.ID,
			},
		},
	}
}

func ParseInstance(in *linodego.LinodeType) (*v1.MachineType, error) {
	return &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         in.ID,
			Description: in.Label,
			CPU:         resource.NewQuantity(int64(in.VCPUs), resource.DecimalExponent),
			RAM:         resource.NewQuantity(int64(in.Memory), resource.BinarySI),
			Disk:        resource.NewQuantity(int64(in.Disk), resource.DecimalExponent),
		},
	}, nil
}
