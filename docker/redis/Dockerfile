FROM redis
COPY redis/redis.conf /etc/redis.conf
RUN ["chmod", "a+r", "/etc/redis.conf"]
CMD redis-server /etc/redis.conf
