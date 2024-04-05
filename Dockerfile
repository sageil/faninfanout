
FROM golang:1.22.2-alpine as builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -o ./entrypoint ./cmd/main.go

FROM alpine:3.17.2 as runner
ARG VERSION=1.0
ENV VERSION=$VERSION
ARG ACCOUNT=nonroot
ARG PORT=3000
ENV PORT=$PORT
RUN addgroup -S ${ACCOUNT} && adduser -S ${ACCOUNT} -G ${ACCOUNT}
WORKDIR /app
COPY --from=builder app/entrypoint .
RUN chown -R ${ACCOUNT}:${ACCOUNT} /app
EXPOSE ${PORT}
USER ${ACCOUNT}:${ACCOUNT}
ENTRYPOINT ["./entrypoint"]


