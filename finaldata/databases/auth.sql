-- Adminer 4.8.1 MySQL 8.0.26 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

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
(1,	NULL,	NULL,	NULL,	'READ_PROFILE_DATA'),
(2,	NULL,	NULL,	NULL,	'EDIT_PROFILE_DATA'),
(3,	NULL,	NULL,	NULL,	'CREATE_CONNECTION'),
(4,	NULL,	NULL,	NULL,	'READ_CONNECTION_STATUS'),
(5,	NULL,	NULL,	NULL,	'READ_CONNECTION_REQUESTS'),
(6,	NULL,	NULL,	NULL,	'EDIT_CONNECTION_STATUS'),
(7,	NULL,	NULL,	NULL,	'CREATE_POST'),
(8,	NULL,	NULL,	NULL,	'READ_NOT_ONLY_PUBLIC_POSTS'),
(9,	NULL,	NULL,	NULL,	'READ_VERIFICATION_REQUESTS'),
(10,	NULL,	NULL,	NULL,	'CREATE_VERIFICATION_REQUEST'),
(11,	NULL,	NULL,	NULL,	'UPDATE_VERIFICATION_REQUEST'),
(12,	NULL,	NULL,	NULL,	'REACT_ON_POST'),
(13,	NULL,	NULL,	NULL,	'REPORT_POST'),
(14,	NULL,	NULL,	NULL,	'READ_REACTIONS'),
(15,	NULL,	NULL,	NULL,	'READ_REPORTS'),
(16,	NULL,	NULL,	NULL,	'DELETE_POST'),
(17,	NULL,	NULL,	NULL,	'DELETE_PROFILE'),
(18,	NULL,	NULL,	NULL,	'CREATE_AGENT_REQUEST'),
(19,	NULL,	NULL,	NULL,	'READ_AGENT_REQUEST'),
(20,	NULL,	NULL,	NULL,	'EDIT_AGENT_REQUEST'),
(21,	NULL,	NULL,	NULL,	'CREATE_AGENT'),
(22,	NULL,	NULL,	NULL,	'READ_API_TOKEN'),
(23,	NULL,	NULL,	NULL,	'AGENT_API_ACCESS'),
(24,	NULL,	NULL,	NULL,	'EDIT_CAMPAIGN_REQUEST'),
(25,	NULL,	NULL,	NULL,	'READ_CAMPAIGN_REQUEST'),
(26,	NULL,	NULL,	NULL,	'MESSAGING');

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
(2,	4),
(2,	5),
(2,	6),
(2,	7),
(2,	8),
(1,	9),
(2,	10),
(1,	11),
(2,	12),
(2,	13),
(2,	14),
(1,	15),
(1,	16),
(1,	17),
(2,	18),
(1,	19),
(1,	20),
(1,	21),
(3,	22),
(4,	23),
(2,	24),
(2,	25),
(1,	26),
(2,	26);

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
(1,	NULL,	NULL,	NULL,	'ADMIN'),
(2,	NULL,	NULL,	NULL,	'REGULAR'),
(3,	NULL,	NULL,	NULL,	'AGENT'),
(4,	NULL,	NULL,	NULL,	'AGENT_API_CLIENT');

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
(1,	1),
(1,	2),
(2,	2),
(3,	2),
(4,	2),
(5,	2),
(6,	2),
(7,	2),
(8,	2),
(1,	3),
(2,	3),
(2,	4);

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `profile_id` bigint unsigned NOT NULL,
  `password` longtext NOT NULL,
  `api_token` longtext,
  `is_deleted` tinyint(1) NOT NULL,
  `is_validated` tinyint(1) NOT NULL,
  `email` varchar(191) NOT NULL,
  `username` varchar(191) NOT NULL,
  `validation_uid` longtext,
  `validation_expire` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `profile_id` (`profile_id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `profile_id`, `password`, `api_token`, `is_deleted`, `is_validated`, `email`, `username`, `validation_uid`, `validation_expire`) VALUES
(1,	NULL,	NULL,	NULL,	1,	'$2y$10$u1tGi0miWp8uPHw/37JPbOJTbNExmAHh9/fRVGswkNLCQyAVTDjya',	'',	0,	1,	'adminko@gmail.com',	'adminko',	'',	'2021-06-15 11:39:08.096'),
(2,	'2021-07-09 14:22:58.297',	'2021-07-10 10:51:56.736',	NULL,	2,	'$2a$10$UiefLwzubTldyXWSAqBFNuEA0x4FyPCmO6PFSTN3HChOFzRQQX.4a',	'$2a$10$lhk9KTO3Hl9eQjaquPm4Jum1JrbmVX6w8wnUBbGFdfamKYMbBfZFW',	0,	1,	'agent@gmail.com',	'agent',	'',	'2021-07-09 14:31:27.635'),
(3,	'2021-07-09 14:23:57.144',	'2021-07-09 14:31:37.123',	NULL,	3,	'$2a$10$9YiCIsJPBgmjRTsVSipJoODhcKTX.N.yEi2SrDftTAysnnIDME7qe',	'',	0,	1,	'b@gmail.com',	'bbb',	'',	'2021-07-09 14:31:37.122'),
(4,	'2021-07-09 14:24:54.534',	'2021-07-09 14:31:48.586',	NULL,	4,	'$2a$10$NDvV0XnMzIDFiEynJ9z/ZOcqwfwuS.6px.nKGPYZ5G925purW8igq',	'',	0,	1,	'c@gmail.com',	'ccc',	'',	'2021-07-09 14:31:48.585'),
(5,	'2021-07-09 14:25:47.641',	'2021-07-09 14:31:58.066',	NULL,	5,	'$2a$10$iwyKl4lNcCoebexTvJcKauvcTRFmcaQuPExz0KygpdkYyBROHVk3e',	'',	0,	1,	'd@gmail.com',	'ddd',	'',	'2021-07-09 14:31:58.066'),
(6,	'2021-07-09 14:27:06.640',	'2021-07-09 14:32:10.057',	NULL,	6,	'$2a$10$s/mz.Lg2kBZb3.nb3oFhu.mG57u6BxtNv0/rhejOuGWc0P.MEL4um',	'',	0,	1,	'e@gmail.com',	'eee',	'',	'2021-07-09 14:32:10.056'),
(7,	'2021-07-09 14:28:39.051',	'2021-07-09 14:32:17.819',	NULL,	7,	'$2a$10$E/TqoKF7LXqtOcgcEvpdYe1jHKbGZdd8ZeARnxlIBxAOk6fqHzF2O',	'',	0,	1,	'f@gmail.com',	'fff',	'',	'2021-07-09 14:32:17.818'),
(8,	'2021-07-09 14:29:47.446',	'2021-07-09 14:32:24.113',	NULL,	8,	'$2a$10$9UX/e8CBDIoMUcCrjS89w.ZSeC3pvrOqIMq6zvV.x26erAX7JtOYO',	'',	0,	1,	'g@gmail.com',	'ggg',	'',	'2021-07-09 14:32:24.112');

-- 2021-07-28 16:43:06
