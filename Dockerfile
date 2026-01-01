FROM  debian:trixie-slim
LABEL maintainer="Rekey <rekey@me.com>"

WORKDIR /app/
ENV TZ=Asia/Shanghai
ADD ./bin/club /app/club
ADD ./config /app/config

RUN apt update && apt install -y tzdata openssl ca-certificates && update-ca-certificates

VOLUME /app/data
VOLUME /app/config
EXPOSE 8888

CMD ["/app/club"]