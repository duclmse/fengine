docker run -dp 5432:5432 \
  -v ~/data/timescale:/var/lib/postgresql/data \
  -e TZ=Asia/Ho_Chi_Minh \
  -e POSTGRES_PASSWORD=1 \
  --name timescale \
  timescale/timescaledb:latest-pg12
