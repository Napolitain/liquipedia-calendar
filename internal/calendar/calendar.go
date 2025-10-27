package calendar

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/Napolitain/liquipedia_calendar/internal/handler"
	"github.com/PuerkitoBio/goquery"
	ics "github.com/arran4/golang-ical"
)

var logger = slog.Default()

// CreateCalendar creates an iCalendar from a goquery document and player query
func CreateCalendar(document *goquery.Document, player handler.Query) (*ics.Calendar, error) {
	// Create iCalendar
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//liquipedia-calendar//en")
	cal.SetVersion("2.0")
	cal.SetLastModified(time.Now())
	// Create events
	matches := document.Find(".infobox_matches_content")
	var UIDs []string
	for i := 0; i < matches.Size(); i++ {
		// Get event info
		teamleft := matches.Eq(i).Find(".team-left a").Eq(0)
		teamright := matches.Eq(i).Find(".team-right span:not(.flag):not(.team-template-image):not(.team-template-team-short) a").Eq(0)
		teamleft_text := teamleft.Text() // TODO: test all match cases
		teamright_text := teamright.Text()
		teamleft_title := teamleft.AttrOr("title", "")
		teamright_title := teamright.AttrOr("title", "")
		contains := false
		for _, p := range player.Players {
			// Filter out events that don't contain at least one player in players
			if teamleft_text == p || teamright_text == p || teamleft_title == p || teamright_title == p {
				contains = true
				break
			}
		}
		if contains == false {
			continue
		}
		matchFormat := matches.Eq(i).Find(".versus abbr").Eq(0).Text()
		timestampStr, exist := matches.Eq(i).Find("[data-timestamp]").Eq(0).Attr("data-timestamp")
		if exist != true {
			return nil, errors.New("Timestamp doesn't exist.")
		}
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return nil, errors.New("Error while converting string to int.")
		}
		tournament := matches.Eq(i).Find(".match-filler div div a").Eq(0).Text()
		// Create a unique and RFC 5545 compliant UID using SHA256 hash
		// Hash the combination of timestamp, teams, and tournament to ensure uniqueness
		// This avoids spaces and special characters in the UID
		uidComponents := timestampStr + "-" + teamleft_text + "-" + teamright_text + "-" + tournament
		hash := sha256.Sum256([]byte(uidComponents))
		uid := hex.EncodeToString(hash[:]) + "@lcalendar"
		flag := false // Ignore identical UIDs for now
		for j := 0; j < len(UIDs); j++ {
			if UIDs[j] == uid {
				flag = true
				break
			}
		}
		if flag {
			continue
		}
		UIDs = append(UIDs, uid)
		// Add event
		event := cal.AddEvent(uid)
		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetStartAt(time.Unix(timestamp, 0))
		event.SetEndAt(time.Unix(timestamp+3600, 0))
		event.SetSummary(teamleft_text + " - " + teamright_text)
		event.SetLocation(tournament + " (" + matchFormat + ")")
	}
	return cal, nil
}
