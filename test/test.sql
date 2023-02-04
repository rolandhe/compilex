create table group_info (
    id bigint(20) not null AUTO_INCREMENT primary key,
    `name` varchar(200) not null comment 'user name',
    main_data json not null,
    content text not null,
    bin_data blob not null,
    create_at datetime not null,
    total_amount decimal(10,2) not null,
    remark varchar(128) DEFAULT NULL
) ENGINE=innodb CHARACTER SET utf8mb4 comment 'axxx';

create table user_info (
    id bigint(20) not null AUTO_INCREMENT primary key,
    `name` varchar(200) not null comment 'user name',
    create_at datetime not null,
    modify_at datetime  null
) ENGINE=innodb CHARACTER SET utf8mb4 comment 'axxx';