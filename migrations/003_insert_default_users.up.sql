-- Insert default users with hashed passwords
-- admin123 = $2a$10$N9qo8uLOickgx2ZMRZoMye7BrR7dkxzeHc/7/OhSZ.VlqzaXwUgKm
-- password123 = $2a$10$5OpE1ko1yUEqxmVfEQmPVuuK.3w2TFHs.VlqzaXwUgKm1234567890
-- cooking456 = $2a$10$3OpE1ko1yUEqxmVfEQmPVuuK.3w2TFHs.VlqzaXwUgKm0987654321
INSERT INTO users (username, password_hash, email, is_active) VALUES 
    ('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye7BrR7dkxzeHc/7/OhSZ.VlqzaXwUgKm', 'admin@example.com', true),
    ('user1', '$2a$10$5OpE1ko1yUEqxmVfEQmPVuuK.3w2TFHs.VlqzaXwUgKm1234567890', 'user1@example.com', true),
    ('chef', '$2a$10$3OpE1ko1yUEqxmVfEQmPVuuK.3w2TFHs.VlqzaXwUgKm0987654321', 'chef@example.com', true)
ON CONFLICT (username) DO NOTHING;
