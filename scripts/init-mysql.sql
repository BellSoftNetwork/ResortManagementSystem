-- Create databases for api-legacy and api-core
CREATE DATABASE IF NOT EXISTS `rms-legacy` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS `rms-core` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Grant privileges to rms user
GRANT ALL PRIVILEGES ON `rms-legacy`.* TO 'rms'@'%';
GRANT ALL PRIVILEGES ON `rms-core`.* TO 'rms'@'%';
FLUSH PRIVILEGES;