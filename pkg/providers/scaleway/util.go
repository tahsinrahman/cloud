package scaleway

import (
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	scaleway "github.com/scaleway/scaleway-cli/pkg/api"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ParseInstance(name string, in *scaleway.ProductServer) (*v1.MachineType, error) {
	out := &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         name,
			Description: in.Arch,
			CPU:         resource.NewQuantity(int64(in.Ncpus), resource.DecimalExponent),
			RAM:         resource.NewMilliQuantity(int64(in.Ram), resource.BinarySI),
		},
	}
	//if in.Baremetal {
	//	out.Category = "BareMetal"
	//} else {
	//	out.Category = "Cloud Servers"
	//}
	return out, nil
}
