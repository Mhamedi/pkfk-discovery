SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = COALESCE($1, 'public')
    AND tablename = $2
ORDER BY indexname;

