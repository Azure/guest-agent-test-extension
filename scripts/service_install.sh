#!/bin/bash

service_name="gatestext.service"

dir=$(pwd)
eval "sed s#{WORKDIR}#$dir# $dir/services/gatestext.service > /lib/systemd/system/gatestext.service"
chmod 755 /lib/systemd/system/$service_name
systemctl daemon-reload