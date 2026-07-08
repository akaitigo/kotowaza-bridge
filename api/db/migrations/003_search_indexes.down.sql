-- Reverse 003_search_indexes.up.sql.
DROP INDEX IF EXISTS idx_kotowaza_japanese_trgm;
DROP INDEX IF EXISTS idx_kotowaza_reading_trgm;
DROP INDEX IF EXISTS idx_kotowaza_meaning_trgm;

-- Restore the btree index that migration 001 created on japanese.
CREATE INDEX IF NOT EXISTS idx_kotowaza_japanese ON kotowaza (japanese);

-- The pg_trgm extension is intentionally left installed: dropping an extension is
-- not always reversible (other objects may depend on it) and re-enabling it via
-- the up migration is cheap and idempotent.
