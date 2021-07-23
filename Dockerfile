FROM golang:1.15-alpine AS build-env

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /  

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./out/envoy-authz

COPY --from=build /out/authz /envoy-authz

WORKDIR "/authz"
EXPOSE 9000 

CMD ["./envoy-authz"]