# liquipedia-calendar

### What is it ?

liquipedia-calendar is a web scrapper / iCalendar server which uses https://liquipedia.net/starcraft2/Liquipedia:Upcoming_and_ongoing_matches as a source and Flask as a iCalendar server.

The end result is a customizable link which permits to subscribe to arbitrary number of players and teams, pro or not as long as it figures on Liquipedia webpage.

You can subscribe to a iCal feed using this [link](https://napolitain.github.io/liquipedia-calendar/)

Example link for Starcraft 2 with Maru and Serral players : https://liquipedia-calendar.oa.r.appspot.com/?query=673d7374617263726166743226703d4d6172752c53657272616c
![image](https://user-images.githubusercontent.com/18146363/134248454-f5817f99-e780-431f-b56d-20a8c4d3dbfd.png)

Fill the querystring by visiting : https://liquipedia.net/

Once the link added in a Calendar App, events are auto generated and look like this.
<img width="766" src="https://user-images.githubusercontent.com/18146363/134247169-57a25f93-66bd-47fd-906e-38641afe084d.png">

### Features supported
* Every e sports games from liquipedia.net (having a dedicated matches page).
* Auto updating calendar integrated into standard calendar apps (such as Apple Calendar, Google Calendar, Microsoft Outlook).
* Specific players, games, teams

### Features being worked on
* A static website to easily make an URL for subscribing (right now, you must fill querystring parameter with some hacks).

### Technical

See [technical architecture here](https://github.com/Napolitain/liquipedia-calendar/blob/master/DESIGN.md)
