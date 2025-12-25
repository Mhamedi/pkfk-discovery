SELECT 
    COUNT(*) as total_rows,
    COUNT(CASE WHEN $3 IN (SELECT $4 FROM $5.$6) THEN 1 END) as matching_rows
FROM $1.$2
LIMIT $7;

