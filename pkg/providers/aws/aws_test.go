package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"pharmer.dev/cloud/pkg/credential"
)

var opts = credential.AWS{
	CommonSpec: credential.CommonSpec{
		Data: map[string]string{
			credential.AWSRegion: "us-east-1",
		},
	},
}

func TestRegion(t *testing.T) {
	g, err := NewClient(opts)
	if err != nil {
		t.Error(err)
		return
	}
	g.session, err = session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("", "", ""),
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = g.ListRegions()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestInstance(t *testing.T) {
	g, err := NewClient(opts)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = g.ListMachineTypes()
	if err != nil {
		t.Error(err)
		return
	}
}
