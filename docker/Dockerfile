FROM alpine:latest
LABEL Osiloke Emoekpere <me@osiloke.com>
RUN apk update && apk add --no-cache bash ca-certificates && rm -rf /var/cache/apk/* 
ADD mail /app/mail
RUN chmod +x /app/mail
# COPY docker-entrypoint.sh /  
WORKDIR /app 
# ENTRYPOINT ["/docker-entrypoint.sh"]
ENTRYPOINT ["/app/mail"]