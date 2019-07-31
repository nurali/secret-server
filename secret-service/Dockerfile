FROM alpine:3.10

RUN mkdir /app
WORKDIR /app
ADD ./.out/secret-service /app/secret-service

CMD ["./secret-service"]
