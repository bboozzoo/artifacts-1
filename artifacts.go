// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mendersoftware/mender-artifact/parser"
	"github.com/mendersoftware/mender-artifact/reader"
	"github.com/mendersoftware/mender-artifact/writer"

	"github.com/urfave/cli"
)

// VERSION of the mender-artifact CLI tool
var VERSION = "0.1"

func writeArtifact(c *cli.Context) error {
	if len(c.String("device-type")) == 0 || len(c.String("artifact-name")) == 0 ||
		len(c.String("update")) == 0 {
		return errors.New("must provide `device-type`, `artifact-name` and `update`")
	}

	he := &parser.HeaderElems{
		Metadata: []byte(""),
	}

	ud := parser.UpdateData{
		P:         &parser.RootfsParser{},
		DataFiles: []string{c.String("update")},
		Type:      "rootfs-image",
		Data:      he,
	}

	name := "mender.tar.gz"
	if len(c.String("output-path")) > 0 {
		name = c.String("output-path")
	}

	devices := c.StringSlice("device-type")

	aw := awriter.NewWriter("mender", c.Int("version"), devices, c.String("artifact-name"))
	return aw.WriteKnown([]parser.UpdateData{ud}, name)
}

func read(aPath string) (*areader.Reader, error) {
	_, err := os.Stat(aPath)
	if err != nil {
		return nil, errors.New("Pathspec '" + aPath +
			"' does not match any files.")
	}

	f, err := os.Open(aPath)
	if err != nil {
		return nil, errors.New("Can not open '" + aPath + "' file.")
	}
	defer f.Close()

	ar := areader.NewReader(f)
	if ar == nil {
		return nil, errors.New("Can not read artifact file.")
	}

	_, err = ar.Read()
	if err != nil {
		return nil, err
	}
	defer ar.Close()

	return ar, nil
}

func readArtifact(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("Nothing specified, nothing read. \nMaybe you wanted" +
			" to say 'artifacts read <pathspec>'?")
	}

	r, err := read(c.Args().First())
	if err != nil {
		return err
	}

	parsers := r.GetWorkers()
	info := r.GetInfo()

	fmt.Printf("Mender artifact:\n")
	fmt.Printf("  Name: %s\n", r.GetArtifactName())
	fmt.Printf("  Format: %s\n", info.Format)
	fmt.Printf("  Version: %d\n", info.Version)
	fmt.Printf("  Compatible devices: '%s'\n", r.GetCompatibleDevices())

	fmt.Printf("\nUpdates:\n")

	for k, p := range parsers {
		fmt.Printf("  %s\n", k)
		fmt.Printf("  Type: '%s'\n", p.GetUpdateType().Type)
		for _, f := range p.GetUpdateFiles() {
			fmt.Printf("  Files:\n")
			fmt.Printf("    %s\n", f.Name)
			fmt.Printf("    size: %d\n", f.Size)
			fmt.Printf("    modified: %s\n", f.Date)
		}
	}
	return nil
}

func validateArtifact(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("Nothing specified, nothing validated. \nMaybe you wanted" +
			" to say 'artifacts validate <pathspec>'?")
	}

	_, err := read(c.Args().First())
	if err != nil {
		return err
	}

	fmt.Println("Artifact file '" + c.Args().First() + "' validated successfully")
	return nil
}

func run() error {
	app := cli.NewApp()
	app.Name = "mender-artifact"
	app.Usage = "Mender artifact read/writer"
	app.UsageText = "mender-artifact [--version][--help] <command> [<args>]"
	app.Version = VERSION

	app.Author = "mender.io"
	app.Email = "contact@mender.io"

	//
	// write
	//
	writeRootfs := cli.Command{
		Name:   "rootfs-image",
		Action: writeArtifact,
	}

	writeRootfs.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "update, u",
			Usage: "Update `FILE`.",
		},
		cli.StringSliceFlag{
			Name: "device-type, t",
			Usage: "Type of device(s) supported by the update. You can specify multiple " +
				"compatible devices providing this parameter multiple times.",
		},
		cli.StringFlag{
			Name:  "artifact-name, n",
			Usage: "Name of the artifact",
		},
		cli.StringFlag{
			Name:  "output-path, o",
			Usage: "Full path to output artifact file.",
		},
		cli.IntFlag{
			Name:  "version, v",
			Usage: "Version of the artifact.",
			Value: 1,
		},
	}

	write := cli.Command{
		Name:  "write",
		Usage: "Writes artifact file.",
		Subcommands: []cli.Command{
			writeRootfs,
		},
	}

	//
	// validate
	//
	validate := cli.Command{
		Name:        "validate",
		Usage:       "Validates artifact file.",
		Action:      validateArtifact,
		UsageText:   "mender-artifact validate [options] <pathspec>",
		Description: "This command validates artifact file provided by pathspec.",
	}

	//
	// read
	//
	read := cli.Command{
		Name:        "read",
		Usage:       "Reads artifact file.",
		Action:      readArtifact,
		UsageText:   "mender-artifact read [options] <pathspec>",
		Description: "This command validates artifact file provided by pathspec.",
	}

	app.Commands = []cli.Command{
		write,
		read,
		validate,
	}
	return app.Run(os.Args)
}

func main() {
	run()
}
