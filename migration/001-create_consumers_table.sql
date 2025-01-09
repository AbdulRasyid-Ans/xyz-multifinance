-- Table consumers
CREATE TABLE IF NOT EXISTS `consumers`(
    `consumer_id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `full_name` VARCHAR(255) NOT NULL,
    `legal_name` VARCHAR(255) NULL,
    `place_of_birth` VARCHAR(255) NOT NULL,
    `dob` DATE NOT NULL,
    `salary` DECIMAL(19, 3) NULL,
    `ktp_image_url` TEXT NULL,
    `selfie_image_url` TEXT NULL,
    `nik` VARCHAR(20) NOT NULL UNIQUE,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL
);
