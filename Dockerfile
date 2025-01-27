FROM golang:1.23-bullseye AS backend-builder

WORKDIR /backend

COPY backend/ .
ARG CGO_ENABLED=1
RUN go mod download
RUN go build -a -installsuffix cgo -o backend server.go

FROM node:23-alpine AS frontend-builder

WORKDIR /frontend

COPY frontend/ .

# empty env var so the frontend uses the browser's url (fallback)
RUN sed -i 's|VITE_BACKEND_URL=.*|VITE_BACKEND_URL=|' .env

RUN npm install --production
RUN npm install npm-run-all
RUN npm run build

# todo: change base image for a lighter one
# currently alpine doesn't have some necessary libs to run the go executable
FROM golang:1.23-bullseye

# RUN apk add --no-cache bash curl ca-certificates tini nginx
RUN apt update
RUN apt install -y tini nginx cron libc6 sqlite3

COPY nginx.conf /etc/nginx/nginx.conf

WORKDIR /

COPY --from=backend-builder /backend/backend ./backend/backend
COPY --from=backend-builder /backend/config_sample/ ./backend/config_sample/
COPY --from=frontend-builder /frontend/dist ./frontend/dist

COPY cron.sh /cron.sh
RUN chmod +x /cron.sh

COPY init.sh /init.sh
RUN chmod +x /init.sh

# Collect data every 6 hours
ADD crontab /etc/cron.d/cron
RUN chmod 0644 /etc/cron.d/cron
RUN touch /var/log/cron.log

# backend
EXPOSE 1323 
# frontend
EXPOSE 80  

# Use tini for better process management
ENTRYPOINT ["/usr/bin/tini", "--"]

CMD ["sh", "-c", "\
    cron && \
    cd backend && ./backend & \
    nginx -g 'daemon off;' & \
    cd / && /init.sh && \
    wait"]

