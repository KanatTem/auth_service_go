DROP TRIGGER IF EXISTS trg_prevent_orphan_user_roles ON user_roles;

DROP FUNCTION IF EXISTS prevent_orphan_user_roles;

DROP TRIGGER IF EXISTS trg_create_default_user_role ON apps;

DROP FUNCTION IF EXISTS create_default_user_role;
