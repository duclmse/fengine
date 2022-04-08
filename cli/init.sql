CREATE TABLE IF NOT EXISTS weather_metrics (
    time             TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    timezone_shift   INT                         NULL,
    city_name        TEXT                        NULL,
    temp_c           DOUBLE PRECISION            NULL,
    feels_like_c     DOUBLE PRECISION            NULL,
    temp_min_c       DOUBLE PRECISION            NULL,
    temp_max_c       DOUBLE PRECISION            NULL,
    pressure_hpa     DOUBLE PRECISION            NULL,
    humidity_percent DOUBLE PRECISION            NULL,
    wind_speed_ms    DOUBLE PRECISION            NULL,
    wind_deg         INT                         NULL,
    rain_1h_mm       DOUBLE PRECISION            NULL,
    rain_3h_mm       DOUBLE PRECISION            NULL,
    snow_1h_mm       DOUBLE PRECISION            NULL,
    snow_3h_mm       DOUBLE PRECISION            NULL,
    clouds_percent   INT                         NULL,
    weather_type_id  INT                         NULL
);

SELECT create_hypertable('weather_metrics', 'time');


COPY weather_metrics (
  time, timezone_shift, city_name, temp_c, feels_like_c, temp_min_c, temp_max_c, pressure_hpa,
  humidity_percent, wind_speed_ms, wind_deg, rain_1h_mm, rain_3h_mm, snow_1h_mm, snow_3h_mm,
  clouds_percent, weather_type_id
) FROM '/var/lib/postgresql/data/weather_data.csv' CSV HEADER;
