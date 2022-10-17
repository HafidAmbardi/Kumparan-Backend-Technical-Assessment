CREATE DATABASE `ricemilkdb2`

-- ricemilkdb2.articles definition

CREATE TABLE `articles` (
  `id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `author_id` varchar(100) DEFAULT NULL,
  `title` varchar(100) DEFAULT NULL,
  `body` varchar(100) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `articles_FK` (`author_id`),
  CONSTRAINT `articles_FK` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ricemilkdb2.authors definition

CREATE TABLE `authors` (
  `id` varchar(100) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;