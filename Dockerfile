FROM golang:1.26-alpine3.23 AS build

RUN apk add git make \
    && apk cache clean

WORKDIR /src

# Install setup dependencies (e.g. betteralign).
COPY Makefile ./
RUN make setup

# Cache Go modules and local vendored dependencies needed for replacements.
COPY go.mod go.sum ./
RUN go mod download

COPY . /src

RUN make

FROM alpine:3.21

RUN apk add bash \
    && apk cache clean \
    && adduser -D -H appuser \
    && mkdir /app \
    && chown appuser /app

COPY --from=build /src/bin/client /usr/bin/client

USER appuser

WORKDIR /app

ENTRYPOINT ["client"]
