# liquipedia-calendar

### Upgraded version of sc2-calendar

I made initially this project in Python / Flask / Redis, but eventually ran into maintenance problems.
You can look nonetheless at its repository.

Link : https://github.com/Napolitain/sc2-calendar

Now, written in Go, and using Memcached mostly because of its free tier, the application is more maintenable and faster.
It is also a cool way to learn Go 😁

### What is it ?

What is it ?
sc2-calendar is a web scrapper / iCalendar server which uses https://liquipedia.net/starcraft2/Liquipedia:Upcoming_and_ongoing_matches as a source and Flask as a iCalendar server.

The end result is a customizable link which permits to subscribe to arbitrary number of players and teams, pro or not as long as it figures on Liquipedia webpage.

You can subscribe to a iCal feed using this link : https://liquipedia-calendar.oa.r.appspot.com/?game=
![image](https://user-images.githubusercontent.com/18146363/134248454-f5817f99-e780-431f-b56d-20a8c4d3dbfd.png)

Fill the querystring by visiting : https://liquipedia.net/

Once the link added in a Calendar App, events are auto generated and look like this.
<img width="766" src="https://user-images.githubusercontent.com/18146363/134247169-57a25f93-66bd-47fd-906e-38641afe084d.png">

### Features supported
* Every e sports games from liquipedia.net (having a dedicated matches page).
* Auto updating calendar integrated into standard calendar apps (such as Apple Calendar, Google Calendar, Microsoft Outlook).

### Features being worked on
* Specific players, specific tournament
* A static website to easily make an URL for subscribing (right now, you must fill querystring parameter with some hacks).

