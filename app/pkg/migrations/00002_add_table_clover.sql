-- +goose Up

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `clovers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `coverage_id` bigint NOT NULL,
  `filename` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `line_num` int NOT NULL,
  `hits` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_coverages_coverage_id` (`coverage_id`),
  CONSTRAINT `fk_coverages_coverage_id` FOREIGN KEY (`coverage_id`) REFERENCES `coverages` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;


-- +goose Down
DROP TABLE IF EXISTS `clovers`;
