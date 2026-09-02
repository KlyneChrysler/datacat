# Build context is the repo root.
FROM golang:1.26 AS build
WORKDIR /src
COPY pkg/ pkg/
COPY services/lite/ services/lite/
WORKDIR /src/services/lite
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/datacat-lite ./cmd/lite

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/datacat-lite /datacat-lite
USER nonroot
ENTRYPOINT ["/datacat-lite"]
