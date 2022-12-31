create table if not exists users (
	username        varchar primary key ,
	password        varchar             ,
	role            varchar             , 
	timestamp       timestamp
);

create table if not exists posts (
    id              serial  primary key ,
    title           varchar not null    ,
    content         text    not null    ,
    timestamp       timestamp
);

insert into users(username, password, role, timestamp) 
values('admin', '240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9', 'admin', NOW())
on conflict do nothing;

insert into users(username, password, role, timestamp) 
values('user', 'e606e38b0d8c19b24cf0ee3808183162ea7cd63ff7912dbb22b5e803286b4446', 'user', NOW())
on conflict do nothing;