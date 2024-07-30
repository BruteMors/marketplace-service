FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

ADD internal/repository/postgres/migrations/*.sql migrations/
ADD pgmigrator/lomsmigration.sh .

RUN chmod +x lomsmigration.sh

ENTRYPOINT ["bash", "lomsmigration.sh"]
