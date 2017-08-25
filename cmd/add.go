// Copyright 2017 The Gop Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// CmdAdd represents add a new dependency package and it's dependencies to this project
var CmdAdd = cli.Command{
	Name:        "add",
	Usage:       "add a new dependency",
	Description: `add a new dependency`,
	Action:      runAdd,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "update, u",
			Usage: "update the dependency package",
		},
	},
}

func copyPkg(srcPkgPath, dstPkgPath string) error {
	return CopyDir(srcPkgPath, dstPkgPath, func(path string) bool {
		return strings.HasPrefix(path, filepath.Join(dstPkgPath, ".git")) ||
			strings.HasPrefix(path, filepath.Join(dstPkgPath, "vendor"))
	})
}

// add add one package to vendor
func add(name, projPath, globalGoPath string, isUpdate bool) error {
	if strings.HasPrefix(name, "../") || filepath.IsAbs(name) || strings.HasPrefix(name, "./") {
		return errors.New("relative pkg and absolute pkg is not supported, only packages on GOPATH")
	}

	absPkgPath := filepath.Join(globalGoPath, "src", name)
	dstPath := filepath.Join(projPath, "src", "vendor", name)

	_, err := os.Stat(absPkgPath)
	if err != nil {
		return err
	}

	info, err := os.Stat(dstPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		fmt.Println("Copying", name)
		err = copyPkg(absPkgPath, dstPath)
		if err != nil {
			return err
		}
	} else if isUpdate {
		if !info.IsDir() {
			return fmt.Errorf("Dest dir %s is a file", dstPath)
		}

		fmt.Println("Copying", name)
		os.RemoveAll(dstPath)
		err = copyPkg(absPkgPath, dstPath)
		if err != nil {
			return err
		}
	} else {
		return nil
	}

	imports, err := ListImports(name, absPkgPath, absPkgPath, "", false)
	if err != nil {
		return err
	}

	for i, imp := range imports {
		var has bool
		for j := 0; j < i; j++ {
			if imports[j] == imp {
				has = true
				break
			}
		}
		if has {
			continue
		}

		if err := add(imp, projPath, globalGoPath, isUpdate); err != nil {
			return err
		}
	}
	return nil
}

func runAdd(ctx *cli.Context) error {
	if len(ctx.Args()) <= 0 {
		return errors.New("You have to indicate more than one package")
	}

	globalGoPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		return errors.New("Not found GOPATH")
	}

	names := ctx.Args()

	_, projectRoot, err := analysisDirLevel()
	if err != nil {
		return err
	}

	for _, name := range names {
		if err = add(name, projectRoot, globalGoPath, ctx.IsSet("update")); err != nil {
			return err
		}
	}
	return nil
}
