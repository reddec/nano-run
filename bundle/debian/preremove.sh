#!/bin/sh

SERVICE="nano-run"
RUNNING_USER="${SERVICE}"

systemctl stop "${SERVICE}".service || echo "failed to stop service"
systemctl disable "${SERVICE}".service || echo "failed to disable service"

if id -u ${RUNNING_USER}; then
  echo "Removing user ${RUNNING_USER}..."
  userdel -r ${RUNNING_USER}
fi


