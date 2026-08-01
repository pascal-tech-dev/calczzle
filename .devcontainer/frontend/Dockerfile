FROM node:24-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /workspace/frontend

RUN chown -R node:node /workspace/frontend

USER node

CMD ["sleep", "infinity"]