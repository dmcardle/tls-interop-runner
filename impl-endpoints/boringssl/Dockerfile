# SPDX-FileCopyrightText: 2020 The tls-interop-runner Authors
# SPDX-License-Identifier: MIT

FROM golang:latest AS builder

RUN apt-get update && \
    apt-get install -y git cmake ninja-build perl && \
    #git clone --branch dc-spec-update https://github.com/xvzcf/boringssl /boringssl
    git clone https://boringssl.googlesource.com/boringssl /boringssl

WORKDIR /boringssl
#RUN git checkout 66929b751f97b3ce2868523bcdde549ea8324b56
RUN git fetch https://boringssl.googlesource.com/boringssl refs/changes/85/45285/27 && \
    git checkout -b change-45285 FETCH_HEAD

RUN mkdir /boringssl/build
WORKDIR /boringssl/build
RUN cmake .. -G"Ninja" && ninja

RUN mkdir /runner
WORKDIR /runner
COPY go.mod go.sum boringssl.go /runner/
RUN go build -o runner boringssl.go

FROM ubuntu:20.04

RUN apt-get update && \
    apt-get install -y net-tools tcpdump ethtool iproute2

COPY --from=builder /boringssl/build/tool/bssl /usr/bin/

COPY --from=builder /runner/runner /runner

ENTRYPOINT [ "/runner" ]
