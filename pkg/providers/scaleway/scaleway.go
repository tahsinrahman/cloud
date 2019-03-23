package scaleway

import (
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/util"
	scaleway "github.com/scaleway/scaleway-cli/pkg/api"
)

type Client struct {
	Data      *ScalewayData
	ParClient *scaleway.ScalewayAPI
	AmsClient *scaleway.ScalewayAPI
}

type ScalewayData v1.CloudProvider

func NewClient(scalewayToken, organization string) (*Client, error) {
	g := &Client{}
	var err error
	g.ParClient, err = scaleway.NewScalewayAPI(organization, scalewayToken, "gen-data", "par1")
	if err != nil {
		return nil, err
	}
	g.AmsClient, err = scaleway.NewScalewayAPI(organization, scalewayToken, "gen-data", "ams1")
	if err != nil {
		return nil, err
	}
	data, err := util.GetDataFormFile("scaleway")
	if err != nil {
		return nil, err
	}
	d := ScalewayData(*data)
	g.Data = &d
	return g, nil
}

func (g *Client) GetName() string {
	return g.Data.Name
}

func (g *Client) GetCredentials() []v1.CredentialFormat {
	return g.Data.Spec.Credentials
}

func (g *Client) GetKubernetes() []v1.KubernetesVersion {
	return g.Data.Spec.Kubernetes
}

func (g *Client) GetRegions() ([]v1.Region, error) {
	regions := []v1.Region{
		{
			Spec: v1.RegionSpec{
				Location: "Paris, France",
				Region:   "par1",
				Zones:    []string{"par1"},
			},
		},
		{
			Spec: v1.RegionSpec{
				Location: "Amsterdam, Netherlands",
				Region:   "ams1",
				Zones:    []string{"ams1"},
			},
		},
	}
	return regions, nil
}

func (g *Client) GetZones() ([]string, error) {
	zones := []string{
		"ams1",
		"par1",
	}
	return zones, nil
}

func (g *Client) GetMachineTypes() ([]v1.MachineType, error) {
	instanceList, err := g.ParClient.GetProductsServers()
	if err != nil {
		return nil, err
	}
	var instances []v1.MachineType
	instancePos := map[string]int{}
	for pos, ins := range instanceList.Servers {
		instance, err := ParseInstance(pos, &ins)
		instance.Spec.Zones = []string{"par1"}
		if err != nil {
			return nil, err
		}
		instances = append(instances, *instance)
		instancePos[instance.Spec.SKU] = len(instances) - 1
	}

	instanceList, err = g.AmsClient.GetProductsServers()
	if err != nil {
		return nil, err
	}
	for pos, ins := range instanceList.Servers {
		instance, err := ParseInstance(pos, &ins)
		if err != nil {
			return nil, err
		}
		if index, found := instancePos[instance.Spec.SKU]; found {
			instances[index].Spec.Zones = append(instances[index].Spec.Zones, "ams1")
		} else {
			instance.Spec.Zones = []string{"ams1"}
			instances = append(instances, *instance)
			instancePos[instance.Spec.SKU] = len(instances) - 1
		}
	}

	return instances, nil
}
