module github.com/u-cto-devops/lguctl

go 1.17

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/Microsoft/go-winio v0.4.9 // indirect
	github.com/alecthomas/assert v0.0.0-20170929043011-405dbfeb8e38 // indirect
	github.com/alecthomas/colour v0.1.0 // indirect
	github.com/alecthomas/repr v0.0.0-20200325044227-4184120f674c // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/atotto/clipboard v0.1.3
	github.com/aws/aws-sdk-go v1.35.19
	github.com/boombuler/barcode v0.0.0-20161226211916-fe0f26ff6d26 // indirect
	github.com/codahale/hdrhistogram v0.9.0 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20180721085148-1ef1cc838816+incompatible // indirect
	github.com/docker/spdystream v0.0.0-20170912183627-bc6354cbbc29 // indirect
	github.com/fatih/color v1.9.0
	github.com/gokyle/hotp v0.0.0-20160218004637-c180d57d286b // indirect
	github.com/gravitational/configure v0.0.0-20160909185025-1db4b84fe9db // indirect
	github.com/gravitational/form v0.0.0-20151109031454-c4048f792f70 // indirect
	github.com/gravitational/kingpin v2.1.11-0.20190130013101-742f2714c145+incompatible // indirect
	github.com/gravitational/oxy v0.0.0-20200916204440-3eb06d921a1d // indirect
	github.com/gravitational/roundtrip v1.0.0 // indirect
	github.com/gravitational/teleport v4.3.10+incompatible
	github.com/gravitational/teleport/api v0.0.0-20220906233424-dc371d91b66e
	github.com/gravitational/trace v1.1.17
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/keybase/go-keychain v0.0.0-20201121013009-976c83ec27a6
	github.com/kylelemons/godebug v0.0.0-20160406211939-eadb3ce320cb // indirect
	github.com/mailgun/minheap v0.0.0-20131208021033-7c28d80e2ada // indirect
	github.com/mailgun/timetools v0.0.0-20141028012446-7e6055773c51 // indirect
	github.com/mailgun/ttlmap v0.0.0-20150816203249-16b258d86efc // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/mdp/rsc v0.0.0-20160131164516-90f07065088d // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pquerna/otp v0.0.0-20160912161815-54653902c20e // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tstranex/u2f v0.0.0-20160508205855-eb799ce68da4 // indirect
	github.com/vulcand/predicate v1.1.0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20151027082146-e0fe6f683076 // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20150808065054-e02fc20de94c // indirect
	github.com/xeipuuv/gojsonschema v0.0.0-20151204154511-3988ac14d6f6 // indirect
	go.opentelemetry.io/otel/trace v1.9.0
	golang.org/x/crypto v0.0.0-20220126234351-aa10faf2a1f8 // indirect
	gopkg.in/ini.v1 v1.60.2
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.20.0-alpha.1.0.20200922235617-829ed199f4e0 // indirect
	k8s.io/klog/v2 v2.60.1 // indirect
	k8s.io/kubectl v0.19.0
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect

)

replace (
	github.com/coreos/go-oidc => github.com/gravitational/go-oidc v0.0.3
	github.com/iovisor/gobpf => github.com/gravitational/gobpf v0.0.1
	github.com/sirupsen/logrus => github.com/gravitational/logrus v0.10.1-0.20171120195323-8ab1e1b91d5f
)

replace github.com/codahale/hdrhistogram => github.com/HdrHistogram/hdrhistogram-go v0.0.0-20200919145931-8dac23c8dac1
