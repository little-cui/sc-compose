package command

import (
	"fmt"
	"github.com/little-cui/sc-compose/model"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	upFile string
)

// NewGetCommand returns the cobra command for "up".
func NewUpCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up [options] <sc ip:port>",
		Short: "Start up",
		Run:   upCommandFunc,
	}

	cmd.Flags().StringVarP(&upFile, "file", "f", "sc-compose.yml", "the spec file of sc compose")
	return cmd
}

func upCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		ExitWithError(ErrInvalidParam, "incorrect inputs", args)
	}

	data, err := ioutil.ReadFile(upFile)
	if err != nil {
		ExitWithError(ErrInvalidParam, err.Error())
	}

	var tmpl model.Compose
	err = yaml.Unmarshal(data, &tmpl)
	if err != nil {
		ExitWithError(ErrUnmarshal, err.Error())
	}

	fmt.Println(tmpl)
}
