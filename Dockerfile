FROM scratch
MAINTAINER Julien Bordellier <julien.bordellier@etixgroup.com>

COPY pepersalt-api pepersalt-api
ENV PORT 3000
EXPOSE 3000
ENTRYPOINT ["/pepersalt-api"]
