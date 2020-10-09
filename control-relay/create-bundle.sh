#!/bin/bash

mkdir -p bin/netconf
mkdir -p bin/plugin-enabled
mkdir -p bin/plugin-repo
cp netconf/* bin/netconf

mkdir -p dist
cd bin
tar  czvf ../dist/control-relay.tgz  *
