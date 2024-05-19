FROM alpine:latest


WORKDIR /app


COPY controller /app/controller


CMD ["./controller"]
