# Build context is the REPO ROOT (services import pkg/ via replace directives):
#   docker build -f services/traffic-sim/Dockerfile .
FROM golang:1.26 AS build
WORKDIR /src
COPY pkg/ pkg/
COPY services/traffic-sim/ services/traffic-sim/
WORKDIR /src/services/traffic-sim
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/traffic-sim ./cmd/sim

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/traffic-sim /traffic-sim
USER nonroot
ENTRYPOINT ["/traffic-sim"]
