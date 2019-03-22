package scaleway

import (
	"github.com/pharmer/pharmer/data"
	"github.com/pharmer/pharmer/hack/pharmer-tools/util"
	scaleway "github.com/scaleway/scaleway-cli/pkg/api"
)

func ParseInstance(name string, in *scaleway.ProductServer) (*data.MachineType, error) {
	out := &data.MachineType{
		SKU:         name,
		Description: in.Arch,
		CPU:         int(in.Ncpus),
	}
	var err error
	out.RAM, err = util.BToGB(int64(in.Ram))
	if err != nil {
		return nil, err
	}
	//if in.Baremetal {
	//	out.Category = "BareMetal"
	//} else {
	//	out.Category = "Cloud Servers"
	//}
	return out, nil
}
