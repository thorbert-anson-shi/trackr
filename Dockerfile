FROM debian:13.4-slim@sha256:cedb1ef40439206b673ee8b33a46a03a0c9fa90bf3732f54704f99cb061d2c5a AS builder

WORKDIR /app

RUN apt-get update && apt-get -y install \
  ca-certificates=20250419 \
  wget=1.25.0-2 \
  --no-install-recommends

RUN wget --progress=dot:giga https://go.dev/dl/go1.26.2.linux-amd64.tar.gz && \
  tar -C /usr/local -xzf go1.26.2.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"
ENV GOPATH="/go"

COPY go.mod go.sum ./
RUN go mod download && go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN /usr/local/go/bin/go build -o trackr tobtoby/trackr


FROM debian:13.4-slim@sha256:cedb1ef40439206b673ee8b33a46a03a0c9fa90bf3732f54704f99cb061d2c5a AS runner

WORKDIR /app

COPY .env service-account-config.json ./
COPY postgresql/ ./postgresql/
COPY --from=builder /app/trackr ./
COPY --from=builder /go/bin/goose ./

CMD [ "./trackr" ]
