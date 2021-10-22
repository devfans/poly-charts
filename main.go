/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os"

	"github.com/polynetwork/bridge-common/log"
	"github.com/urfave/cli/v2"
)

var DESC = ""

func main() {
	app := &cli.App{
		Name:   "charts",
		Usage:  "Generate poly charts",
		Action: start,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Value: "data.txt",
				Usage: "Data input",
			},
			&cli.StringFlag{
				Name:  "path",
				Value: "data",
				Usage: "Data output",
			},
		},
		Before: Init,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("Start error", "err", err)
	}
}

func start(c *cli.Context) error {
	err := Show(c.String("file"), c.String("path"))
	if err != nil {
		log.Error("Failed to show data chart", "err", err)
	}
	return err
}

func Init(ctx *cli.Context) (err error) {
	log.Init()
	DESC = " - " + ctx.String("path")
	return
}

func Show(source, path string) (err error) {
	data, err := feedFile(source)
	if err != nil {
		return
	}
	for c, ps := range data {
		log.Info("Drawing for chain", "chain", c.Name())
		err = Draw(ps, c.Name()+"-"+path+".png")
		if err != nil {
			log.Warn("Draw chart failed", "err", err)
		}
	}
	return
}
