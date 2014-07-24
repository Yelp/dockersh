dockersh
========

A shell which places uses into individual docker containers

Configuration
=============

We use the [XDG](http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)
specification for configuration file locations. The default config file locations are shown below

~/.config/dockersh.json
-----------------------

Local (per user) settings for a specific user's dockersh instance.

Setting name  | Type | Description | Example value
------------- | ---- | ----------- | -------------
container_name  | String. | Mandatory, the name of the container to launch for the user. The %u sequence will interpolate the username | %u/mydockersh
mount_home  | Boolean. | If the user's home directory (as specified in /etc/passwd) should be mounted into the container | true/false
mount_home_to | False or string. | Where to map the user's home directory inside the container | /opt/home/myhomedir
container_username | String | Username which should be used inside the container. Defaults to %u (which is interpolated) | root
shell | String | The shell that should be started for the user inside the container | /bin/bash
blacklist_user_config | Array of Strings | An array of configuration keys to disallow in per user dockershrc files | ['container_username', 'mount_home', 'mount_home_to']

/etc/xdg/dockershrc.json
------------------------

Global settings for all dockersh instances. Allows you to disable settings
in the per-user dockersh.json 

Problems to solve
=================

 * How do we deal with changed settings (i.e. when to recycle the container)
 * Getting multiple shells into the same container (use of nsenter)
 * What becomes PID 1 inside the container? (sh while loop, but zombies?)

Copyright
=========

Copyright (c) 2014 Yelp. Some rights are reserved (see the LICENSE file for more details).

