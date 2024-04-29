CREATE OR REPLACE FUNCTION prevent_orphan_user_roles()
    RETURNS TRIGGER AS $$
DECLARE
    cnt INT;
BEGIN
    SELECT COUNT(*) INTO cnt
    FROM user_roles
    WHERE user_id = OLD.user_id;

    IF cnt <= 1 THEN
        RAISE EXCEPTION 'cannot remove last role for user %', OLD.user_id;
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_orphan_user_roles
    BEFORE DELETE
    ON user_roles
    FOR EACH ROW
EXECUTE FUNCTION prevent_orphan_user_roles();

CREATE OR REPLACE FUNCTION create_default_user_role()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO roles (app_id, name)
    VALUES (NEW.id, 'user');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_default_user_role
    AFTER INSERT
    ON apps
    FOR EACH ROW
EXECUTE FUNCTION create_default_user_role();