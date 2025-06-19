CREATE SUBSCRIPTION auth_subscription
CONNECTION 'host=auth_db port=5432 user=replicator password=replicator_password dbname=auth'
PUBLICATION auth_publication;
