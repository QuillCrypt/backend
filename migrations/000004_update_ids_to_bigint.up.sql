-- Remove social IDs
ALTER TABLE users DROP COLUMN IF EXISTS github_id;
ALTER TABLE users DROP COLUMN IF EXISTS google_id;

-- Drop FK constraints to allow type change
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_id_fkey;
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_receiver_id_fkey;
ALTER TABLE messages_delete DROP CONSTRAINT IF EXISTS messages_delete_sender_id_fkey;
ALTER TABLE messages_delete DROP CONSTRAINT IF EXISTS messages_delete_receiver_id_fkey;

-- Update users table
ALTER TABLE users ALTER COLUMN id DROP DEFAULT;
ALTER TABLE users ALTER COLUMN id TYPE BIGINT USING 0;

-- Update messages table
ALTER TABLE messages ALTER COLUMN sender_id TYPE BIGINT USING 0;
ALTER TABLE messages ALTER COLUMN receiver_id TYPE BIGINT USING 0;
ALTER TABLE messages ADD CONSTRAINT messages_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages ADD CONSTRAINT messages_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE;

-- Update messages_delete table
ALTER TABLE messages_delete ALTER COLUMN sender_id TYPE BIGINT USING 0;
ALTER TABLE messages_delete ALTER COLUMN receiver_id TYPE BIGINT USING 0;
ALTER TABLE messages_delete ADD CONSTRAINT messages_delete_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE messages_delete ADD CONSTRAINT messages_delete_receiver_id_fkey FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_messages_receiver_id ON messages (receiver_id);
CREATE INDEX IF NOT EXISTS idx_messages_delete_receiver_id ON messages_delete (receiver_id);
