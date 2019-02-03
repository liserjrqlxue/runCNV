#!/bin/bash
path=$1
tag=$(basename $path)
prefix=$2
date=$(date --rfc-3339=date)
echo "# perl /share/backup/wangyaoshen/src/ExomeDepth/createScript.sgd.pl $path $prefix/ExomeDepth $tag"
perl /share/backup/wangyaoshen/src/ExomeDepth/createScript.sgd.pl $path $prefix/ExomeDepth $tag
