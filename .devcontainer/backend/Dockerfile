FROM golang:1.26-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        git \
        make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace/backend

CMD ["sleep", "infinity"]