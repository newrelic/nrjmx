FROM ubuntu:20.04 as builder

ARG THRIFT_VERSION='0.16.0'

ENV TZ=Etc/UCT
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update \
    && apt-get install -y automake bison flex g++ git libboost-all-dev libevent-dev libssl-dev libtool make pkg-config curl

RUN curl -L https://github.com/apache/thrift/archive/refs/tags/v${THRIFT_VERSION}.tar.gz | tar -xz

RUN cd thrift-${THRIFT_VERSION} \
    && ./bootstrap.sh \
    && ./configure \
       		--prefix=/usr \
       		--sysconfdir=/etc \
       		--mandir=/usr/share/man \
       		--localstatedir=/var \
       		--without-cpp \
       		--without-qt5 \
       		--without-erlang \
       		--without-nodejs \
       		--without-lua \
       		--without-python \
       		--without-py3 \
       		--without-perl \
       		--without-php \
       		--without-php_extension \
       		--without-dart \
       		--without-ruby \
       		--without-swift \
       		--without-rs \
       		--without-cl \
       		--without-haxe \
       		--without-netstd \
       		--without-d \
       		--disable-tutorial \
    && make \
    && make check \
    && make install \
    && cd /  \
    && rm -rf thrift-${THRIFT_VERSION}

FROM ubuntu:20.04
COPY --from=builder /usr/bin/thrift /usr/bin/thrift
