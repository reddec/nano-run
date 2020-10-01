FROM alpine:3.12
ENV GIN_MODE=release
VOLUME /data
VOLUME /conf.d
EXPOSE 80
COPY templates /ui
COPY nano-run /bin/nano-run
COPY bundle/docker/server.yaml /server.yaml
CMD ["/bin/nano-run", "server", "run", "-f", "-c", "server.yaml"]