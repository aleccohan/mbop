-- so we can generate random uuids.
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists public.registrations
(
    id     uuid  default uuid_generate_v4() not null
        constraint registrations_pk
            primary key,
    org_id varchar                         not null,
    uid    varchar                         not null,
    extra  jsonb default '{}'::jsonb,
    created_at timestamp default now() not null,
    updated_at timestamp default now() not null
);

-- unique index on org_id + uid
create unique index if not exists registrations_org_id_uid_uindex
    on public.registrations (org_id, uid);

-- auto-update the `updated_at` column on update
create function
    update_updated_at()
    returns trigger language plpgsql as $$
begin
    new.updated_at = now();
    return new;
end;
$$;

create trigger registration_updated
    before update on public.registrations
    for each row execute function update_updated_at();
