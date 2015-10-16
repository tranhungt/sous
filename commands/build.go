package commands

import (
	"github.com/opentable/sous/core"
	"github.com/opentable/sous/tools/cli"
	"github.com/opentable/sous/tools/git"
)

func BuildHelp() string {
	return `sous build detects your project type, and tries to find a matching
OpenTable supported stack to build against. Right now it only supports NodeJS
projects. It builds a docker image, tagged and labelled correctly.

sous build does not have any options yet`
}

func Build(sous *core.Sous, args []string) {
	target := "build"
	if len(args) != 0 {
		target = args[0]
	}
	core.RequireGit()
	core.RequireDocker()
	if err := git.AssertCleanWorkingTree(); err != nil {
		cli.Logf("WARNING: Dirty working tree: %s", err)
	}

	feature, context, appInfo := sous.AssembleFeatureContext(target)
	if !sous.BuildIfNecessary(feature, context, appInfo) {
		cli.Successf("Already built: %s", context.DockerTag())
	}
	name := context.CanonicalPackageName()
	cli.Successf("Successfully built %s v%s as %s", name, appInfo.Version, context.DockerTag())
}
