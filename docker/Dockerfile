FROM alpine:3.21

COPY  build/linux/transform-mcp /transform-mcp
COPY config.toml .
COPY resources/security /security

EXPOSE 8080

CMD [ "/transform-mcp" ]
