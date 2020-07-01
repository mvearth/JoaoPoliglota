create database joaopoliglota;

\c joaopoliglota;

create table if not exists translations (
    translation_id bigserial primary key,
    idiom varchar(15),
    standard_key varchar(100),
    translation varchar(400)
);