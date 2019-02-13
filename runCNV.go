package main

import (
	"flag"
	"fmt"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"path"
	"strconv"
	"strings"
)

// os
var (
	//ex, _  = os.Executable()
	//exPath = filepath.Dir(ex)
	pSep = string(os.PathSeparator)
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
		"/share/backup/wangyaoshen/src/ExomeDepth/createScript.sgd.pl",
		"ExomeDepth script to run",
	)
	run2 = flag.String(
		"run2",
		"/share/backup/hanrui/pipeline/CNVkit/bin/run_CNVkit.pl",
		"CNVkit script to run",
	)
	CNVkitControl = flag.String(
		"control",
		"/share/backup/hanrui/pipeline/CNVkit/control/MGISEQ_2000_control/201811/MGISEQ-2000_201811",
		"control of CNVkit",
	)
	submit = flag.Bool(
		"submit",
		false,
		"if auto submit",
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

	runExomeDepth(*run1, *indir, *outdir, *submit)
	runCNVkit(*run2, *indir, *outdir, *CNVkitControl, *submit)
}

func runCNVkit(script, indir, outdir, control string, submit bool) {
	tag := path.Base(indir)
	var args []string
	args = append(args, script)
	args = append(args, strings.Join([]string{outdir, "ExomeDepth", "sample.list.checked"}, pSep))
	args = append(args, control)
	args = append(args, strings.Join([]string{outdir, "CNVkit"}, pSep))
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	args2 = append(args2,
		"-cwd",
		"-l", "vf=31G,p=12",
		"-P", "B2C_SGD",
		"-N", "CNVkit."+tag,
		outdir+"/CNVkit/run.sh",
	)
	if submit {
		fmt.Printf("# qsub %s\n# ", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}

func runExomeDepth(script, indir, outdir string, submit bool) {
	tag := path.Base(indir)
	var args []string
	args = append(args, script, indir)
	args = append(args, strings.Join([]string{outdir, "ExomeDepth"}, pSep))
	args = append(args, tag)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	sampleNum := len(simple_util.File2Array(strings.Join([]string{outdir, "ExomeDepth", "sample.list.checked"}, pSep)))
	if sampleNum > 0 {
		args2 = append(args2,
			"-cwd",
			"-l", "vf="+strconv.Itoa(sampleNum*2)+"G,p="+strconv.Itoa(sampleNum),
			"-P", "B2C_SGD",
			"-N", "ExomeDepth."+tag,
			outdir+"/CNVkit/run.sh",
		)
	} else {
		args2 = append(args2,
			"-cwd",
			"-l", "vf=31G,p=12",
			"-P", "B2C_SGD",
			"-N", "ExomeDepth."+tag,
			outdir+"/CNVkit/run.sh",
		)
	}

	if submit {
		fmt.Printf("# qsub %s\n# ", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}
