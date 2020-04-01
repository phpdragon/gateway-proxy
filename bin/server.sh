#!/usr/bin/env bash

#. /etc/init.d/functions

CURRENT_DATETIME=$(date "+%Y-%m-%d %H:%M:%S")
CURRENT_DIR=$(pwd)
#
#当前脚本目录
APP_BIN_DIR="$(cd "$(dirname "$0")" && pwd)"
#app根目录
APP_ROOT_DIR=$(cd "${APP_BIN_DIR}/../" && pwd)
#
#应用名称
APP_NAME=$(basename "$APP_ROOT_DIR")
APP_BIN_FILE="${APP_BIN_DIR}/${APP_NAME}"
APP_CONF_FILE="${APP_ROOT_DIR}/etc/app.yaml"
#
NOHUT_LOG_FILE="${APP_ROOT_DIR}/log/nohup.log"
NOHUP_APP_BIN_FILE=$(whereis nohup | awk '{print $2}')
#
APP_START_CMD="${NOHUP_APP_BIN_FILE} ${APP_BIN_FILE} -c ${APP_CONF_FILE} > ${NOHUT_LOG_FILE} 2>&1 &"
CHECK_PID_CMD="ps aux | grep \"${APP_NAME}\" | grep -v grep | awk '{print \$2}'"
#
#提示文案
APP_NAME_COM="\033[31m${APP_NAME}\033[0m"
INFO_TIP_PREFIX="\033[32mINFO\033[0m:"
WARN_TIP_PREFIX="\033[33mWARN\033[0m:"
ERROR_TIP_PREFIX="\033[31mERROR\033[0m:"
#
NOHUP_NOT_INSTALLED="${WARN_TIP_PREFIX} nohup has not been installed and is now installed..."
APP_BIN_NOT_EXIST="${ERROR_TIP_PREFIX} The executable file:[${APP_BIN_FILE}] that currently applies ${APP_NAME_COM} does not exist!"
#
RUNNING_TIP_COM="the app ${APP_NAME_COM} is \033[32mrunning\033[0m "
RUNNING_TIP="${INFO_TIP_PREFIX} ${RUNNING_TIP_COM}"
START_AGAIN_TIP="${WARN_TIP_PREFIX} ${RUNNING_TIP}, please don't try again!"
#
STOPED_TIP_COM="the app ${APP_NAME_COM} is not running"
STOPED_TIP="${INFO_TIP_PREFIX} ${STOPED_TIP_COM}!"
STOPED_AGAIN_TIP="${WARN_TIP_PREFIX} ${STOPED_TIP_COM}, please don't try again!"
#
KILL_TIP="${INFO_TIP_PREFIX} exec kill ${APP_NAME_COM} process!"

start() {
  check_nohup
  check_bin_exist

  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    echo -e "${START_AGAIN_TIP}"
    return_curr_dir
    exit 0
  fi

  echo -e "${INFO_TIP_PREFIX} exec start ${APP_NAME_COM} @${CURRENT_DATETIME}"
  echo -e "${INFO_TIP_PREFIX} ${APP_START_CMD}"
  eval "${APP_START_CMD}"

  sleep 3

  check_app_status
}

#检查可执行文件是否存在，命令必须部署目录和bin文件同名
check_bin_exist() {
  if [ ! -f "${APP_BIN_FILE}" ]; then
    echo -e "${APP_BIN_NOT_EXIST}"
    return_curr_dir
    exit 1
  fi
}

#检查nohup是否安装nohup
check_nohup() {
  if [ -f "${NOHUP_APP_BIN_FILE}" ]; then
    return 0
  fi

  echo -e "${NOHUP_NOT_INSTALLED}"

  if (is_centos); then
    # shellcheck disable=SC2034
    yum_bin_exist=$(whereis yum | awk '{print $2}')
    if [ ! -f "${yum_bin_exist}" ]; then
      yum install -y nohup
    fi
  elif (is_mac_os); then
    # shellcheck disable=SC2034
    brew_bin_exist=$(whereis brew | awk '{print $2}')
    if [ ! -f "${brew_bin_exist}" ]; then
      brew install -y nohup
    fi
  elif (is_ubuntu_os); then
    # shellcheck disable=SC2034
    apt_get_bin_exist=$(whereis apt-get | awk '{print $2}')
    if [ ! -f "${apt_get_bin_exist}" ]; then
      apt-get install -y nohup
    fi
  fi
}

is_mac_os() {
  if (uname -a | grep -q "Darwin"); then
    return 0
  fi
  return 1
}

is_ubuntu_os(){
  if [ -f "/etc/issue" ]; then
    # shellcheck disable=SC2002
    if (cat "/etc/issue" | grep -q "CentOS"); then
      return 0
    fi
  fi

  return 1
}

is_centos() {
  if [ -f "/etc/redhat-release" ]; then
    # shellcheck disable=SC2002
    if (cat "/etc/redhat-release" | grep -q "CentOS"); then
      return 0
    fi
  fi

  return 1
}

stop() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x == "${app_pid}"x ]; then
    echo -e "${STOPED_AGAIN_TIP}"
    return_curr_dir
    exit 0
  fi

  echo -e "${KILL_TIP}"
  kill -9 "${app_pid}"

  sleep 3

  check_app_status
}

check_app_status() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    echo -e "${RUNNING_TIP}, pid:${app_pid} !"
  else
    echo -e "${STOPED_TIP}"
  fi
}

return_curr_dir() {
  cd "${CURRENT_DIR}" || exit 1
}

restart() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    stop
  fi

  start
}

case $1 in
start)
  start
  return_curr_dir
  ;;
stop)
  stop
  return_curr_dir
  ;;
restart)
  restart
  return_curr_dir
  ;;
status)
  check_app_status
  return_curr_dir
  ;;
*)
  echo "USAGE:$0 {start|stop|restart|status}"
  exit 1
  ;;
esac
