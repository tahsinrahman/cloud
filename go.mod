module github.com/pharmer/cloud

go 1.12

require (
	cloud.google.com/go v0.39.0 // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.5.0 // indirect
	github.com/Azure/azure-sdk-for-go v29.0.0+incompatible
	github.com/Azure/go-autorest v12.0.0+incompatible
	github.com/JamesClonk/vultr v2.0.0+incompatible
	github.com/appscode/go v0.0.0-20190424183524-60025f1135c9
	github.com/aws/aws-sdk-go v1.19.31
	github.com/creack/goselect v0.0.0-20180501195510-58854f77ee8d // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/digitalocean/godo v1.10.0
	github.com/dnaeon/go-vcr v1.0.1 // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20180920040306-f579f869bbfe // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/juju/ratelimit v1.0.1 // indirect
	github.com/linode/linodego v0.8.0
	github.com/moul/anonuuid v1.1.0 // indirect
	github.com/moul/gotty-client v1.7.0 // indirect
	github.com/onsi/gomega v1.5.0
	github.com/packethost/packngo v0.1.1-0.20190507131943-1343be729ca2
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829 // indirect
	github.com/prometheus/common v0.4.0 // indirect
	github.com/renstrom/fuzzysearch v1.0.2 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/scaleway/scaleway-cli v1.10.2-0.20190329131818-c54911b8b3c5
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/the-redback/go-oneliners v0.0.0-20190417084731-74f7694e6dae
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f // indirect
	golang.org/x/net v0.0.0-20190514140710-3ec191127204
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/sys v0.0.0-20190516014833-cab07311ab81 // indirect
	golang.org/x/tools v0.0.0-20190516015132-d1a3278ee749 // indirect
	google.golang.org/api v0.5.0
	google.golang.org/appengine v1.6.0 // indirect
	google.golang.org/genproto v0.0.0-20190515210553-995ef27e003f // indirect
	gopkg.in/ini.v1 v1.42.0
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb // indirect
	k8s.io/apiextensions-apiserver v0.0.0-20190515024537-2fd0e9006049 // indirect
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22 // indirect
	kmodules.xyz/client-go v0.0.0-20190515205239-a16030cc2e50
	sigs.k8s.io/controller-runtime v0.2.0-alpha.1
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/graymeta/stow => github.com/appscode/stow v0.0.0-20190506085026-ca5baa008ea3
	github.com/renstrom/fuzzysearch => github.com/lithammer/fuzzysearch v1.0.2
	gopkg.in/robfig/cron.v2 => github.com/appscode/cron v0.0.0-20170717094345-ca60c6d796d4
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190311093542-50b561225d70
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
