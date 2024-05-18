-- There are two ways to handling full-text search for that query
-- select title, aid from (
--  select articles.id as aid, articles.title as title, setweight(to_tsvector(articles.title), 'A') || 
--  setweight(to_tsvector(articles.content), 'B') || 
--  setweight(to_tsvector(users.username), 'C') || 
--  setweight(to_tsvector(array_to_string(articles.tags, ', ', '*')), 'B') as document 
--  from articles 
--  join users on users.id = articles.author_id 
--  group by articles.id, users.id) a_search
--  where a_search.document @@ to_tsquery('islam')
--  order by ts_rank(a_search.document , to_tsquery('islam')) DESC;
-- 1- using a materialized view and update it periodically
-- cons of this approach is that it will take time to update the materialized view
-- 2- using a trigger to update the tsvector column on insert and update
--add the tsvector column to the articles table:
ALTER TABLE articles
ADD COLUMN tsv_document tsvector;
-- create a function and a trigger to update this column whenever necessary:
CREATE OR REPLACE FUNCTION update_tsv_document() RETURNS TRIGGER AS $$ BEGIN NEW.tsv_document := setweight(to_tsvector('english', NEW.title), 'A') || setweight(to_tsvector('english', NEW.content), 'B') || setweight(
        to_tsvector('english', array_to_string(NEW.tags, ', ', '*')),
        'B'
    );
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER update_tsv_document_trigger BEFORE
INSERT
    OR
UPDATE ON articles FOR EACH ROW EXECUTE FUNCTION update_tsv_document();
-- Ensure the initial population of the tsvector column:
UPDATE articles
SET tsv_document = setweight(to_tsvector('english', title), 'A') || setweight(to_tsvector('english', content), 'B') || setweight(
        to_tsvector('english', array_to_string(tags, ', ', '*')),
        'B'
    );
-- Create the GIN index on the tsvector column:
CREATE INDEX articles_tsv_document_idx ON articles USING GIN(tsv_document);