#!/bin/bash

CURRENT_DATETIME=$(date "+%Y-%m-%d %H:%M:%S")
CURRENT_DIR=$(pwd)

# shellcheck disable=SC2046
APP_NAME=$(basename $(pwd))
BIN_FILE="${CURRENT_DIR}/bin/${APP_NAME}"
# shellcheck disable=SC2034
CONF_FILE="${CURRENT_DIR}/etc/app.yaml"
NOHUT_LOG_FILE="${CURRENT_DIR}/log/nohup.log"

NOHUP_BIN_FILE=$(whereis nohup | awk '{print $2}')
APP_CMD="${NOHUP_BIN_FILE} ${BIN_FILE} -c ${CONF_FILE} > ${NOHUT_LOG_FILE} 2>&1 &"

HECK_PID_CMD="ps aux | grep \"${APP_NAME}\" | grep -v grep | awk '{print \$2}'"
APP_PID=$(eval "${HECK_PID_CMD}")

#检查nohup是否安装nohup
if [ ! -f "${NOHUP_BIN_FILE}" ] ; then
  echo "Warning: nohup has not been installed and is now installed..."
  /usr/bin/yum -y nohup
fi

#检查可执行文件是否存在，命令必须部署目录和bin文件同名
if [ ! -f ${BIN_FILE} ] ; then
  echo "Error: The executable file:[${BIN_FILE}] that currently applies ${APP_NAME} does not exist!"
  # shellcheck disable=SC2164
  cd "${CURRENT_DIR}"
  exit 1
fi

if [ ""x != "${APP_PID}"x ] ;then
  echo -e "The app: \033[31m${APP_NAME}\033[0m is \033[32mrunning\033[0m!, please don't try again!"
  # shellcheck disable=SC2164
  cd "${CURRENT_DIR}"
  exit 0
fi

echo "Notice: Exec cmd @${CURRENT_DATETIME}"
echo "${APP_CMD}"
# shellcheck disable=SC1036
eval "${APP_CMD}"

sleep 3

APP_PIDX=$(eval "${HECK_PID_CMD}")
if [ ""x != "${APP_PIDX}"x ]; then
 echo -e "Notice: the app \033[31m${APP_NAME}\033[0m is \033[32mrunning\033[0m!"
fi

# shellcheck disable=SC2164
cd "${CURRENT_DIR}"