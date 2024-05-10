# liquipedia-calendar

Watch this video to introduce you to Liquipedia Calendar : https://www.youtube.com/watch?v=wwhYcQSVFdA

### What is it ?

liquipedia-calendar is a web scrapper / iCalendar server which uses https://liquipedia.net/starcraft2/Liquipedia:Upcoming_and_ongoing_matches as a source and Go as a iCalendar server.

The end result is a customizable link which permits to subscribe to arbitrary number of players and teams, pro or not as long as it figures on Liquipedia webpage.

You can subscribe to a iCal feed using this [URL generator](https://napolitain.github.io/liquipedia-calendar/) and [liquipedia.net](https://liquipedia.net/)

Example link for Starcraft 2 with Maru and Serral players : https://liquipedia-calendar.oa.r.appspot.com/?query=673d7374617263726166743226703d4d6172752c53657272616c

Once the link added in a Calendar App, events are auto generated and look like this.
<img width="766" src="https://user-images.githubusercontent.com/18146363/134247169-57a25f93-66bd-47fd-906e-38641afe084d.png">

### Features supported
* Every e sports games from liquipedia.net (having a dedicated matches page).
* Auto updating calendar integrated into standard calendar apps (such as Apple Calendar, Google Calendar, Microsoft Outlook). Notifications handled by calendar application.
* Specific players, games, teams
* A static website written with Next.JS to easily make an URL for subscribing.

### Features being worked on
* For now, none

### Technical

See [technical architecture here](https://github.com/Napolitain/liquipedia-calendar/blob/master/DESIGN.md)
