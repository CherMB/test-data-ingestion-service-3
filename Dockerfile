FROM golang:1.22.4-alpine

LABEL maintainer=CB-Platform \
    email=engineering@cb-platform.com

RUN apk "upgrade" libssl3 libcrypto3

RUN  apk update \
  && apk upgrade \
  && apk add --update coreutils && rm -rf /var/cache/apk/*   \ 
  && apk add --update openjdk11 tzdata curl unzip bash \
  && apk add --no-cache nss \
  && rm -rf /var/cache/apk/*
ENV JAVA_HOME /usr/bin/java
ENV PATH $PATH:/usr/bin/java/bin  

COPY reports_service /

ENTRYPOINT ["/reports_service"]
