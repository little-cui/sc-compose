package usage

import (
	"bytes"
	"fmt"
	"github.com/little-cui/sc-compose/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
)

var (
	commandUsageTemplate *template.Template
	templFuncs           = template.FuncMap{
		"descToLines": func(s string) []string {
			return strings.Split(strings.Trim(s, "\n\t "), "\n")
		},
		"cmdName": func(cmd *cobra.Command, startCmd *cobra.Command) string {
			parts := []string{cmd.Name()}
			for cmd.HasParent() && cmd.Parent().Name() != startCmd.Name() {
				cmd = cmd.Parent()
				parts = append([]string{cmd.Name()}, parts...)
			}
			return strings.Join(parts, " ")
		},
	}
)

func init() {
	commandUsage := `
{{ $cmd := .Cmd }}\
{{ $cmdname := cmdName .Cmd .Cmd.Root }}\
Name:
{{ if not .Cmd.HasParent }}\
{{printf "\t%s - %s" .Cmd.Name .Cmd.Short}}
{{else}}\
{{printf "\t%s - %s" $cmdname .Cmd.Short}}
{{end}}\

Usage:
{{printf "\t%s" .Cmd.UseLine}}
{{ if not .Cmd.HasParent }}\

Version:
{{printf "\t%s" .Version}}
{{end}}\
{{if .Cmd.HasSubCommands}}\


Commands:
{{range .SubCommands}}\
{{ $cmdname := cmdName . $cmd }}\
{{ if .Runnable }}\
{{printf "\t%s\t%s" $cmdname .Short}}
{{end}}\
{{end}}\
{{end}}\
{{ if .Cmd.Long }}\

Descriptions:
{{range $line := descToLines .Cmd.Long}}{{printf "\t%s" $line}}
{{end}}\
{{end}}\
{{if .Cmd.HasLocalFlags}}\

Options:
{{.LocalFlags}}\
{{end}}\
{{if .Cmd.HasInheritedFlags}}\

Global Options:
{{.GlobalFlags}}\
{{end}}
`[1:]

	commandUsageTemplate = template.Must(template.New("command_usage").Funcs(templFuncs).Parse(strings.Replace(commandUsage, "\\\n", "", -1)))
}

func getSubCommands(cmd *cobra.Command) []*cobra.Command {
	var subCommands []*cobra.Command
	for _, subCmd := range cmd.Commands() {
		subCommands = append(subCommands, subCmd)
		subCommands = append(subCommands, getSubCommands(subCmd)...)
	}
	return subCommands
}

func UsageFunc(cmd *cobra.Command) error {
	subCommands := getSubCommands(cmd)
	tabOut := getTabOutWithWriter(os.Stdout)
	commandUsageTemplate.Execute(tabOut, struct {
		Cmd         *cobra.Command
		LocalFlags  string
		GlobalFlags string
		SubCommands []*cobra.Command
		Version     string
	}{
		cmd,
		flagUsages(cmd.LocalFlags()),
		flagUsages(cmd.InheritedFlags()),
		subCommands,
		version.VERSION,
	})
	tabOut.Flush()
	return nil
}

func flagUsages(flagSet *pflag.FlagSet) string {
	x := new(bytes.Buffer)

	flagSet.VisitAll(func(flag *pflag.Flag) {
		if len(flag.Deprecated) > 0 {
			return
		}
		format := ""
		if len(flag.Shorthand) > 0 {
			format = "  -%s, --%s"
		} else {
			format = "   %s   --%s"
		}
		if len(flag.NoOptDefVal) > 0 {
			format = format + "["
		}
		if flag.Value.Type() == "string" {
			// put quotes on the value
			format = format + "=%q"
		} else {
			format = format + "=%s"
		}
		if len(flag.NoOptDefVal) > 0 {
			format = format + "]"
		}
		format = format + "\t%s\n"
		shorthand := flag.Shorthand
		fmt.Fprintf(x, format, shorthand, flag.Name, flag.DefValue, flag.Usage)
	})

	return x.String()
}

func getTabOutWithWriter(writer io.Writer) *tabwriter.Writer {
	aTabOut := new(tabwriter.Writer)
	aTabOut.Init(writer, 0, 8, 1, '\t', 0)
	return aTabOut
}
