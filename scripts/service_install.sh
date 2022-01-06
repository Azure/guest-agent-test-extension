#!/bin/bash
set -euxo pipefail

#           $1          $2
# Usage:  SystemdPath

service_name="gatestext.service"

dir=$(pwd)
cp $dir/bin/gatestext_script_linux /usr/local/bin
cp $dir/services/gatestext.service $1
chmod 644 $1/$service_name
systemctl daemon-reload
