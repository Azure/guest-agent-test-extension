#!/bin/bash

service_name="gatestext.service"

#Require sudo
if [ $EUID != 0 ]; then
    sudo "$0" "$@"
    exit $?
fi

systemctl stop $service_name
rm -r /lib/systemd/system/$service_name