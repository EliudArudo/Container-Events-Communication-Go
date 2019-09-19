docker build -t eliudarudo/go-container-events-communication-consuming-frontend:dev -f  ../github.com/eliudarudo/consuming-frontend/Dockerfile  ../github.com/eliudarudo/consuming-frontend
docker build -t eliudarudo/go-container-events-communication-event-service:dev      -f  ../github.com/eliudarudo/event-service/Dockerfile       ../github.com/eliudarudo/event-service
docker build -t eliudarudo/go-container-events-communication-consuming-backend:dev  -f  ../github.com/eliudarudo/consuming-backend/Dockerfile   ../github.com/eliudarudo/consuming-backend
