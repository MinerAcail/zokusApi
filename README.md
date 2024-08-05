sudo docker exec -it zokusapi-zukudb-1 psql -U justzoku -d testzoku 

delect a table 
DROP TABLE users ;
drop all table
DO $$ 
DECLARE 
    r RECORD;
BEGIN 
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP 
        EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
END $$;



only for gqlgen
go run github.com/99designs/gqlgen generate
go run github.com/99designs/gqlgen generate


