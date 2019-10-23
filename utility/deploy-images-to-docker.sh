export COLLECTIVE_VERSION=v1.0

docker build -t eliudarudo/go-events-communication-consuming-frontend:$COLLECTIVE_VERSION -t eliudarudo/go-events-communication-consuming-frontend:latest -f ../src/github.com/eliudarudo/consuming-frontend/Dockerfile ../src/github.com/eliudarudo/consuming-frontend
docker build -t eliudarudo/go-events-communication-events-service:$COLLECTIVE_VERSION -t eliudarudo/go-events-communication-events-service:latest -f ../src/github.com/eliudarudo/events-service/Dockerfile ../src/github.com/eliudarudo/events-service
docker build -t eliudarudo/go-events-communication-consuming-backend:$COLLECTIVE_VERSION -t eliudarudo/go-events-communication-consuming-backend:latest -f ../src/github.com/eliudarudo/consuming-backend/Dockerfile ../src/github.com/eliudarudo/consuming-backend


docker push eliudarudo/go-events-communication-consuming-frontend:$COLLECTIVE_VERSION
docker push eliudarudo/go-events-communication-events-service:$COLLECTIVE_VERSION
docker push eliudarudo/go-events-communication-consuming-backend:$COLLECTIVE_VERSION

docker push eliudarudo/go-events-communication-consuming-frontend:latest
docker push eliudarudo/go-events-communication-events-service:latest
docker push eliudarudo/go-events-communication-consuming-backend:latest