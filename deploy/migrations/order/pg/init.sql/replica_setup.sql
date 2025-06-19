CREATE SUBSCRIPTION order_subscription
CONNECTION 'host=order_db port=5435 user=replicator password=replicator_password dbname=order'
PUBLICATION order_publication;
