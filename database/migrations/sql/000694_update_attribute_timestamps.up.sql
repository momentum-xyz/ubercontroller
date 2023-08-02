DO $$ BEGIN
ALTER TABLE node_attribute
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
EXCEPTION WHEN duplicate_column THEN
  RAISE NOTICE 'column already exists in node_attribute.';
END $$;

DO $$ BEGIN
ALTER TABLE object_attribute
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
EXCEPTION WHEN duplicate_column THEN
  RAISE NOTICE 'column already exists in object_attribute.';
END $$;

DO $$ BEGIN
ALTER TABLE user_attribute
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
EXCEPTION WHEN duplicate_column THEN
  RAISE NOTICE 'column already exists in user_attribute.';
END $$;

DO $$ BEGIN
ALTER TABLE object_user_attribute
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
EXCEPTION WHEN duplicate_column THEN
  RAISE NOTICE 'column already exists in object_user_attribute.';
END $$;

DO $$ BEGIN
ALTER TABLE user_user_attribute
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
EXCEPTION WHEN duplicate_column THEN
  RAISE NOTICE 'column already exists in user_user_attribute.';
END $$;