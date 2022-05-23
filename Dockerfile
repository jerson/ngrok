FROM golang:alpine

RUN apk add --no-cache git findutils build-base

WORKDIR /app/pgrok

RUN mkdir -p build/

COPY . .

RUN go mod download

ENTRYPOINT [ "scripts/build-all.sh" ]