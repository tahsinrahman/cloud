package linode

import (
	"context"
	"net/http"

	"github.com/linode/linodego"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/util"
	"golang.org/x/oauth2"
)

type Client struct {
	Data   *LinodeData
	Client *linodego.Client
}

type LinodeData v1.CloudProvider

func NewClient(linodeApiToken string) (*Client, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: linodeApiToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)
	g := &Client{
		Client: &client,
	}
	var err error
	data, err := util.GetDataFormFile("linode")
	if err != nil {
		return nil, err
	}
	d := LinodeData(*data)
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

//DataCenter as region
func (g *Client) GetRegions() ([]v1.Region, error) {
	regionList, err := g.Client.ListRegions(context.Background(), &linodego.ListOptions{})
	if err != nil {
		return nil, err
	}
	var regions []v1.Region
	for _, r := range regionList {
		region := ParseRegion(&r)
		regions = append(regions, *region)
	}
	return regions, nil
}

//data.Region.Region as Zone
func (g *Client) GetZones() ([]string, error) {
	regionList, err := g.GetRegions()
	if err != nil {
		return nil, err
	}
	var zones []string
	for _, r := range regionList {
		zones = append(zones, r.Spec.Region)
	}
	return zones, nil
}

func (g *Client) GetMachineTypes() ([]v1.MachineType, error) {
	instanceList, err := g.Client.ListTypes(context.Background(), &linodego.ListOptions{})
	if err != nil {
		return nil, err
	}
	var instances []v1.MachineType
	for _, ins := range instanceList {
		instance, err := ParseInstance(&ins)
		if err != nil {
			return nil, err
		}
		instances = append(instances, *instance)
	}
	return instances, nil
}
