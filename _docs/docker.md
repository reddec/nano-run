# Docker


Check images in [releases](https://github.com/reddec/nano-run/releases)

* Latest one: `reddec/nano-run:latest`


Create Dockerfile inherited from the image and copy configuration and binaries

## Minimal example

**app.yaml**
```yaml
command: '/mybinary --with --some args'
```

**Dockerfile**
```dockerfile
FROM reddec/nano-run
COPY app.yaml /conf.d/app.yaml
COPY mybinary /mybinary
```

**Build & Run**

```bash
docker run --rm -p 127.0.0.1:8080:80 $(docker build -q .)
```

Check it's working by sending test request

```
curl -v -X POST "http://127.0.0.1:8080/app/"
```

* To keep tasks persistent - mount `/data` volume like:
`docker run -v $(pwd)/data:data ...`