create table alexstorm_shockwaves.execution
(
    id        integer not null
        constraint execution_pk
            primary key,
    params    json,
    container text
);