#!/bin/bash

# This is the initial prototype version of dockersh, just to prove that it could be done.

# WARNING: This code is _entirely unsupported_ and a barely functional
#          prototype, please use the .go version if you'd actually like
#          to experiment with this project!

# TODO: Figure out they want for a real shell
REAL_SHELL=/bin/ash

DESIRED_USER="nobody"
# Must use a conistent naming scheme, docker will only let one of these run
# at a time.
DOCKER_NAME="${DESIRED_USER}_shell"

# TODO: Figure out what they want from config
DOCKER_CONTAINER=busybox

DESIRED_USER="vagrant"

# FIXME - We should instead work out the target UID inside the destination container?
DESIRED_UID=$(id -u $DESIRED_USER)
DESIRED_GID=$(id -g $DESIRED_USER)
HOMEDIR=$(eval echo ~$DESIRED_USER)
MYHOSTNAME="$(hostname --fqdn)-${DESIRED_USER}-docker"

PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME" 2>/dev/null)
# If we got here, then the docker is not running.
if [ -z "$PID" ] || [ "$PID" == 0 ]; then
    # If the docker is stopped, we must remove it and start a new one
    docker rm --name="$DOCKER_NAME" >/dev/null 2>&1 # May not be running, just throw away the output
    # TODO: Configur the bind mounts
    # FIXME - If you docker attach to this container, then Ctrl-D, it dies. (This is expected?)
    docker run -t -i -u $DESIRED_USER --hostname="$MYHOSTNAME" --name="$DOCKER_NAME" -v $HOMEDIR:$HOMEDIR:rw -v /etc/passwd:/etc/passwd:ro -v /etc/group:/etc/group:ro -d "$DOCKER_CONTAINER"
    PID=$(docker inspect --format {{.State.Pid}} "$DOCKER_NAME")
fi

# N.B. You need to bobtfish/nsenter version of nsenter for suid/sgid to do the right thing.
sudo nsenter --target "$PID" --mount --uts --ipc --net --pid --setuid $DESIRED_UID --setgid $DESIRED_GID --wd=$HOMEDIR -- "$REAL_SHELL"

