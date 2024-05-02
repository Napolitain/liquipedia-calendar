'use client';

import React from "react";

export default function Home() {
  return (
      <div>
          {/* Top part is an input accepting URLs
           * Entering resets the input
           * Top part is an input accepting URLs
           * Entering resets the input
           */}
          <div className={"part"}>
            <label>Add a player (append link as below, and press enter or space)</label>
            <input placeholder={"https://liquipedia.net/starcraft2/Serral"} autoFocus={true} onKeyDown={onKeyDown} type="text"/>
          </div>
          <div className={"part"}>
            <label>Tracked players (click on element to remove)</label>
              <span id={"tracked-players"}></span>
          </div>
          {/* Bottom part is the output link which, when clicked, is copied to the clipboard.*/}
          <div className={"part"}>
              <label>Output link (automatically copied to clipboard)</label>
              <a href={"https://calendar.google.com/calendar/u/0/r/settings/addbyurl?pli=1"} id={"output-link"}>https://liquipedia-calendar.oa.r.appspot.com/?query=</a>
          </div>
          <div className={"part"}>
              <label id={"final-redirect"} onClick={redirectToGoogleCalendar}>Append to calendar (click)</label>
          </div>
      </div>
  );
}

let trackedPlayers: Record<string, Set<string>> = {}; // Game -> Set of players

/**
 * On space or enter key press, add the URL to the list of tracked players and copy to clipboard.
 */
function onKeyDown(event: React.KeyboardEvent<HTMLInputElement>) {
    /**
     * Extracts the game and player from the input URL.
     * @returns {game, player}
     * @example https://liquipedia.net/starcraft2/Serral -> {game: "starcraft2", player: "Serral"}
     */
    function getQuery(): { game: string; player: string } {
        // Add the URL to the list of tracked players
        const input = event.target as HTMLInputElement;
        const url = input.value;
        // Validate the URL
        if (!validateURL(url)) {
            throw new Error("Invalid URL");
        }
        input.value = "";
        // Get game and player value from the URL
        // Split by /, [3] is the game, [4] is the player
        const parts = url.split("/");
        const game = parts[3];
        const player = parts[4];
        return {game, player};
    }

    function addToTrackedPlayers(game: string, player: string): void {
        // Append to the list of tracked players
        if (!trackedPlayers[game]) {
            trackedPlayers[game] = new Set();
        }
        trackedPlayers[game].add(player);
    }

    function updateTrackedPlayersView() {
        // Update the tracked players input
        const trackedPlayersElement: HTMLElement = document.getElementById("tracked-players");
        trackedPlayersElement.innerHTML = "";
        for (const game in trackedPlayers) {
            const g = trackedPlayers[game];
            g.forEach(player => {
                // Add a sticker game:player to the input
                // If there are multiple players, separate them by commas
                const span = document.createElement("span");
                span.className = "tracked-player";
                span.id = "tracked-player-" + game + "-" + player;
                span.onclick = () => {
                    // Remove the player from the list of tracked players
                    trackedPlayers[game].delete(player);
                    // Update the tracked players input
                    document.getElementById(span.id).remove();
                    // Update the output link
                    updateOutputLink();
                }
                span.textContent = game + ":" + player;
                trackedPlayersElement.appendChild(span);
            })
        }
    }

    if (event.key === "Enter" || event.key === " ") {
        let game: string, player: string;
        try {
            ({game, player} = getQuery());
        } catch (e) {
            console.error(e);
            return;
        }
        addToTrackedPlayers(game, player);
        updateTrackedPlayersView();
        // Update the output link
        updateOutputLink();
    }
}

/**
 * Updates the output link with the current list of tracked players.
 * The URL (CalDAV) is copied to the clipboard.
 * The user is redirected to Google Calendar.
 */
function updateOutputLink() { //  g=starcraft2&p=maru,serral
    const outputLink = document.getElementById("output-link") as HTMLAnchorElement;
    let p = "https://liquipedia-calendar.oa.r.appspot.com/?query=";
    for (const game in trackedPlayers) {
        const g = trackedPlayers[game];
        let n = g.size;
        if (n === 0) {
            continue;
        }
        p += stringToHex("g=" + game + "&p=");
        let i = 0;
        g.forEach(player => {
            p += stringToHex(player);
            if (i < n - 1) {
                p += stringToHex(",");
            }
            i++;
        })
    }
    outputLink.textContent = p;
    // Copy to clipboard
    navigator.clipboard.writeText(p);
}

/**
 * Converts a string to a hexadecimal string.
 * @param str The input string.
 * @returns {string} The hexadecimal string.
 * @example "Serral" -> "53757272616c"
 */
function stringToHex(str: string) {
    let hex = '';
    for (let i = 0; i < str.length; i++) {
        hex += '' + str.charCodeAt(i).toString(16);
    }
    return hex;
}

/**
 * Redirects the user to Google Calendar.
 */
function redirectToGoogleCalendar() {
    // Redirect to Google Calendar
    const outputLink = document.getElementById("output-link") as HTMLAnchorElement;
    outputLink.click();
}

/**
 * Validates the input URL.
 * @param input The input URL.
 * @returns {boolean} True if the URL is valid, false otherwise.
 * @example https://liquipedia.net/starcraft2/Serral -> true
 * @example https://liquipedia.net/starcraft2//Serral -> false
 * @example https://liquipedia.net/starcraft2/Serral/ -> false
 */
function validateURL(input: string): boolean {
    let regex = /https:\/\/liquipedia\.net\/[A-Za-z0-9]+\/[A-Za-z]+/i;
    return regex.test(input);
}
