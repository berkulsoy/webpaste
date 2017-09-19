FROM alpine:latest
ADD webpaste-linux-amd64 /
EXPOSE 8290
WORKDIR /
CMD /webpaste-linux-amd64