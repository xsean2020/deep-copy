package main

import (
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"
	"github.com/xsean2020/deep-copy/deepcopy"
	"golang.org/x/tools/go/packages"
)

var gofileReg = regexp.MustCompile(`.*\.go$`)

func main() {
	app := cli.NewApp()
	app.Name = "deep-copy"
	app.Usage = "generator deep-copy code"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "root",
			Usage: "root directory",
			Value: ".",
		},
		cli.BoolFlag{
			Name:  "reverse",
			Usage: "root directory",
		},
		cli.StringSliceFlag{
			Name:  "method",
			Usage: "copy method",
			Value: &cli.StringSlice{"copy", "deepcopy"},
		},
		cli.StringSliceFlag{
			Name:  "tags",
			Usage: "build tags",
		},
	}

	app.Action = func(c *cli.Context) error {
		root := c.String("root")
		methods := c.StringSlice("method")
		reverse := c.Bool("reverse")
		buildTagsF := c.StringSlice("tags")

		// fmt.Println("root:", root, methods, reverse, buildTagsF)

		maxDep := math.MaxInt32
		if !reverse {
			maxDep = 2
		}

		return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				return nil
			}

			if strings.Count(path, string(os.PathListSeparator)) >= maxDep {
				return filepath.SkipAll
			}

			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			if len(files) == 0 {
				return nil
			}

			var patterns []string
			for _, file := range files {
				if !gofileReg.MatchString(file.Name()) {
					continue
				}
				patterns = append(patterns, path+"/"+file.Name())
			}

			packages, err := load(patterns...)
			if err != nil {
				return err
			}

			if len(packages) < 1 {
				return nil
			}
			p := packages[0]
			groups := deepcopy.ParseGeneratorOptions(p)
			for i, group := range groups {
				if len(group) == 0 {
					continue
				}
				out := strings.ReplaceAll(p.CompiledGoFiles[i], ".go", "_gen_deepcopy.go")
				outF, err := os.Create(out)
				if err != nil {
					return err
				}
				defer outF.Close()
				generator := deepcopy.NewGenerator(buildTagsF, methods)
				if err := generator.Generate(outF, p, group, groups); err != nil {
					return err
				}
			}
			return nil
		})
	}

	app.Run(os.Args)

}

func load(patterns ...string) ([]*packages.Package, error) {
	return packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.NeedName | packages.NeedCompiledGoFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedImports,
	}, patterns...)
}
