#!/bin/sh

SERVICE="nano-run"
RUNNING_USER="${SERVICE}"

if ! id -u ${RUNNING_USER}; then
  echo "Creating user ${RUNNING_USER}..."
  useradd -M -c "${RUNNING_USER} dummy user" -r -s /bin/nologin ${RUNNING_USER}
fi

systemctl enable "${SERVICE}".service || echo "failed to enable service"
systemctl start "${SERVICE}".service || echo "failed to start service"
