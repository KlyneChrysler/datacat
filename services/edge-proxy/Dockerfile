# Build context is the REPO ROOT (services import pkg/ via replace directives):
#   docker build -f services/edge-proxy/Dockerfile .
FROM golang:1.26 AS build
WORKDIR /src
COPY pkg/ pkg/
COPY services/edge-proxy/ services/edge-proxy/
WORKDIR /src/services/edge-proxy
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/edge-proxy ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/edge-proxy /edge-proxy
USER nonroot
ENTRYPOINT ["/edge-proxy"]
