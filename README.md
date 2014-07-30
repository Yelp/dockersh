dockersh
========

If you want to allow multiple users to ssh onto a single box, and enter their
own individual docker container, you normally have to run an ssh daemon in each
container, and have a different port for each user to ssh to.

dockersh is a shell which can be used in /etc/passwd or as an ssh ForceCommand.  
This allows you to have a single ssh process, on the normal ssh port, and gives
a (semi) secure way to connect users into their own individual docker
containers.

*WARNING:* This project was implemented in 48 hours during a Yelp hackathon, it _should not_ be considered
stable, secure or ready for production use - here be dragons. Please expect to get rooted and/or for demons
to fly out of your nose if you use this software on a production host connected to the public internet.

*SECOND WARNING:* Whilst this project goes to some effort to make users inside containers have lowered privileges
and limit their ability to escalate their privilege level inside containers, or on the host machine,
this is *NOT* watertight. It will not be watertight until Docker fully supports user namespaces. Notably,
if you let users pick their own containers to run, they can probably do undesireable things (for example
using a container which allows them to sudo up to root and then writing to /dev/kmem.). We plan to try to
address some of this by limiting suid/sgid permissions within containers, but YMMV.

*THIRD WARNING:* The dockersh binary needs the suid bit set so that it can make the syscalls to adjust
kernel namespaces, so any security issues in the code *will* allow attackers to escalate to root.

Requirements
============

amd64 platforms
---------------

Compiles down into a single binary with no external dependencies - see 'Compiling dockersh' below.

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

    make

and a 'dockersh' binary will be generated in your $GOPATH (or your current
working directory if $GOPATH isn't set). N.B. This binary needs to be moved to where
you would like to install it, and set user + suid

Invoking dockersh
=================

There are two main methods of invoking dockersh. Either:

1. Put the path to dockersh into /etc/shells, and then change the users shell
   in /etc/passwd
1. Set dockersh as the ssh ForceCommand in the users $HOME/.ssh/config, or
   globally in /etc/ssh/ssh_config

*Note:* The dockersh binary needs the suid bit set to operate!

Configuration
=============

NOTE: Currently, the presense of ~/.dockersh.json will _entirely_ override the global config file, settings not present in ~/.dockersh.json
will *not* be merged from global config!

~/.dockersh.json
----------------

Local (per user) settings for a specific user's dockersh instance.

Setting name  | Type | Description | Default value | Example value
------------- | ---- | ----------- | ------------- | -------------
image_name  | String | Mandatory, the name of the image to launch for the user. The %u sequence will interpolate the username | busybox | %u/mydockersh
mount_home_to | String | Where to map the user's home directory inside the container. Empty means don't mount home. | $HOME (from /etc/passwd) | /opt/home/myhomedir
container_username | String | Username which should be used inside the container. Defaults to %u (which is interpolated) | %u | root
shell | String | The shell that should be started for the user inside the container. | /bin/bash | /bin/ash


/etc/dockershrc.json
--------------------

Global settings for all dockersh instances. Allows you to disable settings
in the per-user ~/.dockersh.json

N.B *TODO* Not yet implemented, reading both config files, and allowing the global one to blacklist local settings:

Setting name  | Type | Description | Default value | Example value
------------- | ---- | ----------- | ------------- | -------------
disable_user_config | bool | Set to true to disable ~/.dockersh reading entirely | false | true
blacklist_user_config | Array of Strings | An array of configuration keys to disallow in per user dockershrc files | [] | ['container_username', 'mount_home', 'mount_home_to']

TODO List
=========

 * How do we deal with changed settings (i.e. when to recycle the container)
 * We just run an interactive shell in the root of the container, but if you 'docker attach' to it, then detach, the container goes away.
 * Finish up config settings
   * Fix getpwnam so that we can interpolate the user's shell from /etc/shells (if used in ForceCommand mode!)
   * Add config merging (so user can override global)
   * Add global user config lock out settings
 * Decent test cases
 * Make the darwin nsenter version less crazy
 * suid / sgid binaries inside the container - disable

Contributing
============

Patches are very very welcome!

Please make a branch and send us a pull request.

Please ensure that you use the supplied pre-commit hook to correctly format your code:

    ln -s hooks/pre-commit .git/hooks/pre-commit

Copyright
=========

Copyright (c) 2014 Yelp. Some rights are reserved (see the LICENSE file for more details).

