SELECT COUNT(*), PG_SIZE_PRETTY(PG_TOTAL_RELATION_SIZE('weather_metrics')) FROM weather_metrics;


SELECT
    PG_SIZE_PRETTY(PG_TOTAL_RELATION_SIZE(c.oid)) AS "total_size",
    PG_SIZE_PRETTY(PG_INDEXES_SIZE(c.oid)) AS "indexes_size", *
FROM pg_class c
         LEFT JOIN pg_namespace n ON (n.oid = c.relnamespace)
WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND c.relkind <> 'i' AND nspname !~ '^pg_toast'
ORDER BY PG_TOTAL_RELATION_SIZE(c.oid) DESC;


SELECT * FROM information_schema.tables WHERE table_schema = 'public';

-----------------------------------

SELECT city_name, AVG(temp_c) FROM weather_metrics WHERE time > NOW() - INTERVAL '2 years' GROUP BY city_name;


SELECT city_name, SUM(snow_1h_mm) sum_snow
FROM weather_metrics
WHERE time > NOW() - INTERVAL '5 years'
GROUP BY city_name;

SELECT time_bucket('15 days', time) AS "bucket", city_name, AVG(temp_c)
FROM weather_metrics
WHERE time > NOW() - (12 * INTERVAL '1 month')
GROUP BY bucket, city_name
ORDER BY bucket DESC;


