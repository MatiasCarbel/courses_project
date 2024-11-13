FROM solr:latest

COPY config/ /config/

USER root
RUN chown -R solr:solr /config

USER solr

CMD ["solr-precreate", "courses", "/config"] 