package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/TylerBrock/colorjson"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api <endpoint>",
	Short: "make an authenticated intervals api request",
	RunE: func(c *cobra.Command, args []string) error {
		ctx, err := initCtx()
		if err != nil {
			return err
		}

		req, err := ctx.client.NewRequest("GET", args[0], nil)
		if err != nil {
			return err
		}

		var obj map[string]interface{}
		_, err = ctx.client.Do(req, &obj)
		if err != nil {
			return err
		}

		if ctx.color {
			// Make a custom formatter with indent set
			f := colorjson.NewFormatter()
			f.Indent = 4

			// Marshall the Colorized JSON
			s, _ := f.Marshal(obj)

			fmt.Fprintln(ctx.Out, string(s))
		} else {
			s, _ := json.Marshal(obj)
			fmt.Fprintln(ctx.Out, string(s))
		}

		return nil
	},
}
