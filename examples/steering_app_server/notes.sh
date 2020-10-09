
#other deps
apt install libxml2-dev libxslt1-dev
#python 3
pip install -r requirements.txt

#compile protobuffs
python -m grpc_tools.protoc \
  -I. \
  --python_out=. \
  --grpc_python_out=. \
  ./control_relay_service.proto

#docker image
docker build . \
  -t obbaa/steering_app_server \
  -f ./steering_app.dockerfile

#run steering app
docker run -it \
   --name STEERING_APP_SERVER  \
   --rm \
   --env OBBAA_ADDRESS=192.168.56.102 \
   --env OBBAA_PORT=9292 \
   obbaa/steering_app_server