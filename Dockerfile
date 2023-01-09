FROM node:18.12-slim AS web-builder
WORKDIR /import-tools/web
COPY web/. .
RUN npm i && \
    npm run build && \
    tar zcf dist.tgz dist/

FROM golang:1.16-alpine3.15 AS api-builder
WORKDIR /import-tools
COPY . .
COPY --from=web-builder /import-tools/web/dist.tgz serve/router/
RUN tar -C serve/router -zxf serve/router/dist.tgz && \
    go mod download && \
    CGO_ENABLED=0 go build -trimpath -o bin/import-tools main.go

FROM alpine:3.15
COPY --from=api-builder /import-tools/bin/import-tools /usr/local/bin/
EXPOSE 5000
ENTRYPOINT ["/usr/local/bin/import-tools"]
