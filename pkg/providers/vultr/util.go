package vultr

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pkg/errors"
)

func ParseRegion(in *vultr.Region) *v1.Region {
	return &v1.Region{
		Spec: v1.RegionSpec{
			Location: in.Name,
			Region:   strconv.Itoa(in.ID),
			Zones: []string{
				strconv.Itoa(in.ID),
			},
		},
	}
}

func ParseInstance(in *PlanExtended) (*v1.MachineType, error) {
	out := &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         strconv.Itoa(in.ID),
			Description: in.Name,
			CPU:         resource.NewQuantity(int64(in.VCpus), resource.DecimalExponent),
			Category:    in.Category,
		},
	}
	if in.Deprecated {
		out.Spec.Deprecated = in.Deprecated
	}

	disk, err := strconv.ParseInt(in.Disk, 10, 64)
	if err != nil {
		return nil, errors.Errorf("Parse Instance failed.reasion: %v.", err)
	}
	out.Spec.Disk = resource.NewQuantity(int64(disk), resource.BinarySI)

	ram, err := strconv.ParseInt(in.RAM, 10, 64)
	if err != nil {
		return nil, errors.Errorf("Parse Instance failed.reasion: %v.", err)
	}
	out.Spec.RAM = resource.NewQuantity(int64(ram), resource.BinarySI)

	out.Spec.Zones = []string{}
	for _, r := range in.Regions {
		region := strconv.Itoa(r)
		out.Spec.Zones = append(out.Spec.Zones, region)
	}
	return out, nil
}
