#!/bin/bash

service_name="gatestext.service"

if systemctl is-active $service_name | grep -q "inactive"; then
    systemctl start $service_name
fi
