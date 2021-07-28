-- Adminer 4.8.1 MySQL 8.0.26 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `agent_requests`;
CREATE TABLE `agent_requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `profile_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `profile_id` (`profile_id`),
  KEY `idx_agent_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  PRIMARY KEY (`id`),
  KEY `idx_categories_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `categories` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`) VALUES
(1,	NULL,	NULL,	NULL,	'Influencer'),
(2,	NULL,	NULL,	NULL,	'Sports'),
(3,	NULL,	NULL,	NULL,	'New/media'),
(4,	NULL,	NULL,	NULL,	'Business'),
(5,	NULL,	NULL,	NULL,	'Brand'),
(6,	NULL,	NULL,	NULL,	'Organization');

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

DROP TABLE IF EXISTS `person_interests`;
CREATE TABLE `person_interests` (
  `personal_data_id` bigint unsigned NOT NULL,
  `interest_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`personal_data_id`,`interest_id`),
  KEY `fk_person_interests_interest` (`interest_id`),
  CONSTRAINT `fk_person_interests_interest` FOREIGN KEY (`interest_id`) REFERENCES `interests` (`id`),
  CONSTRAINT `fk_person_interests_personal_data` FOREIGN KEY (`personal_data_id`) REFERENCES `personal_data` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `person_interests` (`personal_data_id`, `interest_id`) VALUES
(2,	1),
(4,	1),
(8,	1),
(3,	2),
(6,	2),
(8,	2),
(3,	3),
(6,	3),
(7,	3),
(3,	4),
(6,	4),
(7,	4),
(8,	4),
(5,	5),
(6,	5),
(5,	6),
(6,	6),
(8,	6),
(5,	7),
(6,	7),
(7,	7),
(8,	7),
(4,	8),
(7,	10),
(2,	12),
(4,	12),
(7,	13);

DROP TABLE IF EXISTS `personal_data`;
CREATE TABLE `personal_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `surname` longtext,
  `telephone` longtext,
  `gender` longtext,
  `birth_date` longtext,
  `profile_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_personal_data_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `personal_data` (`id`, `created_at`, `updated_at`, `deleted_at`, `name`, `surname`, `telephone`, `gender`, `birth_date`, `profile_id`) VALUES
(1,	NULL,	NULL,	NULL,	'Name',	'Surname',	'Telephone',	'Male',	'2021-06-02',	1),
(2,	'2021-07-09 14:22:57.099',	'2021-07-09 14:22:57.099',	NULL,	'aaa',	'aaa',	'065041065',	'male',	'2021-07-01',	2),
(3,	'2021-07-09 14:23:56.968',	'2021-07-09 14:23:56.968',	NULL,	'bbb',	'bbb',	'065041065',	'female',	'2021-07-02',	3),
(4,	'2021-07-09 14:24:54.379',	'2021-07-09 14:24:54.379',	NULL,	'ccc',	'ccc',	'065041065',	'male',	'2021-07-03',	4),
(5,	'2021-07-09 14:25:47.559',	'2021-07-09 14:25:47.559',	NULL,	'ddd',	'ddd',	'065041065',	'female',	'2021-07-04',	5),
(6,	'2021-07-09 14:27:06.486',	'2021-07-09 14:27:06.486',	NULL,	'eee',	'eee',	'065041065',	'male',	'2021-07-05',	6),
(7,	'2021-07-09 14:28:38.956',	'2021-07-09 14:28:38.956',	NULL,	'fff',	'fff',	'065041065',	'female',	'2021-07-06',	7),
(8,	'2021-07-09 14:29:47.365',	'2021-07-09 14:29:47.365',	NULL,	'ggg',	'ggg',	'065041065',	'male',	'2021-07-07',	8);

DROP TABLE IF EXISTS `profile_settings`;
CREATE TABLE `profile_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_private` tinyint(1) DEFAULT NULL,
  `can_receive_message_from_unknown` tinyint(1) DEFAULT NULL,
  `can_be_tagged` tinyint(1) DEFAULT NULL,
  `profile_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_profile_settings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `profile_settings` (`id`, `created_at`, `updated_at`, `deleted_at`, `is_private`, `can_receive_message_from_unknown`, `can_be_tagged`, `profile_id`) VALUES
(1,	NULL,	NULL,	NULL,	1,	1,	1,	1),
(2,	'2021-07-09 14:22:57.096',	'2021-07-09 14:22:57.096',	NULL,	0,	1,	1,	2),
(3,	'2021-07-09 14:23:56.967',	'2021-07-09 14:23:56.967',	NULL,	1,	1,	1,	3),
(4,	'2021-07-09 14:24:54.378',	'2021-07-09 14:24:54.378',	NULL,	0,	1,	1,	4),
(5,	'2021-07-09 14:25:47.558',	'2021-07-09 14:25:47.558',	NULL,	1,	1,	1,	5),
(6,	'2021-07-09 14:27:06.484',	'2021-07-09 14:27:06.484',	NULL,	0,	1,	1,	6),
(7,	'2021-07-09 14:28:38.955',	'2021-07-09 14:28:38.955',	NULL,	0,	1,	1,	7),
(8,	'2021-07-09 14:29:47.364',	'2021-07-09 14:29:47.364',	NULL,	1,	1,	1,	8);

DROP TABLE IF EXISTS `profiles`;
CREATE TABLE `profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(191) NOT NULL,
  `email` varchar(191) NOT NULL,
  `biography` longtext,
  `website` longtext,
  `is_verified` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `profiles` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `email`, `biography`, `website`, `is_verified`) VALUES
(1,	'2021-06-28 00:00:00.000',	'2021-06-28 00:00:00.000',	NULL,	'adminko',	'adminko@gmail.com',	'',	'',	1),
(2,	'2021-07-09 14:22:57.095',	'2021-07-09 14:22:57.095',	NULL,	'agent',	'agent@gmail.com',	'fghgfds',	'https://localhost:83',	0),
(3,	'2021-07-09 14:23:56.965',	'2021-07-09 14:23:56.965',	NULL,	'bbb',	'b@gmail.com',	'fds',	'',	0),
(4,	'2021-07-09 14:24:54.376',	'2021-07-09 14:24:54.376',	NULL,	'ccc',	'c@gmail.com',	'grdx',	'',	0),
(5,	'2021-07-09 14:25:47.557',	'2021-07-09 14:25:47.557',	NULL,	'ddd',	'd@gmail.com',	'gfds',	'',	0),
(6,	'2021-07-09 14:27:06.483',	'2021-07-09 14:27:06.483',	NULL,	'eee',	'e@gmail.com',	'fesg',	'',	0),
(7,	'2021-07-09 14:28:38.954',	'2021-07-09 14:28:38.954',	NULL,	'fff',	'f@gmail.com',	'gdsgsd',	'',	0),
(8,	'2021-07-09 14:29:47.364',	'2021-07-09 14:29:47.364',	NULL,	'ggg',	'g@gmail.com',	'gsdfgs',	'',	0);

DROP TABLE IF EXISTS `verification_requests`;
CREATE TABLE `verification_requests` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `profile_id` bigint unsigned DEFAULT NULL,
  `name` longtext,
  `surname` longtext,
  `verification_status` bigint DEFAULT NULL,
  `image_path` longtext,
  `category_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_verification_requests_deleted_at` (`deleted_at`),
  KEY `fk_verification_requests_category` (`category_id`),
  CONSTRAINT `fk_verification_requests_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 2021-07-28 16:47:11
