dockersh
========

A shell which places uses into individual docker containers

Configuration
=============

~/.dockershrc.json
------------------

Local (per user) settings for a specific user's dockersh instance.

Setting name  | Description   | Example value
------------- | ------------- | -------------
container_name  | String. Mandatory, the name of the container to launch for the user. The %u sequence will interpolate the username | %u/mydockersh
mount_home  | Boolean. If the user's home directory (as specified in /etc/passwd) should be mounted into the container | true/false
mount_home_to | False or string. Where to map the user's home directory inside the container | /opt/home/myhomedir
container_username | String Username which should be used inside the container. Defaults to %u (which is interpolated) | root


/etc/dockershrc.json
--------------------

Global settings for all dockersh instances. Allows you to disable settings
in ~/.dockershrc


Copyright
=========

Copyright (c) 2014 Yelp. Some rights are reserved (see the LICENSE file for more details).

