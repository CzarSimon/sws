service_route="greeter.hello:81/greet"

while [ "1" = "1" ]
do
  curl $service_route
  sleep 0.2
done
