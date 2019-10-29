FROM alpine:3.8

RUN adduser -D ptoe
RUN apk --update add ca-certificates
USER ptoe

ADD ./bin/beaujr/plex2emby-linux_amd64 /usr/local/bin/plex2emby-linux_amd64

ENTRYPOINT ["/usr/local/bin/plex2emby-linux_amd64"]
ARG VCS_REF
LABEL org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/beaujr/plex2emby" \
      org.label-schema.license="Apache-2.0"