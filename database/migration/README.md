# Database Migration Best Practices Guide

## Migration File Naming Convention

```
<sequence>_<description>.sql
Example: 001_create_examples.sql
```

## Up/Down Structure

Always include both UP and DOWN migrations:

```sql
-- +migrate Up
CREATE TABLE ...

-- +migrate Down
DROP TABLE ...
```

## Index Guidelines

### Single Column Indexes
- Index columns used in WHERE clauses
- Index columns used in ORDER BY
- Index foreign key columns

### Composite Indexes
- Place highly selective columns first
- Consider query patterns
- Avoid redundant indexes

### Index Naming Convention
```sql
INDEX `idx_<table>_<column1>_<column2>` (`column1`, `column2`)
```

## Table Design Best Practices

### Always Include
1. **Timestamps**: `created_at`, `updated_at`
2. **Soft Deletes**: `deleted_at` (nullable timestamp)
3. **Primary Key**: Use BIGINT UNSIGNED for scalability

### Charset & Collate
```sql
DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
```

### Engine
```sql
ENGINE=InnoDB
```

## Foreign Key Best Practices

```sql
CREATE TABLE `orders` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    PRIMARY KEY (`id`),
    
    -- Foreign key with constraint
    INDEX `idx_orders_user_id` (`user_id`),
    CONSTRAINT `fk_orders_user_id` 
        FOREIGN KEY (`user_id`) 
        REFERENCES `users` (`id`) 
        ON DELETE RESTRICT 
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## Common Query Patterns & Indexes

### Pattern: Soft Delete with Pagination
```sql
-- Composite index for soft delete queries
INDEX `idx_deleted_created` (`deleted_at`, `created_at`)
```

### Pattern: User-specific queries
```sql
-- For queries: SELECT * FROM orders WHERE user_id = ? AND status = ?
INDEX `idx_orders_user_status` (`user_id`, `status`)
```

### Pattern: Date range queries
```sql
-- For queries: SELECT * FROM orders WHERE created_at BETWEEN ? AND ?
INDEX `idx_orders_created_at` (`created_at`)
```

## Migration Safety Rules

1. **Never modify existing migrations** - create new ones
2. **Always test DOWN migrations** before deploying
3. **Use transactions** for data migrations when possible
4. **Back up data** before destructive operations
5. **Add comments** explaining complex changes

## Example: Adding a Column Safely

```sql
-- +migrate Up
ALTER TABLE `examples` 
ADD COLUMN `description` TEXT NULL AFTER `name`,
ADD INDEX `idx_examples_description` (`description`(100));

-- +migrate Down
ALTER TABLE `examples` 
DROP COLUMN `description`;
```

## Example: Adding Foreign Key

```sql
-- +migrate Up
ALTER TABLE `orders` 
ADD CONSTRAINT `fk_orders_user_id` 
    FOREIGN KEY (`user_id`) 
    REFERENCES `users` (`id`) 
    ON DELETE RESTRICT 
    ON UPDATE CASCADE;

-- +migrate Down
ALTER TABLE `orders` 
DROP FOREIGN KEY `fk_orders_user_id`;
```
