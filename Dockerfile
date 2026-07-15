FROM dhi.io/golang:1.26-alpine3.24-dev@sha256:e48a91483983467f426cae8656aa16be252c6f2e290125e10db01259352a54ca AS builder

WORKDIR /app

ENV GOCACHE=/go-build
ENV GOMODCACHE=/go-mod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go-mod-cache go mod download

COPY . .
RUN --mount=type=cache,target=/go-mod-cache \
  --mount=type=cache,target=/go-build \
  go build -o trackr -ldflags="-s -w" tobtoby/trackr


FROM dhi.io/alpine-base:3.23@sha256:1def3ff29647c43c52f041c378110d513c57c9a5346bec75728205f7bd7e4fe8 AS runner

WORKDIR /app

COPY postgresql/ ./postgresql/
COPY .well-known/ .well-known/
COPY --from=builder /app/trackr ./

CMD [ "./trackr" ]
