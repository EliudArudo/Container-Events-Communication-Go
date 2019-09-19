export NEW_VERSION=v0.1

docker build -t eliudarudo/go-container-events-communication-consuming-frontend:dev -t eliudarudo/go-container-events-communication-consuming-frontend:$NEW_VERSION -f ../github.com/eliudarudo/consuming-frontend/Dockerfile ../github.com/eliudarudo/consuming-frontend
docker build -t eliudarudo/go-container-events-communication-event-service:dev      -t eliudarudo/go-container-events-communication-event-service:$NEW_VERSION      -f ../github.com/eliudarudo/event-service/Dockerfile      ../github.com/eliudarudo/event-service
docker build -t eliudarudo/go-container-events-communication-consuming-backend:dev  -t eliudarudo/go-container-events-communication-consuming-backend:$NEW_VERSION  -f ../github.com/eliudarudo/consuming-backend/Dockerfile  ../github.com/eliudarudo/consuming-backend

docker push eliudarudo/go-container-events-communication-consuming-frontend:dev &&  docker push eliudarudo/go-container-events-communication-consuming-frontend:$NEW_VERSION
docker push eliudarudo/go-container-events-communication-event-service:dev      &&  docker push eliudarudo/go-container-events-communication-event-service:$NEW_VERSION
docker push eliudarudo/go-container-events-communication-consuming-backend:dev  &&  docker push eliudarudo/go-container-events-communication-consuming-backend:$NEW_VERSION