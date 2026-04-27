CREATE TYPE message_type AS ENUM ('CHAT', 'SYSTEM');

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payload TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    type message_type DEFAULT 'CHAT'
);

CREATE INDEX idx_messages_receiver_id ON messages (receiver_id);
