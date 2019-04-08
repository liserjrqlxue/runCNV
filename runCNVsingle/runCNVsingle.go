package main

import (
	"flag"
	"fmt"
	"github.com/liserjrqlxue/simple-util"
	"os"
	"strings"
)

// os
var (
	pSep = string(os.PathSeparator)
)

var (
	outdir = flag.String(
		"outdir",
		"",
		"output workdir, \nworkdir\n\tCNVtype\n\t\tSampleID",
	)
	sampleID = flag.String(
		"sampleID",
		"",
		"sampleID for call CNVs",
	)
	bam = flag.String(
		"bam",
		"",
		"bam of sampleID",
	)
	gender = flag.String(
		"gender",
		"",
		"gender of sampleID",
	)
	run1 = flag.String(
		"run1",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/ExomeDepth/createScript.single.pl",
		"ExomeDepth script to run single",
	)
	run2 = flag.String(
		"run2",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/CNVkit/bin/run_CNVkit.single.pl",
		"CNVkit script to run single",
	)
	ExomeDepthControl = flag.String(
		"rds",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/ExomeDepth/test/all",
		"control of ExomeDepth, -rds.{gender,A}.my.count.rds",
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
	if *sampleID == "" || *bam == "" || *gender == "" {
		flag.Usage()
		os.Exit(0)
	}

	runExomeDepth(*run1, *sampleID, *gender, *bam, *outdir, *ExomeDepthControl, *submit)
	runCNVkit(*run2, *CNVkitControl, *outdir, *sampleID, *bam, *submit)

}

// outdir/tag/sampleID/run.sh

func runExomeDepth(script, sampleID, gender, bam, outdir, control string, submit bool) {
	tag := "ExomeDepth"
	var args []string
	args = append(args, script, sampleID, gender, bam, outdir+pSep+tag, control)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	args2 = append(
		args2,
		"-cwd",
		"-l", "vf=2G,p=1",
		"-P", "B2C_SGD",
		"-N", tag+"."+sampleID,
		strings.Join([]string{outdir, tag, sampleID, "run.sh"}, pSep),
	)

	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}

func runCNVkit(script, control, outdir, sampleID, bam string, submit bool) {
	tag := "CNVkit"
	var args []string
	args = append(args, script, control, outdir+pSep+tag, sampleID, bam)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	args2 = append(
		args2,
		"-cwd",
		"-l", "vf=12G,p=12G",
		"-P", "B2C_SGD",
		"-N", tag+"."+sampleID,
		strings.Join([]string{outdir, tag, sampleID, "run.sh"}, pSep),
	)

	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}
