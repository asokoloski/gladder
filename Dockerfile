FROM busybox:ubuntu-14.04

MAINTAINER Aaron Sokoloski <asokoloski@gmail.com>
ENV GLADDER_HTTP_ADDR :80
EXPOSE 80
ENV GLADDER_DB /data/gladder-db.gkv
VOLUME /data
ADD gladder /usr/bin/gladder
ENTRYPOINT ["/usr/bin/gladder"]

