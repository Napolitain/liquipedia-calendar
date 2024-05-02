# Benchmarks on GCP App Engine F1 instance (2024)

There are multiple cases :
1. instance state due to pod autoscaling
- cold start : 0 instance running, so need to start an instance
- warm start : 1+ instance running, so no need to start an instance
2. cache state
- no cache : cache is empty, so need to fetch data from Liquipedia.net
- cache (game) hit : cache contains the game, so no need to fetch data from Liquipedia.net
- cache (superstar player) hit : cache contains the superstar player, so the calendar is sent from Memcached directly.

Instance/Cache | No cache | Cache (game) hit | Cache (superstar player) hit
--- |----------|------------------| ---
Cold start | TBD      | TBD              | TBD
Warm start | TBD      | TBD              | TBD
