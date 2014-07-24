#!/bin/sh

# Brutally murder all the running containers on a machine

docker ps | awk '{print $1}' | xargs docker kill
docker pa -s | awk '{print $1}' | xargs docker rm

