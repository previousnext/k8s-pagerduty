FROM golang:1.8
ADD workspace /go
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/linux/server/k8s-pagerduty /usr/local/bin/k8s-pagerduty
CMD ["k8s-pagerduty"]
