FROM ubuntu:latest
RUN apt-get update; apt-get install -y fortune
COPY fortuneloop.sh /bin
ENTRYPOINT ["/bin/fortuneloop.sh"]
CMD ["10"]
