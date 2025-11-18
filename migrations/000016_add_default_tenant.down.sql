-- +migrate Down
DELETE FROM tenants WHERE domain = 'localhost';
