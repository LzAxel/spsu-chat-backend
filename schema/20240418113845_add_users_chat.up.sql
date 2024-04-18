CREATE TABLE chat_users (
    user_id BIGINT references users(id),
    chat_id BIGINT references chats(id),
    UNIQUE(user_id, chat_id)
)