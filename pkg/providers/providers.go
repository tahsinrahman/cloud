package providers

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/appscode/go/log"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pharmer/cloud/pkg/cmds/options"
	"github.com/pharmer/cloud/pkg/providers/aws"
	"github.com/pharmer/cloud/pkg/providers/azure"
	"github.com/pharmer/cloud/pkg/providers/digitalocean"
	"github.com/pharmer/cloud/pkg/providers/gce"
	"github.com/pharmer/cloud/pkg/providers/linode"
	"github.com/pharmer/cloud/pkg/providers/packet"
	"github.com/pharmer/cloud/pkg/providers/scaleway"
	"github.com/pharmer/cloud/pkg/providers/vultr"
	"github.com/pharmer/cloud/pkg/util"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/version"
)

const (
	Gce          string = "gce"
	DigitalOcean string = "digitalocean"
	Packet       string = "packet"
	Aws          string = "aws"
	Azure        string = "azure"
	Vultr        string = "vultr"
	Linode       string = "linode"
	Scaleway     string = "scaleway"
)

var supportedProvider = []string{
	"gce",
	"digitalocean",
	"packet",
	"aws",
	"azure",
	"vultr",
	"linode",
	"scaleway",
}

type CloudInterface interface {
	GetName() string
	GetCredentials() []v1.CredentialFormat
	GetKubernetes() []v1.KubernetesVersion
	GetRegions() ([]v1.Region, error)
	GetZones() ([]string, error)
	GetMachineTypes() ([]v1.MachineType, error)
}

func NewCloudProvider(opts *options.GenData) (CloudInterface, error) {
	switch opts.Provider {
	case Gce:
		return gce.NewClient(opts.GCEProjectID, opts.CredentialFile)
	case DigitalOcean:
		return digitalocean.NewClient(opts.DoToken)
	case Packet:
		return packet.NewClient(opts.PacketApiKey)
	case Aws:
		return aws.NewClient(opts.AWSRegion, opts.AWSAccessKeyID, opts.AWSSecretAccessKey)
	case Azure:
		return azure.NewClient(opts.AzureTenantId, opts.AzureSubscriptionId, opts.AzureClientId, opts.AzureClientSecret)
	case Vultr:
		return vultr.NewClient(opts.VultrApiToken)
	case Linode:
		return linode.NewClient(opts.LinodeApiToken)
	case Scaleway:
		return scaleway.NewClient(opts.ScalewayToken, opts.ScalewayOrganization)
	}
	return nil, errors.Errorf("Unknown cloud provider: %s", opts.Provider)
}

//get data from api
func GetCloudProvider(cloudInterface CloudInterface) (*v1.CloudProvider, error) {
	var err error
	cloudData := v1.CloudProvider{
		Spec: v1.CloudProviderSpec{
			Name:        cloudInterface.GetName(),
			Credentials: cloudInterface.GetCredentials(),
			Kubernetes:  cloudInterface.GetKubernetes(),
		},
	}
	cloudData.Spec.Regions, err = cloudInterface.GetRegions()
	if err != nil {
		return nil, err
	}
	cloudData.Spec.MachineTypes, err = cloudInterface.GetMachineTypes()
	if err != nil {
		return nil, err
	}
	return &cloudData, nil
}

//write data in [path to pharmer]/data/files/[provider]/
func WriteCloudProvider(cloudData *v1.CloudProvider, fileName string) error {
	cloudData = util.SortCloudProvider(cloudData)
	dataBytes, err := json.MarshalIndent(cloudData, "", "  ")
	if err != nil {
		return err
	}
	dir, err := util.GetWriteDir()
	if err != nil {
		return err
	}
	err = util.WriteFile(filepath.Join(dir, cloudData.Name, fileName), dataBytes)
	return err
}

//region merge rule:
//	if region doesn't exist in old data, but exists in cur data, then add it
//	if region exists in old data, but doesn't exists in cur data, then delete it
//	if region exist in both, then
//		if field data exists in both cur and old data , then take the cur data
//		otherwise, take data from (old or cur)whichever contains it
//
// instanceType merge rule: same as region rule, except
//		if instance exists in old data, but doesn't exists in cur data, then add it , set the deprecated true
//
//In MergeCloudProvider, we merge only the region and instanceType data
func MergeCloudProvider(oldData, curData *v1.CloudProvider) (*v1.CloudProvider, error) {
	//region merge
	regionIndex := map[string]int{} //keep regionName,corresponding region index in oldData.Regions[] as (key,value) pair
	for index, r := range oldData.Spec.Regions {
		regionIndex[r.Spec.Region] = index
	}
	for index := range curData.Spec.Regions {
		pos, found := regionIndex[curData.Spec.Regions[index].Spec.Region]
		if found {
			//location
			if curData.Spec.Regions[index].Spec.Location == "" && oldData.Spec.Regions[pos].Spec.Location != "" {
				curData.Spec.Regions[index].Spec.Location = oldData.Spec.Regions[pos].Spec.Location
			}
			//zones
			if len(curData.Spec.Regions[index].Spec.Zones) == 0 && len(oldData.Spec.Regions[pos].Spec.Zones) != 0 {
				curData.Spec.Regions[index].Spec.Location = oldData.Spec.Regions[pos].Spec.Location
			}
		}
	}

	//instanceType
	instanceIndex := map[string]int{} //keep SKU,corresponding instance index in oldData.MachineTypes[] as (key,value) pair
	for index, ins := range oldData.Spec.MachineTypes {
		instanceIndex[ins.Spec.SKU] = index
	}
	for index := range curData.Spec.MachineTypes {
		pos, found := instanceIndex[curData.Spec.MachineTypes[index].Spec.SKU]
		if found {
			//description
			if curData.Spec.MachineTypes[index].Spec.Description == "" && oldData.Spec.MachineTypes[pos].Spec.Description != "" {
				curData.Spec.MachineTypes[index].Spec.Description = oldData.Spec.MachineTypes[pos].Spec.Description
			}
			//zones
			if len(curData.Spec.MachineTypes[index].Spec.Zones) == 0 && len(oldData.Spec.MachineTypes[pos].Spec.Zones) == 0 {
				curData.Spec.MachineTypes[index].Spec.Zones = oldData.Spec.MachineTypes[pos].Spec.Zones
			}
			//regions
			//if len(curData.Spec.MachineTypes[index].Spec.Regions)==0 && len(oldData.Spec.MachineTypes[pos].Spec.Regions)!=0 {
			//	curData.Spec.MachineTypes[index].Spec.Regions = oldData.Spec.MachineTypes[pos].Spec.Regions
			//}
			//Disk
			if curData.Spec.MachineTypes[index].Spec.Disk == nil && oldData.Spec.MachineTypes[pos].Spec.Disk != nil {
				curData.Spec.MachineTypes[index].Spec.Disk = oldData.Spec.MachineTypes[pos].Spec.Disk
			}
			//RAM
			if curData.Spec.MachineTypes[index].Spec.RAM == nil && oldData.Spec.MachineTypes[pos].Spec.RAM != nil {
				curData.Spec.MachineTypes[index].Spec.RAM = oldData.Spec.MachineTypes[pos].Spec.RAM
			}
			//category
			if curData.Spec.MachineTypes[index].Spec.Category == "" && oldData.Spec.MachineTypes[pos].Spec.Category != "" {
				curData.Spec.MachineTypes[index].Spec.Category = oldData.Spec.MachineTypes[pos].Spec.Category
			}
			//CPU
			if curData.Spec.MachineTypes[index].Spec.CPU == nil && oldData.Spec.MachineTypes[pos].Spec.CPU != nil {
				curData.Spec.MachineTypes[index].Spec.CPU = oldData.Spec.MachineTypes[pos].Spec.CPU
			}
			//to detect it already added to curData
			instanceIndex[curData.Spec.MachineTypes[index].Spec.SKU] = -1
		}
	}
	for _, index := range instanceIndex {
		if index > -1 {
			//using regions as zones
			if len(oldData.Spec.MachineTypes[index].Spec.Regions) > 0 {
				if len(oldData.Spec.MachineTypes[index].Spec.Zones) == 0 {
					oldData.Spec.MachineTypes[index].Spec.Zones = oldData.Spec.MachineTypes[index].Spec.Regions
				}
				oldData.Spec.MachineTypes[index].Spec.Regions = nil
			}
			curData.Spec.MachineTypes = append(curData.Spec.MachineTypes, oldData.Spec.MachineTypes[index])
			curData.Spec.MachineTypes[len(curData.Spec.MachineTypes)-1].Spec.Deprecated = true
		}
	}
	return curData, nil
}

//get data from api , merge it with previous data and write the data
//previous data written in cloud_old.json
func MergeAndWriteCloudProvider(cloudInterface CloudInterface) error {
	log.Infof("Getting cloud data for `%v` provider", cloudInterface.GetName())
	curData, err := GetCloudProvider(cloudInterface)
	if err != nil {
		return err
	}

	oldData, err := util.GetDataFormFile(cloudInterface.GetName())
	if err != nil {
		return err
	}
	log.Info("Merging cloud data...")
	res, err := MergeCloudProvider(oldData, curData)
	if err != nil {
		return err
	}

	//err = WriteCloudProvider(oldData,"cloud_old.json")
	//if err!=nil {
	//	return err
	//}
	log.Info("Writing cloud data...")
	err = WriteCloudProvider(res, "cloud.json")
	if err != nil {
		return err
	}
	return nil
}

//If kubeData.version exists in old data, then
// 		if kubeData.Envs is empty, then delete it,
//      otherwise, replace it
//If kubeData.version doesn't exists in old data, then append it
func MergeKubernetesSupport(data *v1.CloudProvider, kubeData *v1.KubernetesVersion) (*v1.CloudProvider, error) {
	foundIndex := -1
	for index, k := range data.Spec.Kubernetes {
		if version.CompareKubeAwareVersionStrings(k.Spec.Version.String(), kubeData.Spec.Version.String()) == 0 {
			foundIndex = index
		}
	}
	if foundIndex == -1 { //append
		data.Spec.Kubernetes = append(data.Spec.Kubernetes, *kubeData)
	} else { //replace
		data.Spec.Kubernetes[foundIndex] = *kubeData
	}
	return data, nil
}

func AddKubernetesSupport(opts *options.KubernetesData) error {
	kubeData := &v1.KubernetesVersion{}

	kubeData.Spec.Version = &version.Info{GitVersion: opts.Version}

	kubeData.Spec.Envs = map[string]bool{}
	envs := strings.Split(opts.Envs, ",")
	for _, env := range envs {
		if len(env) > 0 {
			kubeData.Spec.Envs[env] = opts.Deprecated
		}
	}
	for _, name := range supportedProvider {
		if opts.Provider != options.AllProvider && opts.Provider != name {
			continue
		}
		log.Infof("Getting cloud data for `%v` provider", name)
		data, err := util.GetDataFormFile(name)
		if err != nil {
			return err
		}
		log.Infof("Adding kubenetes support for `%v` provider", name)
		data, err = MergeKubernetesSupport(data, kubeData)
		if err != nil {
			return err
		}
		log.Infof("Writing cloud data for `%v` provider", name)
		err = WriteCloudProvider(data, "cloud.json")
		if err != nil {
			return err
		}
	}
	return nil
}
