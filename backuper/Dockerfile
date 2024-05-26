FROM alpine

WORKDIR /root/

COPY backup.sh /root/backup.sh

ARG BACKUPS_TIMEOUT

RUN chmod 755 backup.sh

RUN apk add --no-cache postgresql-client
RUN apk add --no-cache tini

RUN crontab -l | { cat; echo "0 */${BACKUPS_TIMEOUT} * * * sh /root/backup.sh"; } | crontab -

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/sbin/crond", "-f"]