#!/bin/bash
path=$1
tag=$(basename $path)
prefix=${2-$(basename $path)}

#date=$(date --rfc-3339=date)
sh run.ExomeDepth.sh $path $prefix
sh run.CNVkit.sh     $path $prefix
