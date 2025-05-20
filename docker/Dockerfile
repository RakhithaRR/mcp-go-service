FROM alpine:3.21

ARG MCP_USER=wso2
ARG MCP_USER_ID=10500
ARG MCP_USER_GROUP=wso2
ARG MCP_USER_GROUP_ID=10500
ARG MCP_USER_HOME=/home/${MCP_USER}

RUN \
    addgroup -S -g ${MCP_USER_GROUP_ID} ${MCP_USER_GROUP} \
    && adduser -S -u ${MCP_USER_ID} -h ${MCP_USER_HOME} -G ${MCP_USER_GROUP} ${MCP_USER}

WORKDIR ${MCP_USER_HOME}

COPY  build/linux/transform-mcp ./transform-mcp
COPY config.toml ./config.toml
COPY resources/security ./resources/security

RUN chown -R ${MCP_USER}:${MCP_USER_GROUP} ${MCP_USER_HOME}

USER ${MCP_USER}

EXPOSE 8080

CMD [ "./transform-mcp" ]
