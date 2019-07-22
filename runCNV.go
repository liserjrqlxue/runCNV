package main

import (
	"flag"
	"fmt"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"path"
	"path/filepath"
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
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/ExomeDepth/createScript.sgd.pl",
		"ExomeDepth script to run",
	)
	run2 = flag.String(
		"run2",
		"/ifs9/BC_PS/hanrui/pipeline/CNVkit/bin/run_CNVkit.pl",
		"CNVkit script to run",
	)
	run3 = flag.String(
		"run3",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/SMA_WES/run_SMN_CNV_v2.pl",
		"SMA script to run",
	)
	CNVkitControl = flag.String(
		"control",
		"/ifs9/BC_PS/hanrui/pipeline/CNVkit/control/MGISEQ_2000_control/201906/MGISEQ-2000_201906",
		"control of CNVkit",
	)
	SMAControl = flag.String(
		"smn",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/ExomeDepth/workspace/SMA_WES/SMA_v2.txt.control_gene.csv",
		"control of SMA",
	)
	SMAGene = flag.String(
		"geneinfo",
		"/ifs9/BC_PS/hanrui/pipeline/SMA_WES/PP100.gene.info.bed",
		"gene info for SMA",
	)
	submit = flag.Bool(
		"submit",
		false,
		"if auto submit",
	)
	maxThread = flag.Int(
		"thread",
		40,
		"max thread limit, for parallel run and calculate memery usage",
	)
	proj = flag.String(
		"proj",
		"B2C_SGD",
		"project of qsub -P",
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

	runExomeDepth(*run1, *indir, *outdir, *submit, *maxThread)
	runCNVkit(*run2, *indir, *outdir, *CNVkitControl, *submit)
	runSMA(*run3, *indir, *outdir, *SMAGene, *SMAControl, *submit)
}

func runCNVkit(script, indir, outdir, control string, submit bool) {
	tag, _ := filepath.Abs(indir)
	tag = path.Base(tag)
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
		"-P", *proj,
		"-N", "CNVkit."+tag,
		strings.Join([]string{outdir, "CNVkit", "run.sh"}, pSep),
	)
	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}

func runExomeDepth(script, indir, outdir string, submit bool, thread int) {
	tag, _ := filepath.Abs(indir)
	tag = path.Base(tag)
	var args []string
	args = append(args, script, strconv.Itoa(thread), indir)
	args = append(args, strings.Join([]string{outdir, "ExomeDepth"}, pSep))
	args = append(args, tag)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	sampleNum := len(simple_util.File2Array(strings.Join([]string{outdir, "ExomeDepth", "sample.list.checked"}, pSep)))
	if sampleNum > 0 {
		if sampleNum > thread {
			sampleNum = thread
		}
		args2 = append(args2,
			"-cwd",
			"-l", "vf="+strconv.Itoa(sampleNum*2)+"G,p="+strconv.Itoa(sampleNum),
			"-P", *proj,
			"-N", "ExomeDepth."+tag,
			strings.Join([]string{outdir, "ExomeDepth", "run.sh"}, pSep),
		)
	} else {
		args2 = append(args2,
			"-cwd",
			"-l", "vf=31G,p=12",
			"-P", *proj,
			"-N", "ExomeDepth."+tag,
			strings.Join([]string{outdir, "ExomeDepth", "run.sh"}, pSep),
		)
	}

	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}

func runSMA(script, indir, outdir, geneInfo, control string, submit bool) {
	tag, _ := filepath.Abs(indir)
	tag = path.Base(tag)
	var args []string
	args = append(args, script)
	args = append(args, strings.Join([]string{outdir, "ExomeDepth", "sample.list.checked"}, pSep))
	args = append(args, geneInfo, control)
	args = append(args, strings.Join([]string{outdir, "SMA"}, pSep))
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	args2 = append(args2,
		"-cwd",
		"-l", "vf=10G,p=10",
		"-P", *proj,
		"-N", "SMA."+tag,
		strings.Join([]string{outdir, "SMA", "run_SMN_CNV_v2.sh"}, pSep),
	)

	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}
