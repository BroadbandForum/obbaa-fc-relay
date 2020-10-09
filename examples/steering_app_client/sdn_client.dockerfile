FROM golang:1.14.4

RUN apt-get update && \
    apt-get upgrade -y

WORKDIR /sdn_client   

COPY bin/. ./bin/

RUN mkdir data && \
    cd data

COPY data/. ./data/

ENV STEERING_APP_ADDRESS_CLIENT "0.0.0.0"
ENV STEERING_APP_PORT_CLIENT "50055"

CMD ["./bin/sdn_client"]   