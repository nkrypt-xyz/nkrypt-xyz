-- Removing default admin if it was created by the application
DELETE FROM users WHERE user_name = 'admin';

