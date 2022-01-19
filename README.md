# liquipedia-calendar

### What is it ?

I made initially this project in Python / Flask / Redis, but eventually ran into maintenance problems.
You can look nonetheless at its repository.

Link : https://github.com/Napolitain/sc2-calendar

Now, written in Go, and using Memcached mostly because of its free tier, the application is more maintenable and faster.
It is also a cool way to learn Go 😁

You can subscribe to a iCal feed using this link : https://liquipedia-calendar.oa.r.appspot.com/?game=
Fill the querystring by visiting : https://liquipedia.net/

### Features supported
* Every e sports games from liquipedia.net (having a dedicated matches page).
* Auto updating calendar integrated into standard calendar apps (such as Apple Calendar, Google Calendar, Microsoft Outlook).

### Features being worked on
* Specific players, specific tournament
* A static website to easily make an URL for subscribing (right now, you must fill querystring parameter with some hacks).
