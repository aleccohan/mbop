alter table registrations
    add constraint uid_unique
        unique (uid);
