# BUILDER
FROM golang:1.25 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go env && go mod download -x

COPY . .

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags="-s -w" -o /app/guestbook .

# FINAL
FROM scratch AS final

# environment variables
ENV GB_TITLE=Guestbook
ENV GB_FOOTER=
ENV GB_ENTRIES_PER_PAGE=10
ENV GB_PAPER_CSS=false
ENV GB_USE_RATELIMIT=true
ENV GB_RATELIMIT=0.2
ENV GB_BURSTLIMIT=2
ENV PORT=8080

COPY --from=builder /app/guestbook /guestbook
COPY --from=builder /src/templates /templates
COPY --from=builder /src/static /static

VOLUME [ "/data" ]

EXPOSE ${PORT}/tcp

CMD ["/guestbook"]
