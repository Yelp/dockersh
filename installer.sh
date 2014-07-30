#!/bin/bash

echo "dockersh installer - installs prebuilt dockersh binary"
echo ""
echo "To install dockersh"
echo "   docker run -v /usr/local/bin:/target thiscontainer"
echo "If you're using the publicly available (built from source) container, this is:"
echo "   docker run -v /usr/local/bin:/target Yelp/dockersh"
echo ""

if [ -d "/target" ];then
  echo "GOING TO DO INSTALL IN 5 SECONDS, Ctrl-C to abort"
  sleep 5
  cp -a /gopath/src/dockersh/dockersh /target/dockersh
else
  echo "No /target directory found, not installing"
fi

