package aws

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/util"
	"github.com/pkg/errors"
)

type Client struct {
	Data    *AwsData
	Session *session.Session
}

type AwsData v1.CloudProvider

type Ec2Instance struct {
	Family       string      `json:"family"`
	InstanceType string      `json:"instance_type"`
	Memory       json.Number `json:"memory"`
	VCPU         json.Number `json:"vCPU"`
	Pricing      interface{} `json:"pricing"`
}

func NewClient(awsRegionName, awsAccessKeyId, awsSecretAccessKey string) (*Client, error) {
	g := &Client{}
	var err error
	g.Session, err = session.NewSession(&aws.Config{
		Region:      &awsRegionName,
		Credentials: credentials.NewStaticCredentials(awsAccessKeyId, awsSecretAccessKey, ""),
	})
	if err != nil {
		return nil, err
	}
	data, err := util.GetDataFormFile("aws")
	if err != nil {
		return nil, err
	}
	d := AwsData(*data)
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
	//Create new EC2 client
	svc := ec2.New(g.Session)
	regionList, err := svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}
	var regions []v1.Region
	for _, r := range regionList.Regions {
		regions = append(regions, *ParseRegion(r))
	}
	tempSession, err := session.NewSession(&aws.Config{
		Credentials: g.Session.Config.Credentials,
	})
	if err != nil {
		return nil, err
	}
	for pos, region := range regions {
		tempSession.Config.Region = &region.Spec.Region
		svc := ec2.New(tempSession)
		zoneList, err := svc.DescribeAvailabilityZones(nil)
		if err != nil {
			return nil, err
		}
		region.Spec.Zones = []string{}
		for _, z := range zoneList.AvailabilityZones {
			if *z.RegionName != region.Spec.Region {
				return nil, errors.Errorf("Wrong available zone for %v.", region.Spec.Region)
			}
			region.Spec.Zones = append(region.Spec.Zones, *z.ZoneName)
		}
		regions[pos].Spec.Zones = region.Spec.Zones
	}
	return regions, nil
}

func (g *Client) GetZones() ([]string, error) {
	visZone := map[string]bool{}
	regionList, err := g.GetRegions()
	if err != nil {
		return nil, err
	}
	var zones []string
	for _, r := range regionList {
		for _, z := range r.Spec.Zones {
			if _, found := visZone[z]; !found {
				visZone[z] = true
				zones = append(zones, z)
			}
		}
	}
	return zones, nil
}

//https://ec2instances.info/instances.json
//https://github.com/powdahound/ec2instances.info
func (g *Client) GetMachineTypes() ([]v1.MachineType, error) {

	client := &http.Client{}
	req, err := getInstanceRequest()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var instanceList []Ec2Instance
	err = json.Unmarshal(body, &instanceList)
	if err != nil {
		return nil, err
	}
	var instances []v1.MachineType
	for _, ins := range instanceList {
		i, err := ParseInstance(&ins)
		if err != nil {
			return nil, err
		}
		instances = append(instances, *i)
	}
	return instances, nil
}

func getInstanceRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://ec2instances.info/instances.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	return req, nil
}
