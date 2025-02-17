FROM ubuntu:latest
LABEL name="gohub" author="Lyle Zheng"
COPY gohub /app/gohub
WORKDIR /app
ENTRYPOINT ["/app/gohub"]