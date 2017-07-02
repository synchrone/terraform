package command

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/config/module"
	"github.com/mitchellh/cli"
)

// DumpConfigCommand is a Command implementation that dumps the structure of a configuration
type DumpConfigCommand struct {
	Meta
}

func (c *DumpConfigCommand) Run(args []string) int {

	cmdFlags := c.Meta.flagSet("internal-dump config")
	if err := cmdFlags.Parse(args); err != nil {
		return cli.RunResultHelp
	}
	args = cmdFlags.Args()

	var path = "."
	if len(args) > 0 {
		path = args[0]
	}

	var mod *module.Tree
	mod, err := c.Module(path)
	if err == nil && mod == nil {
		fmt.Println(fmt.Sprintf("No terraform files found"))
		return 1
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to load root config module: %s", err))
		return 1
	}

	en := json.NewEncoder(os.Stdout)

	err = en.Encode(struct {
		Resources []*config.Resource
		Variables []*config.Variable
		Outputs   []*config.Output
	}{
		mod.Config().Resources,
		mod.Config().Variables,
		mod.Config().Outputs,
	})

	return 0
}

func (c *DumpConfigCommand) Help() string {
	helpText := `
Usage: terraform analyze config [path]

  Reads and dump a Terraform configuration in json
  form. If no path is specified, the current directory will be used.

Options:

  None

`
	return strings.TrimSpace(helpText)
}

func (c *DumpConfigCommand) Synopsis() string {
	return "Dumps Terraform configuration structure"
}
