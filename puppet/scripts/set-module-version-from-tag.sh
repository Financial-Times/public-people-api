#!/bin/bash

VERSION=`git describe --abbrev=0 --tags`
MODULE=`find ../ -name Modulefile`
sed -e s/^version\ '.*'$/version\ \'${VERSION}\'/ ${MODULE} > Modulefile.tmp && mv Modulefile.tmp ${MODULE}
