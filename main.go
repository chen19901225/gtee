package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	var maxByteStr string
	var BackupCount, isDebug int
	var path string
	app := &cli.App{
		Name:    "gtee",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "max_byte",
				Usage:       "max_byte, eg,max_byte=50m",
				Value:       "50m",
				DefaultText: "50m",
				Destination: &maxByteStr,
			},
			&cli.IntFlag{
				Name:        "backup_count",
				Usage:       "backUpCount",
				Value:       7,
				DefaultText: "7",
				Destination: &BackupCount,
			},
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path to log to",
				Required:    true,
				Destination: &path,
			},
			&cli.IntFlag{
				Name:        "debug",
				Usage:       "debug=0|1",
				Value:       0,
				DefaultText: "0",
				Destination: &isDebug,
			},
		},
		Action: func(c *cli.Context) error {
			return run(
				maxByteStr,  //maxByteStr
				BackupCount, // BackupCount
				path,        // path
				isDebug,     // isDebug
			)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
