-- Table consumer_limits
CREATE TABLE IF NOT EXISTS `consumer_limits`(
    `consumer_limit_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `consumer_id` BIGINT UNSIGNED NOT NULL,
    `tenure` ENUM('1', '2', '3', '6') NOT NULL,
    `limit_amount` DOUBLE NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`consumer_id`)
);
