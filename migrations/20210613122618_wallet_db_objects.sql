-- +migrate Up
SET TIME ZONE 'GMT+6';

--currencies table
create table currencies
(
    id         serial                  not null,
    code       varchar                 not null,
    name       varchar                 not null,
    correct_at timestamp default now() not null,
    archive_fl boolean   default false not null
);

comment on column currencies.id is 'primary key';

comment on column currencies.code is 'currency code by ISO 4217';

comment on column currencies.name is 'currency name';

comment on column currencies.correct_at is 'last correct date';

comment on column currencies.archive_fl is 'is currency in archive';

create unique index currencies_code_uindex
    on currencies (code);

create unique index currencies_id_uindex
    on currencies (id);

alter table currencies
    add constraint currencies_pk
        primary key (id);

--wallet table
create table wallet
(
    id          serial                       not null,
    code        varchar                      not null,
    currency_id int                          not null
        constraint wallet_currency_fk
            references currencies (id),
    cli_name    varchar                      not null,
    balance     numeric(18, 2) default 0     not null,
    create_at   timestamp      default now() not null,
    archive_fl  boolean        default false not null
);

create unique index wallet_code_uindex
    on wallet (code);

create unique index wallet_id_uindex
    on wallet (id);

alter table wallet
    add constraint wallet_pk
        primary key (id);

comment on column wallet.id is 'primary_key';
comment on column wallet.code is 'wallet uniq code as account number';
comment on column wallet.cli_name is 'client_name';
comment on column wallet.balance is 'balance of wallet';
comment on column wallet.create_at is 'wallet created at';
comment on column wallet.archive_fl is 'is archive wallet';

--transactions table
create table transactions
(
    id             serial                  not null,
    wallet_id      int                     not null
        constraint transactions_wallet_fk
            references wallet (id),
    date           timestamp default now() not null,
    amount         numeric(18, 2)          not null,
    operation_type varchar(1)              not null
);

comment on column transactions.id is 'primary key';

comment on column transactions.wallet_id is 'transaction is done wallet';

comment on column transactions.amount is 'amount of transaction';

comment on column transactions.operation_type is 'D - debit; C - credit';

create index transactions_date_index
    on transactions (date desc);

create unique index transactions_id_uindex
    on transactions (id);

create index transactions_op_type_index
    on transactions (operation_type);

create index transactions_wallet_index
    on transactions (wallet_id);

alter table transactions
    add constraint transactions_pk
        primary key (id);

--insert currency dictionary
INSERT INTO public.currencies (id, code, name, correct_at, archive_fl)
VALUES (DEFAULT, 'USD', 'US Dollar', DEFAULT, DEFAULT);

-- +migrate Down
drop table currencies;
drop table wallet;
drop table transactions;
