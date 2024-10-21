set timezone = 'Europe/Moscow';

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
            CREATE TYPE role AS ENUM ('user', 'admin','superAdmin');
        END IF;
END $$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chan_status') THEN
            CREATE TYPE chan_status AS ENUM ('kicked','administrator','left','member','unknown');
        END IF;
    END $$;

DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pub_status') THEN
            CREATE TYPE pub_status AS ENUM ('sent','awaits','error_on_sending','deleted_by_bot','error_on_deleting');
        END IF;
    END $$;

create table if not exists "user"
(
    id           bigint unique,
    tg_username  text                not null,
    created_at   timestamp           not null,
    channel_from varchar(150)        null,
    user_role    role default 'user' not null,
    primary key (id)
);

create unique index if not exists user_tg_username_idx on "user" using btree (tg_username);

create table if not exists channel(
    id int generated always as identity,
    tg_id bigint unique not null,
    channel_name varchar(150) null,
    channel_url varchar(150) null,
    channel_status chan_status not null,
    primary key (id)
);

create table if not exists publication(
    id int generated always as identity,
    channel_id int not null,
    publication_status pub_status default 'awaits' not null,
    text text  not null,
    image varchar(200) null,
    button_url varchar(150) null,
    button_text varchar(150) null,
    publication_date timestamp with time zone not null,
    delete_date timestamp with time zone,
    message_id bigint default null,
    primary key (id),
    foreign key (channel_id)
        references channel (id) on delete cascade
);


ALTER TABLE publication
    ALTER COLUMN publication_status DROP NOT NULL;

alter table publication add column message_id bigint default null;
