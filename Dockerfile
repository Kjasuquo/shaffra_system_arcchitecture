FROM golang:1.22.3-alpine AS build

WORKDIR /app
COPY . .
RUN apk --update --no-cache add g++
RUN go mod tidy
RUN go build -o main cmd/main.go

# Now copy it into our base image.
FROM alpine:3.13
WORKDIR /app
COPY --from=build /app/main .

EXPOSE 5053
CMD ["/app/main"]
