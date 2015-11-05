package commands

import (
	"flag"
	"os"

	"github.com/onsi/say"
)

func NewRunCommand() say.Command {
	var name string
	var command string
	var scriptPath string

	flags := flag.NewFlagSet("run", flag.ContinueOnError)

	flags.StringVar(&name, "name", "", "classroom name")
	flags.StringVar(&command, "c", "", "command to run, parsable by the remote shell")
	flags.StringVar(&scriptPath, "f", "", "script file to run on local filesystem")

	return say.Command{
		Name:        "run",
		Description: "Run a command on all VMs, in parallel",
		FlagSet:     flags,
		Run: func(args []string) {
			validateRequiredArgument("name", name)
			c := newControllerFromEnv()
			if command != "" {
				err := c.RunOnVMs(name, command, false)
				say.ExitIfError("Failed running commands in classroom", err)
			} else if scriptPath != "" {
				err := c.RunOnVMs(name, scriptPath, true)
				say.ExitIfError("Failed running script in classroom", err)
			} else {
				say.Fprintln(os.Stderr, 0, say.Red("run requires either -c or -f flag"))
				os.Exit(1)
			}
		},
	}
}
