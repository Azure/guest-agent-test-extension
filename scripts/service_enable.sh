#!/bin/bash
set -euxo pipefail

service_name="gatestext.service"

systemctl is-active $service_name || systemctl start $service_name
