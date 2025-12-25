-- Seed data for development

-- Insert admin user (password: admin123)
-- Hash generated with: go run tools/generate-password-hash/main.go admin123
-- Default bcrypt hash for "admin123" with cost 10
INSERT INTO users (id, email, password_hash, role) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Insert default settings
INSERT INTO settings (key, value_json) VALUES
('general', '{"app_name": "PK-FK Discovery", "timezone": "UTC", "date_format": "YYYY-MM-DD"}'),
('security', '{"password_min_length": 8, "session_timeout": 3600, "require_2fa": false}'),
('ai', '{"default_provider": null, "rate_limit": 100, "max_tokens": 4000}'),
('registry', '{"default_maturity_level": "L2", "auto_approve": false}'),
('engine', '{"default_sample_mode": true, "default_timeout": 300, "max_concurrency": 5}')
ON CONFLICT (key) DO NOTHING;

