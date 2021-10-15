#!/bin/bash

service_name="gatestext.service"

#Require sudo
if [ $EUID != 0 ]; then
    sudo "$0" "$@"
    exit $?
fi

dir=$(pwd)
eval "sed s#{WORKDIR}#$dir# $dir/services/gatestext.service > /lib/systemd/system/gatestext.service"
chmod 755 /lib/systemd/system/$service_name