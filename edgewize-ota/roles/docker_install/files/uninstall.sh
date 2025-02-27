#!/bin/bash
systemctl stop docker

rm -f /usr/bin/docker*
rm -f /usr/bin/containerd*
rm -f /usr/bin/ctr
rm -f /usr/bin/runc

rm -rf /etc/systemd/system/docker.service

rm -rf /var/lib/docker*

rm -rf /var/run/docker*

rm -rf /var/run/docker.pid

systemctl daemon-reload

