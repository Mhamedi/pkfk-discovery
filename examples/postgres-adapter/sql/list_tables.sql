SELECT 
    table_schema,
    table_name,
    table_type
FROM information_schema.tables
WHERE table_schema = COALESCE($1, 'public')
    AND table_type = 'BASE TABLE'
ORDER BY table_schema, table_name;

