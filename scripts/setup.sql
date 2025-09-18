CREATE DATABASE IF NOT EXISTS occupation_db
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE occupation_db;

CREATE TABLE IF NOT EXISTS occupations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    seq VARCHAR(20) NOT NULL UNIQUE COMMENT '职业编号',
    gbm VARCHAR(20) COMMENT 'GBM编码',
    name VARCHAR(200) NOT NULL COMMENT '职业名称',
    level TINYINT NOT NULL COMMENT '层级: 1-大类, 2-中类, 3-小类, 4-细类',
    parent_seq VARCHAR(20) COMMENT '父级编号',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_seq (parent_seq),
    INDEX idx_level (level),
    FOREIGN KEY (parent_seq) REFERENCES occupations(seq) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;