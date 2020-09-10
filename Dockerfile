FROM alpine:3.12
VOLUME /data
VOLUME /conf.d
EXPOSE 80
COPY nano-run /bin/nano-run
COPY bundle/docker/server.yaml /server.yaml
CMD ["/bin/nano-run", "server", "run", "-f", "-c", "server.yaml"]