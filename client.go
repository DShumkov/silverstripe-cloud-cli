package main

import (
	"flag"
	"os"
)

func check(e error) {
	if e != nil {
		Log.Err(e.Error())
	}
}

func main() {

	var configPath, projectID string
	var debug bool
	var envID, reference, refType string

	deployCmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	deployCmd.BoolVar(&debug, "debug", false, "debug mode ON/OFF")
	deployCmd.StringVar(&configPath, "config", "config.json", "path to your config json file")
	deployCmd.StringVar(&envID, "envID", "", "Environment ID")
	deployCmd.StringVar(&reference, "ref", "", "Reference for deployment "+RedColour+"(required)"+ResetColour)
	deployCmd.StringVar(&refType, "type", "", "Type of reference [branch]")
	deployCmd.StringVar(&projectID, "projectID", "", "ID of your project")

	fetchCmd := flag.NewFlagSet("fetch", flag.ExitOnError)
	fetchCmd.BoolVar(&debug, "debug", false, "debug mode ON/OFF")
	fetchCmd.StringVar(&configPath, "config", "config.json", "path to your config json file")
	fetchCmd.StringVar(&projectID, "projectID", "", "ID of your project")

	if len(os.Args) < 2 {
		Log.Err("expected 'fetch' or 'deploy' commands")
	}
	switch os.Args[1] {
	case "fetch":
		fetchCmd.Parse(os.Args[2:])
	case "deploy":
		deployCmd.Parse(os.Args[2:])
	default:
		Log.Err("expected 'fetch' or 'deploy' commands")
	}

	config := Config{}
	config.ConfigFile = configPath
	config.Init()

	if debug || config.Debug {
		Log.Ok("Debug mode ON")
		Log.DebugOn()
	}

	if "" == projectID {
		projectID = config.ProjectID

		if "" == projectID {
			Log.Err("Bad project ID")
		}
	}

	switch os.Args[1] {
	case "fetch":
		RunFetch(&config, projectID)
	case "deploy":
		if "" == envID {
			envID = config.EnvID
			if "" == envID {
				Log.Err("Bad env ID")
			}
		}

		if "" == reference {
			Log.Err("Bad deployment reference")
		}

		if "" == refType {
			refType = config.ReferenceType
			if "" == refType {
				Log.Err("Bad deployment reference type")
			}
		}
		req := DeploymentRequest{}
		req.Ref = reference
		req.RefType = refType
		req.Bypass = true
		req.BypassAndStart = true
		req.Send(&config, envID, projectID)
	default:
		Log.Err("expected 'fetch' or 'deploy' subcommands")
	}

}
