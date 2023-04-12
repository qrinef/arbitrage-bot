create table pools
(
    id           serial
        primary key,
    block_number integer     not null,
    factory      varchar(42) not null,
    address      varchar(42) not null,
    token0       varchar(42) not null,
    token1       varchar(42) not null,
    constraint pools_pk
        unique (block_number, factory, address, token0, token1)
);
