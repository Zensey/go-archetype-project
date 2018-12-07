#FROM go-archetype-project
FROM golang:1.11.2-stretch

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
ADD . /app/
WORKDIR /app

RUN set -eux; \
    apt-get update; \
	apt-get install -y postgresql-9.6 gosu nano redis-server screen rsyslog tmux; \
	useradd -r docker; \
	cat /app/redis.conf.add >> /etc/redis/redis.conf; \
	echo "host all  all    172.17.0.0/24  md5" >> /etc/postgresql/9.6/main/pg_hba.conf; \
	echo "host all  all    127.0.0.1/32  md5" >> /etc/postgresql/9.6/main/pg_hba.conf; \
	echo "host all  all    ::1/128       md5" >> /etc/postgresql/9.6/main/pg_hba.conf; \
	echo "listen_addresses='*'" >> /etc/postgresql/9.6/main/postgresql.conf; \
	sed -i -e "s/'\(.*\)'$/'\1 -w'/" /etc/postgresql/9.6/main/pg_ctl.conf; \
	sed -i -e "s/# en_US.UTF-8/en_US.UTF-8/" /etc/locale.gen; \
	cp /app/task.rsyslog /etc/rsyslog.d; \
	sed -i -e "/IncludeConfig.*conf/ a \$IncludeConfig /etc/rsyslog.d/*.rsyslog" /etc/rsyslog.conf; \
	chmod +x /app/docker_entrypoint.sh; \
	locale-gen; \
	echo "PATH=\$PATH:\$GOPATH/bin:/usr/local/go/bin" >> /etc/profile; \
	chmod +x /app/tmux.sh; \
	make all \
	;


#VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
#USER root
#RUN make api
#CMD ["/app/build/_workspace/bin/api"]
#RUN "/bin/bash"

ENTRYPOINT ["/app/docker_entrypoint.sh"]
EXPOSE 5432


