docker build -t eliudarudo/go-events-communication-consuming-frontend:dev -f ../src/github.com/eliudarudo/consuming-frontend/Dockerfile ../src/github.com/eliudarudo/consuming-frontend
docker build -t eliudarudo/go-events-communication-events-service:dev -f ../src/github.com/eliudarudo/events-service/Dockerfile ../src/github.com/eliudarudo/events-service
docker build -t eliudarudo/go-events-communication-consuming-backend:dev -f ../src/github.com/eliudarudo/consuming-backend/Dockerfile ../src/github.com/eliudarudo/consuming-backend

