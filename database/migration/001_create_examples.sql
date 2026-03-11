-- +migrate Up
-- Create examples table with best practices

CREATE TABLE IF NOT EXISTS `examples` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` TIMESTAMP NULL DEFAULT NULL,
    
    -- Primary key
    PRIMARY KEY (`id`),
    
    -- Indexes for common queries
    INDEX `idx_examples_name` (`name`),
    INDEX `idx_examples_created_at` (`created_at`),
    INDEX `idx_examples_deleted_at` (`deleted_at`),
    
    -- Composite index for soft delete queries
    INDEX `idx_examples_deleted_created` (`deleted_at`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS `examples`;
