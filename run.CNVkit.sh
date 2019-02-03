#!/bin/bash
path=$1
tag=$(basename $path)
prefix=$2
date=$(date --rfc-3339=date)
echo "# perl /share/backup/hanrui/pipeline/CNVkit/bin/run_CNVkit.pl $prefix/ExomeDepth/sample.list.checked /share/backup/hanrui/pipeline/CNVkit/control/MGISEQ_2000_control/201811/MGISEQ-2000_201811 $prefix/CNVkit"
perl /share/backup/hanrui/pipeline/CNVkit/bin/run_CNVkit.pl $prefix/ExomeDepth/sample.list.checked /share/backup/hanrui/pipeline/CNVkit/control/MGISEQ_2000_control/201811/MGISEQ-2000_201811 $prefix/CNVkit
echo "# submit cmd:"
echo qsub -cwd -l vf=30G,p=12 -P B2C_SGD -N CNVkit.$tag $prefix/CNVkit/run.sh
