FROM golang:1.12.7-stretch

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update
RUN apt-get install \
    pv

ENTRYPOINT [ "/bin/bash" ]