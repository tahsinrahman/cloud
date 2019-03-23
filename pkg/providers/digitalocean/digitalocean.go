package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/util"
	"golang.org/x/oauth2"
)

type Client struct {
	Data   *DigitalOceanData
	Client *godo.Client
	Ctx    context.Context
}

type DigitalOceanData v1.CloudProvider

func NewClient(doToken string) (*Client, error) {
	g := &Client{
		Ctx:  context.Background(),
		Data: &DigitalOceanData{},
	}
	var err error
	g.Client = getClient(g.Ctx, doToken)

	data, err := util.GetDataFormFile("digitalocean")
	if err != nil {
		return nil, err
	}
	d := DigitalOceanData(*data)
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
	regionList, _, err := g.Client.Regions.List(g.Ctx, &godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	regions := []v1.Region{}
	for _, region := range regionList {
		r := ParseRegion(&region)
		regions = append(regions, *r)
	}
	return regions, nil
}

//Rgion.Slug is used as zone name
func (g *Client) GetZones() ([]string, error) {
	regionList, _, err := g.Client.Regions.List(g.Ctx, &godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	zones := []string{}
	for _, region := range regionList {
		zones = append(zones, region.Slug)
	}
	return zones, nil
}

func (g *Client) GetMachineTypes() ([]v1.MachineType, error) {
	sizeList, _, err := g.Client.Sizes.List(g.Ctx, &godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	instances := []v1.MachineType{}
	for _, s := range sizeList {
		ins, err := ParseSizes(&s)
		if err != nil {
			return nil, err
		}
		instances = append(instances, *ins)
	}
	return instances, nil
}

func getClient(ctx context.Context, doToken string) *godo.Client {
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: doToken,
	}))
	return godo.NewClient(oauthClient)
}
