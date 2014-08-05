dockersh
========

A user shell for isolated environments.

What is this?
=============

dockersh is designed to be used as a login shell on machines with multiple interactive users.

When a user invokes dockersh, it will bring up a Docker container (if not already running), and
then spawn a new interactive shell in the container's namespace.

dockersh can be used as a shell in /etc/passwd or as an ssh ForceCommand.


This allows you to have a single ssh process, on the normal ssh port, and gives
a secure way to connect users into their own individual docker
containers.

Why do I want this?
===================

You want to allow multiple users to ssh onto a single box, but you'd like some isolation
between those users. With dockersh each user enters their
own individual docker container (acting like a lightweight VM), with their homedirectory mounted from the host
system (so that user data is persistent between container restarts), but with it's own kernel namespaces for
processes and networking.

This can be used for more easily seperating each user's processes from the rest of the system
and having further per user constraints (e.g. memory limit all of the user's processes,
or limit their aggregate bandwidth etc)

Normally to give users individual containers you have to run an ssh daemon in each
container, and either have have a different port for each user to ssh to or some nasty
Forcecommand hacks (which only work with agent forwarding from the client).

Dockersh eliminates the need for any of these techiques by acting like a regular
shell which can be used in /etc/passwd or as an ssh ForceCommand.  
This allows you to have a single ssh process, on the normal ssh port, and gives
a (hopefully) secure way to connect users into their own individual docker
containers.

SECURITY WARNING
================

dockersh tries hard to drop all priviliges as soon as possible, including disabling 
the suid, sgid, raw sockets and mknod capabilities of the target process (and all children)

*WARNING:* This project was implemented in 48 hours during a Yelp hackathon, it _should not_ be considered
stable, secure or ready for production use - here be dragons. Please expect to get rooted and/or for demons
to fly out of your nose if you use this software on a production host connected to the public internet.

*SECOND WARNING:* Whilst this project goes to some effort to make users inside containers have lowered privileges
and limit their ability to escalate their privilege level inside containers, or on the host machine,
this is *NOT* watertight. It will not be watertight until Docker fully supports user namespaces. Notably,
if you let users pick their own containers to run, they can probably do undesireable things.

*THIRD WARNING:* The dockersh binary needs the suid bit set so that it can make the syscalls to adjust
kernel namespaces, so any security issues in this code *will* allow attackers to escalate to root.

Requirements
============

amd64 platforms
---------------

Compiles down into a single binary with no external dependencies - see 'Installation' below.

darwin
------

dockersh tries to support Mac environments e.g. boot2docker (however at this time the solution is less
optimum).

dockersh requires a patched version of the 'nsenter' utility on the target machine *if* you want to
use it from darwin (E.g. in boot2docker). This version of nsenter needs to be installed inside the
boot2docker VM.

It is recommended that
you remove any version of nsenter you have installed currently, then invoke dockersh, which will
tell you how to install the patched version.

Installation
============

With docker
-----------

(This is the recommended method).
Build the Dockerfile in the local directory into an image, and run it like this:

    $ docker build .
    # Progress, takes a while the first time..
    ....
    Successfully built 3006a08eef2e 
    $ docker run -v /usr/local/bin:/target 3006a08eef2e

Without docker
--------------

You need to install golang (tested on >= 1.3), then you should just be able to run:

    go get
    make

and a 'dockersh' binary will be generated in your $GOPATH (or your current
working directory if $GOPATH isn't set). N.B. This binary needs to be moved to where
you would like to install it (recommended /usr/local/bin), and owned by root + u+s
(suid). This is done automatically if you use the Docker based installed, but
you need to do it manually if you're compiling the binary yourself.

Invoking dockersh
=================

There are two main methods of invoking dockersh. Either:

1. Put the path to dockersh into /etc/shells, and then change the users shell
   in /etc/passwd (e.g. chsh myuser -s /usr/local/bin/dockersh)
1. Set dockersh as the ssh ForceCommand in the users $HOME/.ssh/config, or
   globally in /etc/ssh/ssh_config

*Note:* The dockersh binary needs the suid bit set to operate!

Configuration
=============

We use (https://code.google.com/p/gcfg/)[gocfg] to read configs in an ini style format.

Each file has a [docker] block in it, and zero or more [user "foo"] blocks.

This can be used to enable or disable setting on a per user basis.

Each user can (optionally) have a per user config, although this can be disabled.

~/.dockersh
-----------

Local (per user) settings for a specific user's dockersh instance.

Settings should be contained prefixed with a [docker] block

Setting name  | Type | Description | Default value | Example value
------------- | ---- | ----------- | ------------- | -------------
imagename  | String | Mandatory, the name of the image to launch for the user. The %u sequence will interpolate the username | busybox | ubuntu, or %u/mydockersh
mounthome | Bool | If the users home directory should be mounted in the target container | true | false
mounttmp | String | If /tmp should be mounted into the target container (so that ssh agent forwarding works). N.B. Security risk | false | true
mounthometo | String | Where to map the user's home directory inside the container. | $HOME (from /etc/passwd) | /opt/home/myhomedir
mounthomefrom | String | Where to map the user's home directory from on the host. | $HOME (from /etc/passwd) | /opt/home/%u
containerusername | String | Username which should be used inside the container. Defaults to %u (which is interpolated) | %u | root
shell | String | The shell that should be started for the user inside the container. | /bin/ash | /bin/bash
mountdockersocket | Bool | If to mount the docker socket from the host. (DANGEROUS) | false | true
dockersocket | String | The location of the docker socket from the host. | /var/run/docker.sock | /opt/docker/var/run/docker.sock
entrypoint | String | The entrypoint for the persistent process to keep the container running | internal | /sbin/yoursupervisor

/etc/dockershrc
---------------

Global settings for all dockersh instances. Allows you to set settings for all users (in a [docker block])
or specific users (in a [user "username"] block), on enable setting settings in per user  ~/.dockersh

Setting name  | Type | Description | Default value | Example value
------------- | ---- | ----------- | ------------- | -------------
enableuserconfig | bool | Set to true to enable reading of per user ~/.dockersh files | false | true
enableuserimagename | bool | Set to true to enable reading of imagename parameter from ~/.dockersh files | false | true
enableusermounthome | bool | Set to true to enable reading of mounthome parameter from ~/.dockersh files | false | true
enableusermounttmp | bool | Set to true to enable reading of mounttmp parameter from ~/.dockersh files | false | true
enableusermounthometo | bool | Set to true to enable reading of mounthometo parameter from ~/.dockersh files | false | true
enableusermounthomefrom | bool | Set to true to enable reading of mounthomefrom parameter from ~/.dockersh files | false | true
enableusercontainerusername | bool | Set to true to enable reading of containerusername parameter from ~/.dockersh files | false | true
enableusershell | bool | Set to true to enable reading of shell parameter from ~/.dockersh files | false | true
enableuserentrypoint | bool | Set to true to enable users to set their own supervisor daemon / entry point to the container for PID 1 | false | true

Example configs
---------------

Note the liberal use of the blacklistuserconfig

Sets up a fairly restricted shell environment, with one admin user being allowed additional privs, set the following /etc/dockersh 

    [dockersh]
    imagename = ubuntu:precise
    shell = /bin/bash
    mounthome

    [user "someadminguy"]
    mounttmp
    mountdockersocket
    
In a less restrictive environment, you may allow users to choose their own container and shell, from a 'shell' container
they have uploaded to the registry, and have ssh agent forwarding working, with the following /etc/dockersh

    [dockersh]
    imagename = "%u/shell"
    mounthome
    mounttmp
    enableuserconfig
    blacklistuserconfig = imagename
    blacklistuserconfig = mounthometo
    blacklistuserconfig = mountdockersocket
    blacklistuserconfig = dockersocket

    [user "someadminguy"]
    mountdockersocket

And an example user's ~/.dockersh

    [dockersh]
    shell = /bin/zsh

Or just allowing your users to run whatever container they want:

    [dockersh]
    mounthome
    mounttmp
    enableuserconfig
    blacklistuserconfig = imagename
    blacklistuserconfig = mounthometo
    blacklistuserconfig = mountdockersocket
    blacklistuserconfig = dockersocket
    
TODO
====

 * How do we deal with changed settings (i.e. when to recycle the container)
    * Document just kill 1 inside the container?
 * Fix up go panics when eixting the root container.
 * More config settings?
 * getpwnam so that we can interpolate the user's shell from /etc/shells (if used in ForceCommand mode!)
 * Change config over to be INI style
    * This would be nicer, as we could also add global per user config as [username] type sections
 * Decent test cases
 * Make the darwin nsenter version less crazy - or kill as less features?
 * Allow setting the max memory for the container's processes
 * Allow setting the CMD of the root process
 * Allow setting the entrypoint of the root process to be something other than "internal"
 * Find a better way to make ssh agent sockets work than to bind /tmp
 * Expose ability to mount additional volumes in the config
 * Expose ability to pass arbitrary options to docker in the config.

Contributing
============

Patches are very very welcome!

This is our first real Go project, so we apologise about the shoddy quality of the code.

Please make a branch and send us a pull request.

Please ensure that you use the supplied pre-commit hook to correctly format your code
with go fmt:

    ln -s hooks/pre-commit .git/hooks/pre-commit

Copyright
=========

Copyright (c) 2014 Yelp. Some rights are reserved (see the LICENSE file for more details).

