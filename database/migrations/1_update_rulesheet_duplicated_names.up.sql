UPDATE rulesheets AS r SET name = concat(r.name, concat("-deleted-", r.id)) WHERE r.deleted_at IS NOT NULL and name not like "%-deleted-%";
