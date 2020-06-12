package cmd

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tomlazar/go-intervals"
)

type CommandCtx struct {
	ctx    context.Context
	client *intervals.Client

	In  io.ReadCloser
	Out io.Writer
	Err io.Writer

	color bool
}

func initCtx() (*CommandCtx, error) {
	burl := viper.GetString("intervals.url")
	apik := viper.GetString("intervals.key")

	if burl == "" {
		return nil, errors.New("intervals url not configured")
	}
	if apik == "" {
		return nil, errors.New("intervals api key not configured")
	}

	c := intervals.Config{
		APIKey:       apik,
		IntervalsURL: burl,
	}

	client := intervals.NewClient(c, http.DefaultClient)

	var out io.Writer = os.Stdout
	var color bool
	if viper.GetBool("color") == false && isTerminal(os.Stdout) {
		out = colorable.NewColorable(os.Stdout)
		color = true
	}

	return &CommandCtx{
		In:    os.Stdin,
		Out:   out,
		Err:   os.Stderr,
		color: color,

		client: client,
		ctx:    context.Background(),
	}, nil
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
