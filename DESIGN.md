# System Design

### Google Cloud

The application is hosted on Google Cloud Platform (GCP). The application is a Go application running on Google App Engine (GAE), which also provides the Memcached service.
The application is served over HTTPS using Google Cloud Load Balancer (GCLB). The application is monitored using Google Cloud Monitoring (GCM) and Google Cloud Logging (GCL).

### Cache miss

When a user (Google Calendar) requests the Liquipedia Calendar, the application will first check if the data is in the cache (Memcached). If it is not, the application will fetch the data from Liquipedia.net and store it in the cache.
This will take some time as it depends synchronously on the response time of Liquipedia.net. Additionally, it exposes Liquipedia.net to the risk of being overwhelmed by requests from my own application.

```mermaid
sequenceDiagram
    participant Google Calendar
    box Liquipedia Calendar
    participant CalDAV
    participant Memcached
    end
    participant Liquipedia.net
    Google Calendar->>CalDAV: GET
    CalDAV->>Memcached: GET (SUPERSTAR PLAYER)
    Memcached-->>CalDAV: MISS
    CalDAV->>Memcached: GET (GAME)
    Memcached-->>CalDAV: MISS
    CalDAV->>Liquipedia.net: GET
    Liquipedia.net->>CalDAV: RESPONSE
    CalDAV->>Memcached: CREATE (GAME HTML)
    CalDAV-->>CalDAV: PARSE HTML, FILTER PLAYERS
    CalDAV-->>CalDAV: CREATE CALENDAR
    CalDAV-->>Memcached: CREATE (SUPERSTAR PLAYER)
    CalDAV-->>Google Calendar: CALENDAR
```

### Cache hit (game)

Once a first request has been made, the data will be stored in the cache. Subsequent requests will be served from the cache, which is much faster than fetching the data from Liquipedia.net.
The data will expire after a set period of time, for now 3 hours. There will be 2 cache keys types: one for the superstar player and one for the game.
If the player requested isn't found in the cache, but the game is, the application will receive HTML data from the cache.
The application will then filter players from the retrieved HTML data from the cache, and create the calendar.
The cache not only speeds up the application by removing external network calls, but also protects Liquipedia.net from being overwhelmed by requests.

```mermaid
sequenceDiagram
    participant Google Calendar
    box Liquipedia Calendar
        participant CalDAV
        participant Memcached
    end
    Google Calendar->>CalDAV: GET
    CalDAV->>Memcached: GET (SUPERSTAR PLAYER)
    Memcached-->>CalDAV: MISS
    CalDAV->>Memcached: GET (GAME)
    Memcached-->>CalDAV: RESPONSE (HTML)
    CalDAV-->>CalDAV: PARSE HTML, FILTER PLAYERS
    CalDAV-->>CalDAV: CREATE CALENDAR
    CalDAV-->>Memcached: CREATE (SUPERSTAR PLAYER)
    CalDAV-->>Google Calendar: CALENDAR
```

### Cache hit (superstar player)

If the player requested is found in the cache, the application will receive the calendar data from the cache. This is the fastest way to serve the data, as it doesn't require any external network calls,
and the data is already in the format required by Google Calendar. There is no HTML parsing, filtering or calendar creation required.

```mermaid
sequenceDiagram
    participant Google Calendar
    box Liquipedia Calendar
        participant CalDAV
        participant Memcached
    end
    Google Calendar->>CalDAV: GET
    CalDAV->>Memcached: GET (SUPERSTAR PLAYER)
    Memcached-->>CalDAV: RESPONSE (CALENDAR)
    CalDAV-->>Google Calendar: CALENDAR
```
