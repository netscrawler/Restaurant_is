CREATE SUBSCRIPTION menu_subscription
CONNECTION 'host=menu_db port=5434 user=replicator password=replicator_password dbname=menu'
PUBLICATION menu_publication;
