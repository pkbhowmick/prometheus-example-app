FROM alpine:latest

COPY prometheus-example-app .

ENTRYPOINT ["./prometheus-example-app"]

