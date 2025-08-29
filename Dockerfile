# BUILDER
FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go env && go mod download -x

COPY . .

RUN ls config

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags="-s -w" -o /app/guestbook .

# FINAL
FROM scratch AS final

COPY --from=builder /app/guestbook /guestbook
COPY --from=builder /src/templates /templates
COPY --from=builder /src/static /static
COPY --from=builder /src/config /config

VOLUME [ "/data" ]

EXPOSE 8080

# default environment variables
ENV PORT=8080

CMD ["/guestbook"]
