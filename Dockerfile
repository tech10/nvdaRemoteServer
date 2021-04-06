FROM golang:alpine as build

RUN apk add --no-cache gcc musl-dev
RUN mkdir /app
COPY . /app/
RUN cd /app && go build -buildmode=pie "-asmflags=all='-trimpath=`pwd`'" -ldflags '-w -s -linkmode external -extldflags "-static"' -o nvdaRemoteServer .

FROM scratch

COPY --from=build /app/nvdaRemoteServer /nvdaRemoteServer

EXPOSE 6837
CMD ["/nvdaRemoteServer"]
