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
		"output workdir:\nworkdir\n\tCNVtype\n\t\tSampleID\n\t\t\trun.sh",
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
		"gender of sampleID,[F/M]",
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
	run3 = flag.String(
		"run3",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/wangyaoshen/SMA_WES/run_SMN_CNV_v2.single.pl",
		"SMA script to run single sample",
	)
	ExomeDepthControl = flag.String(
		"rds",
		"/ifs7/B2C_SGD/PROJECT/PP12_Project/ExomeDepth/test/all",
		"control of ExomeDepth, -rds.{gender,A}.my.count.rds",
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
	proj = flag.String(
		"proj",
		"B2C_SGD",
		"project of qsub -P",
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
	runSMA(*run3, *bam, *SMAGene, *SMAControl, *outdir, *sampleID, *submit)
}

// outdir/tag/sampleID/run.sh

func runExomeDepth(script, sampleID, gender, bam, outdir, control string, submit bool) {
	tag := "ExomeDepth"
	var args []string
	args = append(args, script, sampleID, gender, bam, outdir+pSep+tag, control)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	runSubmit(outdir, tag, sampleID, "2G", "1", submit)
}

func runCNVkit(script, control, outdir, sampleID, bam string, submit bool) {
	tag := "CNVkit"
	var args []string
	args = append(args, script, control, outdir+pSep+tag, sampleID, bam)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	runSubmit(outdir, tag, sampleID, "12G", "10", submit)
}

func runSMA(script, bam, geneInfo, control, outdir, sampleID string, submit bool) {
	tag := "SMA"
	var args []string
	args = append(args, script, bam, geneInfo, control, outdir+pSep+tag, sampleID)
	fmt.Printf("# perl %s\n", strings.Join(args, " "))
	simple_util.RunCmd("perl", args...)

	runSubmit(outdir, tag, sampleID, "10G", "1", submit)
}

func runSubmit(outdir, tag, sampleID, vf, p string, submit bool) {
	var args2 []string
	args2 = append(
		args2,
		"-cwd",
		"-l", "vf="+vf+",p="+p,
		"-P", *proj,
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
