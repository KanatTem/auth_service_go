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

CREATE OR REPLACE FUNCTION remove_user_if_no_roles()
    RETURNS TRIGGER AS $$
BEGIN
    -- Check if the user has any remaining roles
    IF NOT EXISTS (
        SELECT 1
        FROM user_roles
        WHERE user_id = OLD.user_id
    ) THEN
        -- Delete the user if no roles remain
        DELETE FROM users WHERE id = OLD.user_id;
    END IF;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_remove_user_if_no_roles
    AFTER DELETE ON user_roles
    FOR EACH ROW
EXECUTE FUNCTION remove_user_if_no_roles();

