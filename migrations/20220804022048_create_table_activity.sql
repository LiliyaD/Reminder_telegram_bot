-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.activities(
   chat_id bigint NOT NULL,
   act_name varchar(120) NOT NULL,
   begin_date date,
   end_date date,
   times_per_day smallint,
   quantity_per_time real,
   PRIMARY KEY (chat_id, act_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.activities;
-- +goose StatementEnd
