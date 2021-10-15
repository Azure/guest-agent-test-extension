#!/bin/bash

service_name="gatestext.service"

#Require sudo
if [ $EUID != 0 ]; then
    sudo "$0" "$@"
    exit $?
fi

systemctl daemon-reload
systemctl start $service_name