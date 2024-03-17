CREATE TABLE "todolist" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT NOT NULL,
  "status" status DEFAULT 'OK' NOT NULL,
  "created_at" TIMESTAMP DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE TRIGGER todolist_on_update
BEFORE UPDATE ON todolist
FOR EACH ROW
EXECUTE FUNCTION on_update();

