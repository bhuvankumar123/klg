package main

import (
	"bytes"
	"fmt"
	liblog "log"
	"os"
	"strings"

	app "github.com/bhuvankumar123/klg"
	"github.com/bhuvankumar123/klg/cmd/ldflags"
	"github.com/bhuvankumar123/klg/crud"
	"github.com/bhuvankumar123/klg/proxy"
	"github.com/pkg/errors"
	"github.com/unbxd/go-base/utils/log"
	"github.com/urfave/cli/v2"
)

func badge(cx *cli.Context) (err error) {
	fmt.Println(" ----------------------------------")
	fmt.Println("|          Go-Starter!!!           |")
	fmt.Println(" ----------------------------------")
	fmt.Println("> Module: 		", ldflags.Module)
	fmt.Println("> GitHash: 		", ldflags.GitHash)
	fmt.Println("> GitBranch: 		", ldflags.GitBranch)
	fmt.Println("> GitStat: 		", ldflags.GitStat)
	fmt.Println("> BuildDate: 		", ldflags.BuildDate)
	fmt.Println("> Version: 		", ldflags.Version)
	fmt.Println(" ----------------------------------")
	return
}

func errorstack(errorstr string) string {
	parts := strings.Split(errorstr, ": ")

	var buff bytes.Buffer

	for ix, p := range parts {
		buff.WriteRune('\n')

		for i := 0; i <= ix; i++ {
			buff.WriteRune(' ')
		}

		buff.WriteString("> ")
		buff.WriteString(p)
		if ix > 3 {
			break
		}
	}

	for i := 4; i < len(parts); i++ {
		buff.WriteString(parts[i])
		buff.WriteString(": ")
	}

	return buff.String()
}

// Command Start
func beforeStart(cx *cli.Context) (ax *app.App, err error) {
	logger, err := log.NewZapLogger(
		log.ZapWithLevel(cx.String("log.level")),
		log.ZapWithEncoding(cx.String("log.encoding")),
		log.ZapWithOutput([]string{cx.String("log.output")}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initliase logging")
	}

	pb, err := proxy.NewProxyBinder(
		logger,
		cx.String("proxy.downstream"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy binder")
	}

	// Create log binder
	mb, err := crud.NewHTTPBinder(
		logger,
		cx.String("mongo.uri"),
		cx.String("mongo.database"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create log binder")
	}

	ax, err = app.NewApp(
		app.WithCustomLogger(logger),
		app.WithHTTPTransport(
			cx.String("http.host"),
			cx.String("http.port"),
			cx.StringSlice("http.monitor"),
		),
		app.WithHTTPBinder(pb),
		app.WithHTTPBinder(mb),
	)
	return
}

func actionStart(cx *cli.Context, ax *app.App) (err error) {
	return ax.Open(cx.Context)
}

// main function
func main() {
	var ax *app.App

	err := (&cli.App{
		Name:    "Go-Starter!",
		Usage:   "Go-Starter is a template for all future Go Development at Unbxd.",
		Version: ldflags.Version,
		Before:  badge,
		Flags:   flags(),
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "starts the overpass server",
				Before: func(cx *cli.Context) (err error) {
					ax, err = beforeStart(cx)
					return err
				},
				Action: func(cx *cli.Context) (err error) {
					return actionStart(cx, ax)
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "prints version details of the binary",
			},
		},
		Action: func(cx *cli.Context) (err error) {
			fmt.Println("Go Starter is a template for all future Go Development at Unbxd")
			return
		},
	}).Run(os.Args)

	if err != nil {
		fmt.Println("Something Went Wrong. Failed to start Overpass.: " + err.Error())
		liblog.Fatalf(
			"-- \nfailed to startoverpass. \n--\n Caused By:\n%s\n--",
			errorstack(err.Error()),
		)
	}
}
