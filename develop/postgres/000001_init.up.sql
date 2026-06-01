create table if not exists store_message
(
    id            uuid   not null primary key default uuidv7(),
    tenant_id     text   not null,
    external_id   text   not null,
    message       text   not null,
    timestamp     text   not null,
    constraint store_message_unique unique (tenant_id, external_id)
);

create index if not exists store_message_per_tenant_pagination on store_message (tenant_id, id);

create table if not exists store_lock
(
    kind    text    not null,
    constraint store_lock_unique_kind unique (kind)
);

insert into store_lock (kind) values ('update_messages') on conflict do nothing;
