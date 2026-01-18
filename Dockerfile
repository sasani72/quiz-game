FROM ubuntu:latest
LABEL authors="arvan"

ENTRYPOINT ["top", "-b"]