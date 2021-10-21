#!/bin/bash

service_name="gatestext.service"

systemctl stop $service_name
rm -r /lib/systemd/system/$service_name
systemctl daemon-reload