## Sample Dockerfile to package the whole application (backend and frontend) as a Docker image.
# Sample build command:
# docker build --rm -t exter .

FROM node:13.6-alpine3.11 AS builder_fe
RUN apk add jq sed
RUN mkdir /build
COPY . /build
COPY ./fe-gui/src/config_one_image.json /build/fe-gui/src/config.json
RUN cd /build \
    && export APP_NAME=`jq -r '.name' appinfo.json` \
    && export APP_SHORTNAME=`jq -r '.shortname' appinfo.json` \
    && export APP_INITIAL=`jq -r '.initial' appinfo.json` \
    && export APP_VERSION=`jq -r '.version' appinfo.json` \
    && export APP_DESC=`jq -r '.description' appinfo.json` \
    # && export BUILD=`date +%Y%m%d%H%M` \
    && cd /build/fe-gui \
    && rm -rf dist node_modules \
    && sed -i s/#{pageTitle}/"$APP_NAME v$APP_VERSION"/g public/index.html \
    && sed -i s/#{appName}/"$APP_NAME"/g public/index.html \
    && sed -i s/#{appInitial}/"$APP_INITIAL"/g public/index.html \
    && sed -i s/#{appShortname}/"$APP_SHORTNAME"/g public/index.html \
    && sed -i s/#{appDescription}/"$APP_DESC"/g public/index.html \
    && sed -i s/#{appVersion}/"$APP_VERSION"/g public/index.html \
    && sed -i s/#{appName}/"$APP_NAME"/g src/config.json \
    && sed -i s/#{appInitial}/"$APP_INITIAL"/g src/config.json \
    && sed -i s/#{appShortname}/"$APP_SHORTNAME"/g src/config.json \
    && sed -i s/#{appDescription}/"$APP_DESC"/g src/config.json \
    && sed -i s/#{appVersion}/"$APP_VERSION"/g src/config.json \
    && npm install && npm run build
    # && sed -i 's/exter - Exter v0.1.0/GoVueAdmin v0.1.1 b'$BUILD/g public/index.html \

FROM golang:1.13-alpine AS builder_be
RUN apk add git build-base jq sed
RUN mkdir /build
COPY . /build
RUN cd /build \
    && export APP_NAME=`jq -r '.name' appinfo.json` \
    && export APP_SHORTNAME=`jq -r '.shortname' appinfo.json` \
    && export APP_INITIAL=`jq -r '.initial' appinfo.json` \
    && export APP_VERSION=`jq -r '.version' appinfo.json` \
    && export APP_DESC=`jq -r '.description' appinfo.json` \
    && cd /build/be-api \
    && sed -i s/#{appName}/"$APP_NAME"/g config/application.conf \
    && sed -i s/#{appInitial}/"$APP_INITIAL"/g config/application.conf \
    && sed -i s/#{appShortname}/"$APP_SHORTNAME"/g config/application.conf \
    && sed -i s/#{appDescription}/"$APP_DESC"/g config/application.conf \
    && sed -i s/#{appVersion}/"$APP_VERSION"/g config/application.conf \
    && go build -o main

FROM alpine:3.10
LABEL maintainer="Thanh Nguyen <btnguyen2k@gmail.com>"
RUN mkdir -p /app/frontend
COPY --from=builder_be /build/be-api/main /app/main
COPY --from=builder_be /build/be-api/config /app/config
COPY --from=builder_fe /build/ge-gui/dist /app/frontend
RUN apk add --no-cache -U tzdata bash ca-certificates \
    && update-ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime \
    && chmod 711 /app/main \
    && rm -rf /var/cache/apk/*
WORKDIR /app
EXPOSE 8000
CMD ["/app/main"]
#ENTRYPOINT /app/main
