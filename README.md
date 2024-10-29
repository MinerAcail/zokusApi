    sudo docker exec -it zokusapi-zukudb-1 psql -U justzoku -d testzoku 

    //delect a table AND ALL RELATED TABLES WHEN USING CASCADE WHEN YOU HAVE Constraint fk_.....
    DROP TABLE TABLENAME, TABLENAME2 CASCADE;
    DROP TABLE cart_items, cart_item_groups CASCADE;
DROP TABLE carts, cart_items, cart_item_groups,delivery_meta_data,delivery_statuses,groups,orders,purchase_meta_data CASCADE;
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


only for GORM or psql

const { hasAccess:swapperHasAccess } = useProtect([ProtectType.swapper]); 



 Function Explain
GetAllFavoriteProduct => mostFavoritedProduct
To find the most favorited product, you can use SQL queries to count the number of times each product appears in the FavoriteProduct table and then select the product with the highest count.

