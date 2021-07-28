-- Adminer 4.8.1 MySQL 8.0.26 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `agent_products`;
CREATE TABLE `agent_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `quantity` bigint unsigned NOT NULL,
  `price_per_item` double NOT NULL,
  `valid_from` datetime(3) NOT NULL,
  `is_valid` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_agent_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `agent_products` (`id`, `created_at`, `updated_at`, `deleted_at`, `user_id`, `product_id`, `quantity`, `price_per_item`, `valid_from`, `is_valid`) VALUES
(1,	'2021-07-10 10:53:05.999',	'2021-07-10 10:53:05.999',	NULL,	1,	1,	30,	85,	'2021-07-10 10:53:05.999',	1);

DROP TABLE IF EXISTS `campaign_stats`;
CREATE TABLE `campaign_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `campaign_id` bigint unsigned DEFAULT NULL,
  `post_link` longtext,
  `statistics_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_campaign_stats_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `influencer_stats`;
CREATE TABLE `influencer_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `campaign_stat_id` bigint unsigned DEFAULT NULL,
  `influencer_username` longtext,
  `statistics_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_influencer_stats_deleted_at` (`deleted_at`),
  KEY `fk_campaign_stats_influencer_stat` (`campaign_stat_id`),
  CONSTRAINT `fk_campaign_stats_influencer_stat` FOREIGN KEY (`campaign_stat_id`) REFERENCES `campaign_stats` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `interest_stats`;
CREATE TABLE `interest_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `campaign_stat_id` bigint unsigned DEFAULT NULL,
  `interest` longtext,
  `statistics_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_interest_stats_deleted_at` (`deleted_at`),
  KEY `fk_campaign_stats_interest_stat` (`campaign_stat_id`),
  CONSTRAINT `fk_campaign_stats_interest_stat` FOREIGN KEY (`campaign_stat_id`) REFERENCES `campaign_stats` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `items`;
CREATE TABLE `items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_id` bigint unsigned NOT NULL,
  `quantity` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_items_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `timestamp` datetime(3) NOT NULL,
  `full_price` double NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `agent_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `privileges`;
CREATE TABLE `privileges` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_privileges_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `privileges` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`) VALUES
(1,	NULL,	NULL,	NULL,	'CREATE_PRODUCT'),
(2,	NULL,	NULL,	NULL,	'DELETE_PRODUCT'),
(3,	NULL,	NULL,	NULL,	'UPDATE_PRODUCT'),
(4,	NULL,	NULL,	NULL,	'READ_PRODUCT'),
(5,	NULL,	NULL,	NULL,	'CREATE_ORDER'),
(6,	NULL,	NULL,	NULL,	'CREATE_TOKEN'),
(7,	NULL,	NULL,	NULL,	'READ_POSTS'),
(8,	NULL,	NULL,	NULL,	'READ_CAMPAIGNS'),
(9,	NULL,	NULL,	NULL,	'CREATE_CAMPAIGN'),
(10,	NULL,	NULL,	NULL,	'EDIT_CAMPAIGN'),
(11,	NULL,	NULL,	NULL,	'DELETE_CAMPAIGN');

DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext NOT NULL,
  `picture_path` longtext NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `products` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `picture_path`) VALUES
(1,	'2021-07-10 10:53:05.990',	'2021-07-10 10:53:05.990',	NULL,	'kavabon',	'8c8e7c00-7636-4f05-9244-bb8cdbd13cd6.png');

DROP TABLE IF EXISTS `role_privileges`;
CREATE TABLE `role_privileges` (
  `role_id` bigint unsigned NOT NULL,
  `privilege_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`privilege_id`),
  KEY `fk_role_privileges_privilege` (`privilege_id`),
  CONSTRAINT `fk_role_privileges_privilege` FOREIGN KEY (`privilege_id`) REFERENCES `privileges` (`id`),
  CONSTRAINT `fk_role_privileges_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `role_privileges` (`role_id`, `privilege_id`) VALUES
(2,	1),
(2,	2),
(2,	3),
(1,	4),
(2,	4),
(1,	5),
(2,	6),
(2,	7),
(2,	8),
(2,	9),
(2,	10),
(2,	11);

DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `roles` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`) VALUES
(1,	NULL,	NULL,	NULL,	'CUSTOMER'),
(2,	NULL,	NULL,	NULL,	'AGENT');

DROP TABLE IF EXISTS `statistics`;
CREATE TABLE `statistics` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `num_of_likes` bigint unsigned DEFAULT NULL,
  `num_of_dislikes` bigint unsigned DEFAULT NULL,
  `num_of_visits` bigint unsigned DEFAULT NULL,
  `num_of_comments` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_statistics_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_user_roles_role` (`role_id`),
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `user_roles` (`user_id`, `role_id`) VALUES
(1,	2);

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `api_token` longtext,
  `address` longtext,
  `is_validated` tinyint(1) NOT NULL,
  `validation_uid` longtext,
  `validation_expire` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `email`, `password`, `api_token`, `address`, `is_validated`, `validation_uid`, `validation_expire`) VALUES
(1,	'2021-06-28 00:00:00.000',	'2021-07-10 10:52:25.955',	NULL,	'agent@gmail.com',	'$2a$10$LsQPBxzsx/IbxK9PutpiZ.aOqQ4SsUUnGU4qhPf0dtWdggIzFMJ1W',	'13752fef-094c-4df1-895d-3f44504f3851',	'address',	1,	'',	'2021-06-28 00:00:00.000');

-- 2021-07-28 16:42:48
