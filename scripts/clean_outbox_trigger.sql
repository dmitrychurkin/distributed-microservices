-- 1. Create the function
CREATE OR REPLACE FUNCTION delete_outbox_row()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM outbox WHERE id = NEW.id;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- 2. Create the trigger
-- Use 'AFTER INSERT' so the row exists long enough for the WAL to log it
CREATE TRIGGER trigger_delete_outbox_after_insert
AFTER INSERT ON outbox
FOR EACH ROW
EXECUTE FUNCTION delete_outbox_row();
