#!/usr/bin/env bash

# shellcheck disable=SC2034
CURRENT_DATETIME=$(date "+%Y-%m-%d %H:%M:%S")
CURRENT_DIR=$(pwd)

# shellcheck disable=SC2046
APP_NAME=$(basename $(pwd))
BIN_FILE="${CURRENT_DIR}/bin/${APP_NAME}"
# shellcheck disable=SC2034
CONF_FILE="${CURRENT_DIR}/etc/app.yaml"
NOHUT_LOG_FILE="${CURRENT_DIR}/log/nohup.log"

NOHUP_BIN_FILE=$(whereis nohup | awk '{print $2}')
# shellcheck disable=SC2034
APP_CMD="${NOHUP_BIN_FILE} ${BIN_FILE} -c ${CONF_FILE} > ${NOHUT_LOG_FILE} 2>&1 &"

HECK_PID_CMD="ps aux | grep \"${APP_NAME}\" | grep -v grep | awk '{print \$2}'"
APP_PID=$(eval "${HECK_PID_CMD}")

if [ ""x == "${APP_PID}"x ] ;then
  echo -e "Notice: the app \033[31m${APP_NAME}\033[0m is not \033[31mrunning\033[0m!, please don't try again!"
  # shellcheck disable=SC2164
  cd "${CURRENT_DIR}"
  exit 0
fi

# shellcheck disable=SC2086
kill -9 "${APP_PID}"

sleep 5

APP_PIDX=$(eval "${HECK_PID_CMD}")
if [ ""x == "${APP_PIDX}"x ]; then
  echo -e "Notice: the app \033[31m${APP_NAME}\033[0m is not \033[31mrunning\033[0m!"
fi

# shellcheck disable=SC2164
cd "${CURRENT_DIR}"