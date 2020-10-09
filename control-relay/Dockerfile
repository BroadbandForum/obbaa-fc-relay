FROM ubuntu:18.04

RUN apt-get update && \
    apt-get upgrade -y 
    
ADD dist/control-relay.tgz /control_relay

WORKDIR /control_relay

CMD ["./control-relay"]    