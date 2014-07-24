#!/bin/bash

# TODO: Figure out they want for a real shell
REAL_SHELL=/bin/ash

# Must use a conistent naming scheme, docker will only let one of these run
# at a time.
DOCKER_NAME="$USER_shell"

# TODO: Figure out what they want from config
DOCKER_CONTAINER=busybox

PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME")
# If we got here, then the docker is not running.
if [ -z "$PID" ] || [ "$PID" == 0 ]; then
    # If the docker is stopped, we must remove it and start a new one
    docker rm --name="$DOCKER_NAME"
    # TODO: Configur the bind mounts
    docker run -t -i --name="$DOCKER_NAME" -v "$HOME":/root/:rw -d "$DOCKER_CONTAINER"
    PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME")
fi

sudo nsenter --target $PID --mount --uts --ipc --net --pid -- $REAL_SHELL
