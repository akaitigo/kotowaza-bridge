-- Create kotowaza table
CREATE TABLE IF NOT EXISTS kotowaza (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    japanese TEXT NOT NULL,
    reading TEXT NOT NULL,
    meaning TEXT NOT NULL,
    origin TEXT NOT NULL DEFAULT '',
    usage_example TEXT NOT NULL DEFAULT '',
    cultural_note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create equivalent table
CREATE TABLE IF NOT EXISTS equivalent (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kotowaza_id UUID NOT NULL REFERENCES kotowaza(id) ON DELETE CASCADE,
    language TEXT NOT NULL,
    expression TEXT NOT NULL,
    literal_meaning TEXT NOT NULL DEFAULT '',
    explanation TEXT NOT NULL DEFAULT '',
    UNIQUE(kotowaza_id, language, expression)
);

CREATE INDEX IF NOT EXISTS idx_equivalent_kotowaza_id ON equivalent(kotowaza_id);
CREATE INDEX IF NOT EXISTS idx_kotowaza_japanese ON kotowaza(japanese);
