package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/subscriptions"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2017-12-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/util"
)

type Client struct {
	Data           *AzureData
	SubscriptionId string
	GroupsClient   subscriptions.Client
	VmSizesClient  compute.VirtualMachineSizesClient
}

type AzureData v1.CloudProvider

func NewClient(tenantId, subscriptionId, clientId, clientSecret string) (*Client, error) {
	g := &Client{
		SubscriptionId: subscriptionId,
	}
	var err error
	data, err := util.GetDataFormFile("azure")
	if err != nil {
		return nil, err
	}
	d := AzureData(*data)
	g.Data = &d

	baseURI := azure.PublicCloud.ResourceManagerEndpoint
	config, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantId)
	if err != nil {
		return nil, err
	}

	spt, err := adal.NewServicePrincipalToken(*config, clientId, clientSecret, baseURI)
	if err != nil {
		return nil, err
	}
	g.GroupsClient = subscriptions.NewClient()
	g.GroupsClient.Authorizer = autorest.NewBearerAuthorizer(spt)

	g.VmSizesClient = compute.NewVirtualMachineSizesClient(subscriptionId)
	g.VmSizesClient.Authorizer = autorest.NewBearerAuthorizer(spt)
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
	regionList, err := g.GroupsClient.ListLocations(context.Background(), g.SubscriptionId)
	var regions []v1.Region
	for _, r := range *regionList.Value {
		region := ParseRegion(&r)
		regions = append(regions, *region)
	}
	return regions, err
}

func (g *Client) GetZones() ([]string, error) {
	regions, err := g.GetRegions()
	if err != nil {
		return nil, err
	}
	visZone := map[string]bool{}
	var zones []string
	for _, r := range regions {
		for _, z := range r.Spec.Zones {
			if _, found := visZone[z]; !found {
				zones = append(zones, z)
				visZone[z] = true
			}
		}
	}
	return zones, nil
}

func (g *Client) GetMachineTypes() ([]v1.MachineType, error) {
	zones, err := g.GetZones()
	if err != nil {
		return nil, err
	}
	var instances []v1.MachineType
	//to find the position in instances array
	instancePos := map[string]int{}
	for _, zone := range zones {
		instanceList, err := g.VmSizesClient.List(context.Background(), zone)
		if err != nil {
			return nil, err
		}
		for _, ins := range *instanceList.Value {
			instance, err := ParseInstance(&ins)
			if err != nil {
				return nil, err
			}
			pos, found := instancePos[instance.Spec.SKU]
			if found {
				instances[pos].Spec.Zones = append(instances[pos].Spec.Zones, zone)
			} else {
				instancePos[instance.Spec.SKU] = len(instances)
				instance.Spec.Zones = []string{zone}
				instances = append(instances, *instance)
			}
		}
	}
	return instances, nil
}
