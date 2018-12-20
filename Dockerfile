# Build our go application
FROM golang:1.10 AS go-build
RUN go get github.com/internet-research-labs/fish-cake/server
WORKDIR /go/src/github.com/internet-research-labs/fish-cake/bin
ENV GOBIN=/go/bin
RUN go install -ldflags "-linkmode external -extldflags -static" -a run.go

# Compile our javascript
FROM node:11.2.0-slim AS node-build
COPY --from=go-build  /go/src/github.com/internet-research-labs/fish-cake/ /fish-cake
WORKDIR /fish-cake/front
RUN npm install
RUN npm run ayy

# Copy things from the build stage into the running application
FROM scratch
COPY --from=go-build /go/bin/run /run
COPY --from=node-build /fish-cake/static/ /static
EXPOSE 8080
CMD ["/run", "-static", "/static/"]
