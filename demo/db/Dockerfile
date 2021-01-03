FROM postgres:11

ARG TESTDATA_URL='https://sp.postgresqltutorial.com/wp-content/uploads/2019/05/dvdrental.zip'

ENV PGPASSWORD=$POSTGRES_PASSWORD

WORKDIR /var/dvdrental

RUN apt-get update && \
    apt-get install -y curl unzip

RUN curl -o dvdrental.zip ${TESTDATA_URL} && \
    unzip dvdrental.zip

COPY restore.sh .

CMD [ "/bin/sh", "restore.sh", "dvdrental.tar" ]
    
