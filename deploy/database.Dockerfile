FROM postgres:16

RUN apt-get update && apt-get install -y curl \
    && curl https://install.citusdata.com/community/deb.sh | bash \
    && apt-get -y install postgresql-16-citus-12.1