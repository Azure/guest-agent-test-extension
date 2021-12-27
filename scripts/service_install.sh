#!/bin/bash
set -euxo pipefail

service_name="gatestext.service"

dir=$(pwd)
cp $dir/bin/gatestext_script_linux /usr/local/bin
eval "sed s#{WORKDIR}#$dir# $dir/services/gatestext.service > /lib/systemd/system/gatestext.service"
chmod 644 /lib/systemd/system/$service_name
systemctl daemon-reload
