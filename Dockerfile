FROM alpine
RUN apk add --update curl && \
    curl -OL https://github.com/db-journey/journey/releases/download/v2.0.0/journey-linux-amd64 && \
    mv journey-linux-amd64 /journey && \
    chmod +x /journey && \
    apk del curl && \
    rm -rf /var/cache/apk/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENTRYPOINT /journey
