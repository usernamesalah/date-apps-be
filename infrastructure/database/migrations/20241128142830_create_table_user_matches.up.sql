CREATE TABLE user_matches (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `user_uid` varchar(27) NOT NULL,
    `match_uid` varchar(27) NOT NULL,
    `match_type` varchar(20) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT current_timestamp(),
    `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`),
    FOREIGN KEY (`match_uid`) REFERENCES users(`uid`),
    INDEX `user_matches` (`user_uid`, `match_uid`, `match_type`)
);
