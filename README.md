dockersh
========

A shell which places uses into individual docker containers.

This is designed to be a (semi) secure method of giving users individual docker containers
on your host, without needing to run ssh in the containers (and have users connect to different ports).

*WARNING:* This project was implemented in 48 hours during a Yelp hackathon, it _should not_ be considered
stable, secure or ready for production use - here be dragons. Please expect to get rooted and/or for demons
to fly out of your nose if you use this software on a production host connected to the public internet.

*SECOND WARNING:* Whilst this project goes to some effort to make users inside containers low priviliged
(and therefore not be able to escalate their privelidge level inside containers or on the host machine),
this is *NOT* watertight, and will not be watertight until Docker fully supports user namespaces. Notably,
if you let users pick their own containers to run, they can probably do undesireable things (for example
using a container which allows them to sudo up to root and then writing to /dev/kmem).

Compiling dockersh
==================

You need to install golang (tested on >= 1.3), then you should just be able to run:

    go install

and a 'dockersh' binary will be generated in your $GOPATH (or .)

NOTE: dockersh requires a patched version of the 'nsenter' utility currently. It is recommended that
you remove any version of nsenter you have installed currently, then invoke dockersh, which will
tell you how to install the patched version.

Configuration
=============

We use the [XDG](http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html)
specification for configuration file locations. The default config file locations are shown below

~/.config/dockersh.json
-----------------------

Local (per user) settings for a specific user's dockersh instance.

Setting name  | Type | Description | Default value | Example value
------------- | ---- | ----------- | ------------- | -------------
image_name  | String | Mandatory, the name of the image to launch for the user. The %u sequence will interpolate the username | busybox | %u/mydockersh
mount_home_to | String | Where to map the user's home directory inside the container. Empty means don't mount home. | $HOME (from /etc/passwd) | /opt/home/myhomedir
container_username | String | Username which should be used inside the container. Defaults to %u (which is interpolated) | $USER | root
shell | String | The shell that should be started for the user inside the container | /bin/bash | /bin/ash
blacklist_user_config | Array of Strings | An array of configuration keys to disallow in per user dockershrc files | [] | ['container_username', 'mount_home', 'mount_home_to']

/etc/xdg/dockershrc.json
------------------------

Global settings for all dockersh instances. Allows you to disable settings
in the per-user dockersh.json 

Problems to solve
=================

 * How do we deal with changed settings (i.e. when to recycle the container)
 * Getting multiple shells into the same container (use of nsenter)
 * What becomes PID 1 inside the container? (sh while loop, but zombies?)

Contributing
============

Patches are very very welcome!

Please make a branch and send us a pull request.

Please ensure that you use the supplied pre-commit hook to correctly format your code:

    ln -s hooks/pre-commit .git/hooks/pre-commit

Copyright
=========

Copyright (c) 2014 Yelp. Some rights are reserved (see the LICENSE file for more details).

