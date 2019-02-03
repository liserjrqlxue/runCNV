package main

import (
	"flag"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"path"
)

var (
	indir = flag.String(
		"indir",
		"",
		"wes batch workdir",
	)
	outdir = flag.String(
		"outdir",
		"",
		"output workdir, default is basename of -indir",
	)
	run1 = flag.String(
		"run1",
		"run.ExomeDepth.sh",
		"ExomeDepth script to run",
	)
	run2 = flag.String(
		"run2",
		"run.CNVkit.sh",
		"CNVkit script to run",
	)
)

func main() {
	flag.Parse()
	if *indir == "" {
		flag.Usage()
		os.Exit(0)
	}
	if *outdir == "" {
		*outdir = path.Base(*indir)
	}
	simple_util.RunCmd(*run1, *indir, *outdir)
	simple_util.RunCmd(*run2, *indir, *outdir)
}