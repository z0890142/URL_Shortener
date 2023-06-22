# Copyright (c) 2020 Wisers DMA Team. All rights reserved.

FROM --platform=linux/amd64 mysql:5.7

COPY . /docker-entrypoint-initdb.d

RUN chmod -R 0755 /docker-entrypoint-initdb.d
