#!/bin/bash

REAL_SHELL=/bin/ash

DOCKER_NAME=kwa_shell

DOCKER_CONTAINER=busybox

PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME")
if [ -z "$PID" ] || [ "$PID" == 0 ]; then
    docker run -t -i --name="$DOCKER_NAME" -d "$DOCKER_CONTAINER"
    PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME")
fi

sudo nsenter --target $PID --mount --uts --ipc --net --pid -- $REAL_SHELL
#nsenter --target $PID --uts --ipc --net --pid -- $REAL_SHELL
