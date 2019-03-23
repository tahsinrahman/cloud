package aws

import (
	"github.com/appscode/go/log"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
)

func ParseInstance(in *Ec2Instance) (*v1.MachineType, error) {
	out := &v1.MachineType{
		Spec: v1.MachineTypeSpec{
			SKU:         in.InstanceType,
			Description: in.InstanceType,
			Category:    in.Family,
		},
	}
	cpu, err := in.VCPU.Int64()
	if err != nil {
		log.Warning("ParseInstance failed, intance ", in.InstanceType, ". Reason: ", err)
		cpu = -1
	}
	out.Spec.CPU = resource.NewQuantity(cpu, resource.DecimalExponent)
	ram, err := resource.ParseQuantity(in.Memory.String() + "GiB")
	if err != nil {
		return nil, errors.Errorf("ParseInstance failed, intance %v. Reason: %v.", in.InstanceType, err)
	}
	out.Spec.RAM = &ram
	return out, nil
}

func ParseRegion(in *ec2.Region) *v1.Region {
	return &v1.Region{
		Spec: v1.RegionSpec{
			Region: *in.RegionName,
		},
	}
}
