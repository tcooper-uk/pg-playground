CREATE SUBSCRIPTION dvdrental_sub
    CONNECTION 'host=dvdrental-primary port=5432 dbname=dvdrental user=replicator password=replicator' 
	PUBLICATION dvdrental_pub;