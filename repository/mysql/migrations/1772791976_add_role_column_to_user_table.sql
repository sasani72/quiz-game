-- +migrate Up
-- MYSQL 8.0 sets default role to `user` for old records, since its defined before `admin`
ALTER TABLE `users` ADD COLUMN `role` ENUM('user', 'admin') NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`;