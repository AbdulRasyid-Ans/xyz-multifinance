-- Table transactions
CREATE TABLE IF NOT EXISTS `transactions`(
    `transaction_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `consumer_id` BIGINT UNSIGNED NOT NULL,
    `loan_id` BIGINT UNSIGNED NULL,
    `description` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `modified_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`consumer_id`),
    FOREIGN KEY (`loan_id`) REFERENCES `loans`(`loan_id`)
);
