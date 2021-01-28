FROM golang:1.12.7-alpine3.10 AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/RGuide-backend
COPY . .
ENV GOPATH="/go/src/RGuide-backend"
RUN go get -d ./...
RUN GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/RGuide-backend/main .
EXPOSE 8080
ENTRYPOINT  ["./main"]