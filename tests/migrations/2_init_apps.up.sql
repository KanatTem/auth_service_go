INSERT INTO users (email, pass_hash)
VALUES ('test@email.test', '$2a$10$WAllLHJJRTP6gIA2fLdBYeZJ8tKSEIwQoNzrIbBCYff7CwMIh/KO2')
    ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT
    u.id,
    r.id
FROM users AS u
         CROSS JOIN roles AS r
WHERE
    u.email   = 'test@email.test'
  AND r.name = 'admin'
  AND r.app_id = 1
ON CONFLICT (user_id, role_id) DO NOTHING;

