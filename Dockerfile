FROM debian
RUN apt update && \
    apt install -y curl && \
    curl -OL --fail https://github.com/db-journey/journey/releases/download/v2.1.1/journey.linux-amd64.tar.gz && \
    tar xvzf journey.linux-amd64.tar.gz && \
    mv build/journey.linux-amd64 /usr/local/bin/journey && \
    chmod +x /usr/local/bin/journey && \
    apt remove -y curl && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
ENTRYPOINT /usr/local/bin/journey
