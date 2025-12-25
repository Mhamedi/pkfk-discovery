SELECT 
    COUNT(*) as total_rows,
    COUNT(DISTINCT $3) as distinct_values,
    COUNT(*) - COUNT($3) as null_count,
    MIN($3) as min_value,
    MAX($3) as max_value,
    AVG(CASE WHEN $3 IS NOT NULL THEN CAST($3 AS NUMERIC) ELSE NULL END) as avg_value
FROM $1.$2
LIMIT $4;

