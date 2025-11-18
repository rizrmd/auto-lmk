
INSERT INTO tenants (domain, name, status)
VALUES ('localhost', 'Auto LMK Development', 'active')
ON CONFLICT (domain) DO NOTHING;
