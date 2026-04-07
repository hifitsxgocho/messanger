-- +goose Up
CREATE TABLE conversations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_a_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_b_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_conversation UNIQUE (user_a_id, user_b_id)
);

CREATE INDEX idx_conv_user_a ON conversations(user_a_id);
CREATE INDEX idx_conv_user_b ON conversations(user_b_id);

-- +goose Down
DROP TABLE IF EXISTS conversations;
