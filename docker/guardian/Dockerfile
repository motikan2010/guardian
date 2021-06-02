FROM golang

ADD workdir .

RUN sed -i 's,POSTGRES_HOST,db,g' appsettings.live.json && \
    sed -i 's,POSTGRES_USER,guardian,g' appsettings.live.json && \
    sed -i 's,POSTGRES_DB,guardiandb,g' appsettings.live.json

CMD ["./guardian"]

EXPOSE 80 443