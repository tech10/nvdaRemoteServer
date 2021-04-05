FROM golang:alpine as build

RUN mkdir /app
COPY . /app/
RUN cd /app && go build -o nvdaRemoteServer .

FROM alpine

COPY --from=build /app/nvdaRemoteServer /bin/nvdaRemoteServer

EXPOSE 6837 6837
CMD ["/bin/nvdaRemoteServer"]
