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
	run2 = flag.String(
		"run2",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/CNVkit/bin/run_CNVkit.single.pl",
		"CNVkit script to run single",
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
	if *sampleID == "" || *bam == "" {
		flag.Usage()
		os.Exit(0)
	}

	runCNVkit(*run2, *CNVkitControl, *outdir, *sampleID, *bam, *submit)

}

func runCNVkit(script, control, outdir, sampleID, bam string, submit bool) {
	var args []string
	args = append(args, script, control, outdir, sampleID, bam)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	var args2 []string
	args2 = append(
		args2,
		"-cwd",
		"-l", "vf=12G,p=12G",
		"-P", "B2C_SGD",
		"-N", "CNVkit."+sampleID,
		strings.Join([]string{outdir, "CNVkit", sampleID, "run.sh"}, pSep),
	)

	if submit {
		fmt.Printf("# qsub %s\n", strings.Join(args2, " "))
		simple_util.RunCmd("qsub", args2...)
	} else {
		fmt.Printf("# submit cmd:\nqsub %s\n", strings.Join(args2, " "))
	}
}
