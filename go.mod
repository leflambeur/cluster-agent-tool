module github.com/leflambeur/cluster-agent-tool

go 1.15

replace k8s.io/client-go v12.0.0+incompatible => k8s.io/client-go v0.19.0

require (
	github.com/AlecAivazis/survey/v2 v2.2.7
	github.com/dustin/go-humanize v1.0.0
	github.com/rancher/log v0.1.2
	github.com/rancher/norman v0.0.0-20200930000340-693d65aaffe3
	github.com/rancher/types v0.0.0-20201223181453-72359190db1b
	github.com/urfave/cli v1.22.5
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v12.0.0+incompatible
)
