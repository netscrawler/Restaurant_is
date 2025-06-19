CREATE SUBSCRIPTION user_subscription
CONNECTION 'host=user_db port=5439 user=replicator password=replicator_password dbname=user'
PUBLICATION user_publication;
