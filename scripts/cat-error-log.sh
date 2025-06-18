#!/bin/bash
set -e 
if [ "$1" = "last" ]; then
   WORKDIR=$(ls -lt | grep workdir- | head -1 | awk '{print $9}')
else
   WORKDIR=$(ls -lt | grep workdir- | head -2 | tail -1  | awk '{print $9}')
fi
grep -A5  "ERROR" $WORKDIR/info.log