package util

import (
	"encoding/json"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/version"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/appscode/go/runtime"
	"github.com/pharmer/cloud/pkg/apis/cloud/v1"
	"github.com/pkg/errors"
)

func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return errors.Errorf("failed to create dir `%s`. Reason: %v", dir, err)
	}
	return nil
}

func ReadFile(name string) ([]byte, error) {
	dataBytes, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, errors.Errorf("failed to read `%s`.Reason: %v", name, err)
	}
	return dataBytes, nil
}

func WriteFile(filename string, bytes []byte) error {
	err := ioutil.WriteFile(filename, bytes, 0666)
	if err != nil {
		return errors.Errorf("failed to write `%s`. Reason: %v", filename, err)
	}
	return nil
}

// versions string formate is `1.1.0,1.9.0`
//they are comma separated, no space allowed
func ParseVersions(versions string) []string {
	v := strings.Split(versions, ",")
	return v
}

func MBToGB(in int64) (float64, error) {
	gb, err := strconv.ParseFloat(strconv.FormatFloat(float64(in)/1024, 'f', 2, 64), 64)
	return gb, err
}

func BToGB(in int64) (float64, error) {
	gb, err := strconv.ParseFloat(strconv.FormatFloat(float64(in)/(1024*1024*1024), 'f', 2, 64), 64)
	return gb, err
}

// write directory is [path]/pharmer/data/files
func GetWriteDir() (string, error) {
	dir := filepath.Join(runtime.GOPath(), "src/github.com/pharmer/pharmer/data/files")
	return dir, nil
}

//getting provider data from cloud.json file
//data contained in [path to pharmer]/data/files/[provider]/cloud.json
func GetDataFormFile(provider string) (*v1.CloudProvider, error) {
	data := v1.CloudProvider{}
	dir, err := GetWriteDir()
	if err != nil {
		return nil, err
	}
	dir = filepath.Join(dir, provider, "cloud.json")
	dataBytes, err := ReadFile(dir)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func SortCloudProvider(data *v1.CloudProvider) *v1.CloudProvider {
	sort.Slice(data.Spec.Regions, func(i, j int) bool {
		return data.Spec.Regions[i].Spec.Region < data.Spec.Regions[j].Spec.Region
	})
	for index := range data.Spec.Regions {
		sort.Slice(data.Spec.Regions[index].Spec.Zones, func(i, j int) bool {
			return data.Spec.Regions[index].Spec.Zones[i] < data.Spec.Regions[index].Spec.Zones[j]
		})
	}
	sort.Slice(data.Spec.MachineTypes, func(i, j int) bool {
		return data.Spec.MachineTypes[i].Spec.SKU < data.Spec.MachineTypes[j].Spec.SKU
	})
	for index := range data.Spec.MachineTypes {
		sort.Slice(data.Spec.MachineTypes[index].Spec.Zones, func(i, j int) bool {
			return data.Spec.MachineTypes[index].Spec.Zones[i] < data.Spec.MachineTypes[index].Spec.Zones[j]
		})
	}
	sort.Slice(data.Spec.Kubernetes, func(i, j int) bool {
		return version.CompareKubeAwareVersionStrings(
			data.Spec.Kubernetes[i].Spec.Version.String(), data.Spec.Kubernetes[j].Spec.Version.String()) < 0
	})
	return data
}
