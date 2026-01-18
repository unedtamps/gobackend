CREATE TYPE "status" AS ENUM (
  'OK',
  'DELETED'
);

create or replace function on_update()
returns trigger
as $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
language 'plpgsql'
;
