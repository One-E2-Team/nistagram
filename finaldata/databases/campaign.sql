-- Adminer 4.8.1 MySQL 8.0.26 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `campaign_parameters`;
CREATE TABLE `campaign_parameters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `start` datetime(3) NOT NULL,
  `end` datetime(3) NOT NULL,
  `campaign_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_campaign_parameters_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `campaign_parameters` (`id`, `created_at`, `updated_at`, `deleted_at`, `start`, `end`, `campaign_id`) VALUES
(1,	'2021-07-10 10:55:41.097',	'2021-07-10 10:55:41.097',	NULL,	'2021-07-10 00:00:00.000',	'2021-07-10 00:00:00.000',	1),
(2,	'2021-07-10 11:01:16.964',	'2021-07-10 11:01:16.964',	NULL,	'2021-07-10 00:00:00.000',	'2021-07-11 00:00:00.000',	2);

DROP TABLE IF EXISTS `campaign_requests`;
CREATE TABLE `campaign_requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `influencer_id` bigint unsigned NOT NULL,
  `request_status` bigint NOT NULL,
  `campaign_parameters_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_campaign_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `campaign_requests` (`id`, `created_at`, `updated_at`, `deleted_at`, `influencer_id`, `request_status`, `campaign_parameters_id`) VALUES
(1,	'2021-07-10 10:55:41.098',	'2021-07-10 10:55:41.098',	NULL,	4,	1,	1);

DROP TABLE IF EXISTS `campaigns`;
CREATE TABLE `campaigns` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `post_id` longtext NOT NULL,
  `agent_id` bigint unsigned NOT NULL,
  `campaign_type` bigint NOT NULL,
  `start` datetime(3) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_campaigns_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `campaigns` (`id`, `created_at`, `updated_at`, `deleted_at`, `post_id`, `agent_id`, `campaign_type`, `start`) VALUES
(1,	'2021-07-10 10:55:41.094',	'2021-07-10 10:55:41.094',	NULL,	'60e97c5fe5d18defbd0eac66',	2,	0,	'2021-07-10 00:00:00.000'),
(2,	'2021-07-10 11:01:16.963',	'2021-07-10 11:01:16.963',	NULL,	'60e97dbde5d18defbd0eac68',	2,	1,	'2021-07-10 00:00:00.000');

DROP TABLE IF EXISTS `interests`;
CREATE TABLE `interests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  KEY `idx_interests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `interests` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`) VALUES
(1,	NULL,	NULL,	NULL,	'Sports'),
(2,	NULL,	NULL,	NULL,	'Alcohol'),
(3,	NULL,	NULL,	NULL,	'Food'),
(4,	NULL,	NULL,	NULL,	'Gaming'),
(5,	NULL,	NULL,	NULL,	'Linux'),
(6,	NULL,	NULL,	NULL,	'Movies'),
(7,	NULL,	NULL,	NULL,	'Music'),
(8,	NULL,	NULL,	NULL,	'Nature'),
(9,	NULL,	NULL,	NULL,	'Programming'),
(10,	NULL,	NULL,	NULL,	'Shopping'),
(11,	NULL,	NULL,	NULL,	'Windows'),
(12,	NULL,	NULL,	NULL,	'XML'),
(13,	NULL,	NULL,	NULL,	'Youtube');

DROP TABLE IF EXISTS `parameters_interests`;
CREATE TABLE `parameters_interests` (
  `campaign_parameters_id` bigint unsigned NOT NULL,
  `interest_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`campaign_parameters_id`,`interest_id`),
  KEY `fk_parameters_interests_interest` (`interest_id`),
  CONSTRAINT `fk_parameters_interests_campaign_parameters` FOREIGN KEY (`campaign_parameters_id`) REFERENCES `campaign_parameters` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_parameters_interests_interest` FOREIGN KEY (`interest_id`) REFERENCES `interests` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `parameters_interests` (`campaign_parameters_id`, `interest_id`) VALUES
(1,	2),
(2,	2),
(1,	3),
(2,	3),
(1,	4),
(2,	4),
(2,	5),
(2,	6),
(2,	7);

DROP TABLE IF EXISTS `timestamps`;
CREATE TABLE `timestamps` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `campaign_parameters_id` bigint unsigned DEFAULT NULL,
  `time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_timestamps_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `timestamps` (`id`, `created_at`, `updated_at`, `deleted_at`, `campaign_parameters_id`, `time`) VALUES
(1,	'2021-07-10 10:55:41.099',	'2021-07-10 10:55:41.099',	NULL,	1,	'2021-07-10 10:00:41.814'),
(2,	'2021-07-10 11:01:16.965',	'2021-07-10 11:01:16.965',	NULL,	2,	'2021-07-10 11:00:17.691');

-- 2021-07-28 16:43:24
