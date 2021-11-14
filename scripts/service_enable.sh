#!/bin/bash
set -euxo pipefail

service_name="gatestext.service"

systemctl is-active $service_name | grep -q "active" || systemctl start $service_name
