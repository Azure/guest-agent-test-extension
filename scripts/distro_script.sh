#!/bin/bash
set -euo pipefail

distro=None
if [ -f /etc/os-release ]; then
    # freedesktop and systemd
    . /etc/os-release
    distro=$ID
elif type lsb_release >/dev/null 2>&1; then
    # linuxbase.org
    distro=$(lsb_release -si)
elif [ -f /etc/lsb-release ]; then
    # For some versions of Debian/Ubuntu without lsb_release command
    . /etc/lsb-release
    distro=$DISTRIB_ID
elif [ -f /etc/debian_version ]; then
    # Older Debian/Ubuntu/etc.
    distro=Debian
elif [ -f /etc/SuSe-release ]; then
    # Older SuSE/etc.
    distro=Suse
elif [ -f /etc/centos-release ]; then
    # Older CentOS.
    distro=CentOS
elif [ -f /etc/redhat-release ]; then
    # Older Red Hat.
    distro=RedHat
else
    # Fall back to uname, e.g. "Linux <distro>", also works for BSD, etc.
    distro=$(uname -s)
fi

echo $distro
