FROM node:current-alpine AS build-frontend
LABEL maintainer Ascensio System SIA <support@onlyoffice.com>
ARG BACKEND_GATEWAY
ARG BACKEND_GATEWAY_WS
ARG DOC_SERVER
ARG WORD_FILE
ARG SLIDE_FILE
ARG SPREADSHEET_FILE
ENV BACKEND_GATEWAY=$BACKEND_GATEWAY \
    BACKEND_GATEWAY_WS=$BACKEND_GATEWAY_WS \
    DOC_SERVER=$DOC_SERVER \
    WORD_FILE=$WORD_FILE \
    SLIDE_FILE=$SLIDE_FILE \
    SPREADSHEET_FILE=$SPREADSHEET_FILE
WORKDIR /usr/src/app
COPY ./frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

FROM golang:alpine AS build-gateway
WORKDIR /usr/src/app
COPY backend .
RUN go build services/gateway/main.go

FROM golang:alpine AS build-auth
WORKDIR /usr/src/app
COPY backend .
RUN go build services/auth/main.go

FROM golang:alpine AS build-builder
WORKDIR /usr/src/app
COPY backend .
RUN go build services/builder/main.go

FROM golang:alpine AS build-callback
WORKDIR /usr/src/app
COPY backend .
RUN go build services/callback/main.go

FROM golang:alpine AS gateway
WORKDIR /usr/src/app
COPY --from=build-gateway \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 4044
CMD ["./main", "server"]

FROM golang:alpine AS auth
WORKDIR /usr/src/app
COPY --from=build-auth \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 5051
CMD ["./main", "server"]

FROM golang:alpine AS builder
WORKDIR /usr/src/app
COPY --from=build-builder \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 6060
CMD ["./main", "server"]

FROM golang:alpine AS callback
WORKDIR /usr/src/app
COPY --from=build-callback \
     /usr/src/app/main \
     /usr/src/app/main
EXPOSE 5044
CMD ["./main", "server"]

FROM nginx:alpine AS frontend
COPY --from=build-frontend \
    /usr/src/app/build \
    /usr/share/nginx/html
EXPOSE 80
