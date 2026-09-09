FROM ubuntu:24.04
LABEL name="gohub" author="Lyle Zheng"
COPY build/gohub /app/gohub
WORKDIR /app
ENTRYPOINT ["/app/gohub"]
