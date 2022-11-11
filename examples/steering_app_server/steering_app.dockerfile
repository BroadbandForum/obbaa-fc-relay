FROM python:3.7-slim as builder

WORKDIR /install

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    gcc \
    g++ \
    make \
    python3-dev

RUN python -m venv env

COPY . .
RUN env/bin/pip install \
  -r requirements.txt 

RUN env/bin/python -m grpc_tools.protoc \
  -I. \
  --python_out=. \
  --grpc_python_out=. \
  ./control_relay_service.proto

RUN python -m grpc_tools.protoc \
  -I. \
  --python_out=. \
  --grpc_python_out=. \
  ./tr477_cpri_service.proto

RUN env/bin/python -m grpc_tools.protoc \
  -I. \
  --python_out=. \
  --grpc_python_out=. \
  ./tr477_cpri_message.proto


FROM python:3.7-slim as app

RUN apt-get update && \
    apt-get upgrade -y

WORKDIR /app

COPY --from=builder /install/ .

ENV STEERING_APP_PORT "50053"

CMD [ "env/bin/python", "./steering_app.py" ]