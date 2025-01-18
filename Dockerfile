FROM golang:1.23-bullseye AS backend-builder

WORKDIR /app/backend

COPY backend/ .
ARG CGO_ENABLED=1
RUN go mod download
RUN go build -a -installsuffix cgo -o backend server.go

FROM node:23-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/ .
RUN npm install --production
RUN npm install npm-run-all
RUN npm run build

# todo: change base image for a lighter one
# currently alpine doesn't have some necessary libs to run the go executable
FROM golang:1.23-bullseye

# RUN apk add --no-cache bash curl ca-certificates tini nginx
RUN apt update
RUN apt install -y tini nginx cron libc6

COPY nginx.conf /etc/nginx/nginx.conf

WORKDIR /app

COPY --from=backend-builder /app/backend/backend ./backend/backend
COPY --from=backend-builder /app/backend/config/ ./backend/config/
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

COPY cron.sh /app/cron.sh
RUN chmod +x /app/cron.sh

COPY init.sh /app/init.sh
RUN chmod +x /app/init.sh

# Collect data every 6 hours
RUN echo "0 */6 * * * /app/cron.sh" > /etc/cron.d/app-cron

# backend
EXPOSE 1323 
# frontend
EXPOSE 80  

# Use tini for better process management
ENTRYPOINT ["/usr/bin/tini", "--"]

CMD ["sh", "-c", "\
    service cron start && \
    ./backend/backend & \
    nginx -g 'daemon off;' & \
    /app/init.sh"]

