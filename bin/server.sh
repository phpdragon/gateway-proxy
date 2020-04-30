#!/usr/bin/env bash

#. /etc/init.d/functions
source /etc/profile

CURRENT_DATETIME=$(date "+%Y-%m-%d %H:%M:%S")
CURRENT_DIR=$(pwd)
#app脚本目录
APP_BIN_DIR="$(cd "$(dirname "$0")" && pwd)"
#app根目录
APP_ROOT_DIR=$(cd "${APP_BIN_DIR}/../" && pwd)
#
#应用名称
APP_NAME=$(basename "$APP_ROOT_DIR")

#应用可执行文件路径
APP_BIN_FILE_PATH="${APP_BIN_DIR}/${APP_NAME}"

#
NOHUT_LOG_FILE="${APP_ROOT_DIR}/log/nohup.log"
NOHUP_BIN=$(whereis nohup | awk '{print $2}')
#
APP_STARTUP_CMD="${NOHUP_BIN} ${APP_BIN_FILE_PATH} >> ${NOHUT_LOG_FILE} 2>&1 &"
CHECK_PID_CMD="ps -fu${USER} | grep '${APP_BIN_FILE_PATH}' | grep -v grep | awk '{print \$2}'"
#
#提示文案
APP_NAME_COM="\033[31m${APP_NAME}\033[0m"
INFO_TIP_PREFIX="\033[32mINFO\033[0m:"
WARN_TIP_PREFIX="\033[33mWARN\033[0m:"
ERROR_TIP_PREFIX="\033[31mERROR\033[0m:"
#
NOHUP_NOT_INSTALLED="${WARN_TIP_PREFIX} nohup has not been installed and is now installed..."
APP_EXEC_CMD_NOT_EXIST="${ERROR_TIP_PREFIX} The executable file:[${APP_BIN_FILE_PATH}] that currently applies ${APP_NAME_COM} does not exist!"
#
WARIT_RUNNING_TIP="${INFO_TIP_PREFIX} waiting for the app ${APP_NAME_COM} startup"
RUNNING_TIP_COM="the app ${APP_NAME_COM} is \033[32mrunning\033[0m "
RUNNING_TIP="${INFO_TIP_PREFIX} ${RUNNING_TIP_COM}"
START_AGAIN_TIP="${WARN_TIP_PREFIX} ${RUNNING_TIP_COM}, please don't try again!"
#
WARIT_STOP_TIP="${INFO_TIP_PREFIX} waiting for the app ${APP_NAME_COM} stop"
STOPED_TIP_COM="the app ${APP_NAME_COM} is not running"
IS_STOPED_TIP="the app ${APP_NAME_COM} is stoped!"
STOPED_AGAIN_TIP="${WARN_TIP_PREFIX} ${STOPED_TIP_COM}, please don't try again!"
#
KILL_TIP="${INFO_TIP_PREFIX} exec kill ${APP_NAME_COM} process!"

# 进入app根目录下,解决程序中一些资源路径问题
cd "${APP_ROOT_DIR}" || exit 1

function start() {
  check_nohup

  check_bin_exist

  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    echo -e "${START_AGAIN_TIP} pid:\033[32m${app_pid}\033[0m"
    return_curr_dir
    exit 0
  fi

  echo -e "${INFO_TIP_PREFIX} exec start ${APP_NAME_COM} @${CURRENT_DATETIME}"
  echo -e "${INFO_TIP_PREFIX} ${APP_STARTUP_CMD}"
  eval "${APP_STARTUP_CMD}"

  printf "${WARIT_RUNNING_TIP}"
  app_pid=$(eval "${CHECK_PID_CMD}")
  while [[ -z "${app_pid}" ]]
  do
      printf "."
      sleep 1
  done

  echo ""
  echo -e "${RUNNING_TIP}, pid:\033[32m${app_pid}\033[0m !"
}

#检查可执行文件是否存在，应用部署目录必须和bin文件同名
function check_bin_exist() {
  if [ ! -f "${APP_BIN_FILE_PATH}" ]; then
    echo -e "${APP_EXEC_CMD_NOT_EXIST}"
    return_curr_dir
    exit 1
  fi
}

#检查nohup是否安装nohup
function check_nohup() {
  if [ -f "${NOHUP_BIN}" ]; then
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

function is_mac_os() {
  if (uname -a | grep -q "Darwin"); then
    return 0
  fi
  return 1
}

function is_ubuntu_os() {
  if [ -f "/etc/issue" ]; then
    # shellcheck disable=SC2002
    if (cat "/etc/issue" | grep -q "CentOS"); then
      return 0
    fi
  fi

  return 1
}

function is_centos() {
  if [ -f "/etc/redhat-release" ]; then
    # shellcheck disable=SC2002
    if (cat "/etc/redhat-release" | grep -q "CentOS"); then
      return 0
    fi
  fi

  return 1
}

function stop() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x == "${app_pid}"x ]; then
    echo -e "${STOPED_AGAIN_TIP}"
    return_curr_dir
    exit 0
  fi

  echo -e "${KILL_TIP}"
  kill_cmd="kill ${app_pid}"
  echo -e "${INFO_TIP_PREFIX} ${kill_cmd}"
  eval ${kill_cmd}

  printf "${WARIT_STOP_TIP}"
  while [[ -n "${app_pid}" ]]
  do
      printf "."
      sleep 1
      app_pid=$(eval "${CHECK_PID_CMD}")
  done

  echo ""
  echo -e "${IS_STOPED_TIP}"
}

function check_app_status() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    echo -e "${RUNNING_TIP}, pid:\033[32m${app_pid}\033[0m !"
  else
    echo -e "${STOPED_TIP}"
  fi
}

function return_curr_dir() {
  cd "${CURRENT_DIR}" || exit 1
}

function restart() {
  app_pid=$(eval "${CHECK_PID_CMD}")
  if [ ""x != "${app_pid}"x ]; then
    stop
  fi

  start
}

case $1 in
'start')
  start
  ;;
'stop')
  stop
  ;;
'restart')
  restart
  ;;
'status')
  check_app_status
  ;;
*)
  echo "USAGE:$0 {start|stop|restart|status}"
  exit 1
  ;;
esac

return_curr_dir

exit 0
