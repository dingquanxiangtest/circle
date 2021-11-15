CREATE TABLE `schema` (
    `id` VARCHAR(36) NOT NULL,
    `table_name` VARCHAR(42) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20),
    `deleted_at` BIGINT(20),
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `propertie`(
  `id` VARCHAR(36) NOT NULL,
  `schema_id` VARCHAR(36) NOT NULL,
  `name`   VARCHAR(42) NOT NULL,
  `type`      VARCHAR(36) NOT NULL,
  `length`    INT(2) NOT NULL,
  `decimal`   INT(1) NOT NULL,
  `no_null`   TINYINT(1) NOT NULL,
  `comment`   TEXT,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20),
  `deleted_at` BIGINT(20),
   PRIMARY KEY (`id`),
   INDEX `nk_schema_id`(`schema_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `method` (
  `id` VARCHAR(36) NOT NULL,
  `schema_id` VARCHAR(36) NOT NULL,
  `func` VARCHAR(54) NOT NULL,
  `where` TEXT,
  `in` TEXT,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20),
  PRIMARY KEY (`id`),
  INDEX `nk_schema_id`(`schema_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `api` (
  `id` VARCHAR(36) NOT NULL,
  `uri` VARCHAR(62) NOT NULL,
  `method` VARCHAR(6) NOT NULL,
  `func_id` VARCHAR(36) NOT NULL,
  `created_at` BIGINT(20) NOT NULL,
  `updated_at` BIGINT(20),
  PRIMARY KEY (`id`),
  UNIQUE KEY `nk_method_uri`(`method`,`uri`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
