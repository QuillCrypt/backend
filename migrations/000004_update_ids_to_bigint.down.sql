-- Drop FK constraints
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_receiver_id_fkey;
ALTER TABLE messages_delete DROP CONSTRAINT IF EXISTS messages_delete_sender_id_fkey;
ALTER TABLE messages_delete DROP CONSTRAINT IF EXISTS messages_delete_receiver_id_fkey;

-- Revert users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS github_id TEXT UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS google_id TEXT UNIQUE;

ALTER TABLE users ALTER COLUMN id TYPE UUID USING gen_random_uuid();
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- Revert messages table
ALTER TABLE messages ALTER COLUMN sender_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages ALTER COLUMN receiver_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages ADD CONSTRAINT messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages ADD CONSTRAINT messages_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE;

-- Revert messages_delete table
ALTER TABLE messages_delete ALTER COLUMN sender_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages_delete ALTER COLUMN receiver_id TYPE UUID USING gen_random_uuid();
ALTER TABLE messages_delete ADD CONSTRAINT messages_delete_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages_delete ADD CONSTRAINT messages_delete_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
DROP INDEX IF EXISTS idx_messages_receiver_id;
DROP INDEX IF EXISTS idx_messages_delete_receiver_id;
CREATE INDEX idx_messages_receiver_id ON messages (receiver_id);
CREATE INDEX idx_messages_delete_receiver_id ON messages_delete (receiver_id);
