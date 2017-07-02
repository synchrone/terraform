package command

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mitchellh/cli"
)

// DumpPlanCommand is a Command implementation that dumps the structure of a configuration
type DumpPlanCommand struct {
	Meta
}

func (c *DumpPlanCommand) Run(args []string) int {

	cmdFlags := c.Meta.flagSet("internal-dump plan")
	if err := cmdFlags.Parse(args); err != nil {
		return cli.RunResultHelp
	}
	args = cmdFlags.Args()

	path := "terraform.tfplan"
	if len(args) > 0 {
		path = args[0]
	}

	f, err := os.Open(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading file: %s", err))
		return 1
	}
	defer f.Close()

	plan, err := terraform.ReadPlan(f)
	if err != nil {
		if _, err := f.Seek(0, 0); err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading file: %s", err))
			return 1
		}

		plan = nil
	}

	en := json.NewEncoder(os.Stdout)
	err = en.Encode(plan)

	return 0
}

func (c *DumpPlanCommand) Help() string {
	helpText := `
Usage: terraform internal-dump plan [path]

  Reads and dump a Terraform plan in json
  form. If no path is specified, terraform.tfplan will be used.

Options:

  None

`
	return strings.TrimSpace(helpText)
}

func (c *DumpPlanCommand) Synopsis() string {
	return "Dumps Terraform plan structure"
}
