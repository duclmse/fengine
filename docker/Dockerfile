FROM golang:1.18-alpine AS builder

ARG TZ
ARG SVC
ARG GOARCH
ARG GOARM
ARG VERSION
ARG COMMIT
ARG TIME

WORKDIR /viot
COPY . .
RUN apk update && apk add make tree \
 && make $SVC \
 && tree -L 3 \
 && mv build/viot-$SVC /exe


###########################################################
FROM scratch

ENV TZ=$TZ

# Certificates are needed so that mailing util can work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /

ENTRYPOINT ["/exe"]
