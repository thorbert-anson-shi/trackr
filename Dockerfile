FROM dhi.io/golang:1@sha256:a473fd3819e71b8ec2d81b38014892c937c66b06d94f5a5d49ad393f0cf2e6a2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN /usr/local/go/bin/go build -o trackr -ldflags="-s -w" tobtoby/trackr


FROM dhi.io/alpine-base:3.23@sha256:1def3ff29647c43c52f041c378110d513c57c9a5346bec75728205f7bd7e4fe8 AS runner

WORKDIR /app

COPY postgresql/ ./postgresql/
COPY --from=builder /app/trackr ./

CMD [ "./trackr" ]
