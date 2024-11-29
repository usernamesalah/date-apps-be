BEGIN;

CREATE TABLE
    premium_config (
        `id` bigint (20) unsigned NOT NULL AUTO_INCREMENT,
        `uid` varchar(27) NOT NULL,
        name VARCHAR(255) NOT NULL,
        description TEXT,
        price int NOT NULL,
        quota int DEFAULT 0, -- if NULL / 0 will be unlimited
        expired_day int DEFAULT 0, -- if NULL / 0 will be forever
        is_active boolean DEFAULT true,
        `created_at` datetime NOT NULL DEFAULT current_timestamp(),
        `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
        PRIMARY KEY (`id`),
        UNIQUE KEY `config_uid_unique` (`uid`),
        UNIQUE KEY `config_name_unique` (`name`),
        UNIQUE KEY `config_price_quota_unique` (`price` , `quota`, `expired_day`),
        INDEX `config_price_quota_idx` (`price`, `quota`, `expired_day`)
    );

INSERT INTO premium_config (uid, name, description, price, quota, expired_day)
VALUES 
('basic123', 'Basic Plan', 'Basic subscription plan', 100, 10, 30),
('standar123', 'Standard Plan', 'Standard subscription plan', 200, 20, 60),
('premium123', 'Premium Plan', 'Premium subscription plan', 300, 0, 0);


CREATE TABLE user_premium (
    `id` bigint (20) unsigned NOT NULL AUTO_INCREMENT,
    `uid` varchar(27) NOT NULL,
    `user_uid` varchar(27) NOT NULL,
    `premium_config_uid` varchar(27) NOT NULL,
    `started_at` date DEFAULT NULL,
    `ended_at` date DEFAULT NULL,
    `quota` int DEFAULT 0, -- avoid change on premium_config
    `created_at` datetime NOT NULL DEFAULT current_timestamp(),
    `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_uid`) REFERENCES users(`uid`),
    FOREIGN KEY (`premium_config_uid`) REFERENCES premium_config(`uid`),
    UNIQUE KEY `user_premium_uid_unique` (`uid`),
    INDEX `user_premium_user_uid_config_uid_start_end_at_idx` (`user_uid` ,`premium_config_uid`,`started_at`, `ended_at`)
);

COMMIT;