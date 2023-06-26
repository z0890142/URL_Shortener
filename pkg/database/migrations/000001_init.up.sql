CREATE TABLE IF NOT EXISTS `url_mapping`
(
    `url_id` char(6) NOT NULL primary key,
    `original_url` varchar(512) NOT NULL,
    `expired_at` timestamp NOT NULL,
    `expired` tinyint(1) NOT NULL DEFAULT 0,
    `created_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS `available_key`
(
    `key` char(6) NOT NULL primary key,
    `created_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS `allocated_key` (
    `key` char(6)  NOT NULL primary key,
    `created_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP
);