# Build context is the REPO ROOT (services import pkg/ via replace directives):
#   docker build -f services/enforcement/Dockerfile .
FROM golang:1.26 AS build
WORKDIR /src
COPY pkg/ pkg/
COPY services/enforcement/ services/enforcement/
WORKDIR /src/services/enforcement
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/enforcement ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/enforcement /enforcement
USER nonroot
ENTRYPOINT ["/enforcement"]
