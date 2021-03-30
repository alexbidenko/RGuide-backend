FROM golang:1.15-alpine AS build
# Support CGO and SSL
RUN apk --no-cache add gcc g++ make
RUN apk add git
WORKDIR /go/src/RGuide-backend
COPY . .
ENV GOPATH="/go/src/RGuide-backend"
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN go get github.com/asaskevich/govalidator
RUN go get gorm.io/gorm
RUN go get gorm.io/driver/postgres
RUN GOOS=linux go build -ldflags="-s -w" -o main .

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/RGuide-backend/main .
EXPOSE 8010
ENTRYPOINT  ["./main"]