## MCP Transformer

This repo has the go service that will transform the payload from the gateway to the relevant OpenAPI format based on the input schema.

You need Go 1.23.x to build this.


### Build
```
make build
```

### Build docker image

```docker
docker build -t go-mcp:0.1 . -f docker/Dockerfile
```