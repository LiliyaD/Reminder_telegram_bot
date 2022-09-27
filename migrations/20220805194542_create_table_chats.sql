-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.chats(
   chat_id bigint NOT NULL,
   user_name varchar(100),
   PRIMARY KEY (chat_id)
);

ALTER TABLE IF EXISTS public.activities
   ADD CONSTRAINT fkchat FOREIGN KEY (chat_id) REFERENCES public.chats(chat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.activities
   DROP CONSTRAINT fkchat;

DROP TABLE IF EXISTS public.chats;
-- +goose StatementEnd
