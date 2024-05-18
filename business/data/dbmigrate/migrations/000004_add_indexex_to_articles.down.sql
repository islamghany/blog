-- Drop the GIN index on the tsvector column:
DROP INDEX IF EXISTS articles_tsv_document_idx;
-- Remove the trigger:
DROP TRIGGER IF EXISTS update_tsv_document_trigger ON articles;
-- Remove the function:
DROP FUNCTION IF EXISTS update_tsv_document();
-- Remove the tsvector column:
ALTER TABLE articles DROP COLUMN IF EXISTS tsv_document;