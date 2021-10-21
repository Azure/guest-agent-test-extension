#!/bin/bash

service_name="gatestext.service"

systemctl is-active $service_name | grep "inactive" && systemctl start $service_name
