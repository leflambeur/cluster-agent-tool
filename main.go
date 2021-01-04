package main

import (
	"cluster-agent-tool/cattle"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := &cli.App{

		Name:  "cluster-agent-tool",
		Usage: "Redeploy Rancher Agent Workloads",

		Commands: []cli.Command{
			{
				Name: "get",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "token",
						Usage: "Rancher API Bearer Token",
					},
					cli.StringFlag{
						Name:  "cluster",
						Usage: "Rancher Cluster ID - c-XXXXX",
					},
					cli.StringFlag{
						Name:  "server",
						Usage: "Rancher Server URL - https://<rancher>.com",
					},
					cli.BoolFlag{
						Name:  "apply",
						Usage: "Set if you would like to automatically apply the updated Workload YAML",
					},
				},
				Aliases: []string{"g"},
				Usage:   "Gets a copy of the Workload Deployment YAML from Rancher Server",
				Action: func(c *cli.Context) error {
					cattleDeployment, err := cattle.GetDeploymentSetup()
					if err != nil{
						return err
					}
					fmt.Println(cattleDeployment)
					return nil
				},
			},
			{
				Name: "apply",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:      "deployment",
						TakesFile: true,
						Usage:     "Path to Deployment File",
						Required:  true,
					},
					cli.StringFlag{
						Name:      "config",
						TakesFile: true,
						Usage:     "Path to Kubeconfig file",
					},
				},
				Aliases: []string{"a"},
				Usage:   "Apply a copy of a Rancher Agent Workload Deployment YAML",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name: "status",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:      "config",
						TakesFile: true,
						Usage:     "Path to Kubeconfig file",
					},
					cli.BoolFlag{
						Name:  "verbose",
						Usage: "Verbose",
					},
				},
				Aliases: []string{"s"},
				Usage:   "Get status of current cattle agent pods",
				Action: func(c *cli.Context) error {
					status, err := cattle.PrintStatus(c.Bool("verbose"))
					if err != nil {
						log.Println(err)
					}
					fmt.Println(status)
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
