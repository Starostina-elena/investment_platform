CREATE OR REPLACE FUNCTION trg_set_project_completed()
    RETURNS trigger AS
    $$
    BEGIN
        IF NEW.current_money IS NOT NULL AND NEW.wanted_money IS NOT NULL AND NEW.current_money >= NEW.wanted_money THEN
            NEW.is_completed := TRUE;
        END IF;
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

    CREATE TRIGGER projects_set_completed
    BEFORE INSERT OR UPDATE OF current_money, wanted_money ON projects
    FOR EACH ROW
    EXECUTE FUNCTION trg_set_project_completed();