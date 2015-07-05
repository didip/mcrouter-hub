#!/bin/bash

set -e

yum install -y epel-release git bzip2-devel libevent-devel libcap-devel scons unzip gflags-devel \
openssl-devel bison flex snappy-devel numactl-devel cyrus-sasl-devel cmake libtool \
glibc-devel.i686 glibc-devel.x86_64 gcc gcc-c++ zlib-devel autoconf automake \
double-conversion double-conversion-devel boost boost-devel glog glog-devel thrift thrift-devel golang memcached

#Folly
if [ ! -d "/opt/folly-0.47.0" ]; then
    cd /opt && curl https://codeload.github.com/facebook/folly/tar.gz/v0.47.0 | tar xvz
    cd folly-0.47.0/folly

    export LD_LIBRARY_PATH="/opt/folly-0.47.0/folly/lib:$LD_LIBRARY_PATH"
    export LD_RUN_PATH="/opt/folly-0.47.0/folly/lib"
    export LDFLAGS="-L/opt/folly-0.47.0/folly/lib -L/usr/local/lib -ldl"
    export CPPFLAGS="-I/opt/folly-0.47.0/folly/include"
    autoreconf -ivf
    ./configure
    make && make install
fi

#Ragel
if [ ! -d "/opt/ragel-6.9" ]; then
    cd /opt && curl http://www.colm.net/files/ragel/ragel-6.9.tar.gz | tar zx
    cd /opt/ragel-6.9 && ./configure && make && make install
fi

#McRouter
if ! command -v mcrouter >/dev/null; then
    cd /opt && git clone https://github.com/facebook/mcrouter.git || true
    cd mcrouter/mcrouter
    export LDFLAGS="-L/opt/folly-0.47.0/folly/lib -L/usr/local/lib -ldl"
    export CXXFLAGS="-fpermissive"
    autoreconf --install && ./configure
    make && make install
    mkdir -p /var/spool/mcrouter
    mkdir -p /var/mcrouter/stats
    mcrouter --help
fi