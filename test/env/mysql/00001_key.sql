DROP TABLE IF EXISTS `key`;
CREATE TABLE IF NOT EXISTS `key`
(
    `id` bigint NOT NULL AUTO_INCREMENT  primary key COMMENT '自增ID',
    `key` char(6) NOT NULL,
    `used` tinyint(1) NOT NULL DEFAULT 0,
    `created_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`              timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);