FROM golang:latest AS build
ADD . /src
RUN cd /src && go build -o /go/bin/digitalocean-dyndns

FROM gcr.io/distroless/base:debug
COPY --from=build /go/bin/digitalocean-dyndns /
ENTRYPOINT ["/digitalocean-dyndns"]
