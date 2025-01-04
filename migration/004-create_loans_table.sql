-- Table loans
CREATE TABLE IF NOT EXISTS `loans`(
    `loan_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `consumer_limit_id` BIGINT UNSIGNED NOT NULL,
    `consumer_id` BIGINT UNSIGNED NOT NULL,
    `loan_amount` DOUBLE NOT NULL,
    `paid_loan_amount` DOUBLE NOT NULL DEFAULT 0,
    `contract_number` VARCHAR(255) NOT NULL UNIQUE,
    `interest_rate` INT NOT NULL DEFAULT 0,
    `interest_amount` DOUBLE NOT NULL DEFAULT 0,
    `paid_interest_amount` DOUBLE NOT NULL DEFAULT 0,
    `loan_status` ENUM('on_going', 'finish', 'late') NOT NULL DEFAULT 'on_going',
    `due_date` DATE NOT NULL,
    `installment` INT NOT NULL DEFAULT 0,
    `asset_name` VARCHAR(255) NOT NULL,
    `merchant_id` BIGINT UNSIGNED NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL,
    FOREIGN KEY (`consumer_limit_id`) REFERENCES `consumer_limits`(`consumer_limit_id`),
    FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`consumer_id`),
    FOREIGN KEY (`merchant_id`) REFERENCES `merchants`(`merchant_id`)
);
