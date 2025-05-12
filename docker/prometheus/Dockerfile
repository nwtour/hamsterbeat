FROM alpine:latest
RUN apk update
RUN apk add --no-cache bash
RUN apk add --no-cache prometheus
COPY prometheus/prometheus.yml /etc/prometheus.yml
CMD ["/usr/bin/prometheus", "--config.file=/etc/prometheus.yml"]
