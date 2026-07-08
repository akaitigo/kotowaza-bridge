-- Enable trigram matching so the substring (LIKE '%term%') search on the
-- kotowaza text columns can use an index instead of a sequential scan.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- The search endpoint filters japanese, reading and meaning with a
-- leading-wildcard LIKE, which a standard btree index cannot serve. GIN trigram
-- indexes accelerate these substring matches on every branch of the OR.
CREATE INDEX IF NOT EXISTS idx_kotowaza_japanese_trgm ON kotowaza USING gin (japanese gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_kotowaza_reading_trgm ON kotowaza USING gin (reading gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_kotowaza_meaning_trgm ON kotowaza USING gin (meaning gin_trgm_ops);

-- Drop the btree index created in 001: it cannot serve leading-wildcard LIKE and
-- is now superseded by the trigram index on the same column.
DROP INDEX IF EXISTS idx_kotowaza_japanese;
