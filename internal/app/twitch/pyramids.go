package twitch

import (
	"math/rand"
	"regexp"
	"strings"
	"sync"

	"github.com/Adeithe/go-twitch/irc"
)

var (
	spaces *regexp.Regexp = regexp.MustCompile("\\s+")
	facts  []string       = []string{
		"sodaG Fact #1: Giraffes are really tall.",
		"sodaG Fact #2: Giraffes eat leaves.",
		"sodaG Fact #3: Giraffes only need to drink water once every couple of days. They get most of their water from their plant-based diet-which is good considering their height makes the process of drinking difficult (and, if a lion happens upon a drinking giraffe, even dangerous).",
		"sodaG Fact #4: Female giraffes often return to where they were born to give birth. Once there, their calves receive a rough welcome into the world, falling over five feet to the ground.",
		"sodaG Fact #5: Male giraffes have pps.",
		"sodaG Fact #6: Female giraffes have vagenes.",
		"sodaG Fact #7: Fortunately, baby giraffes can stand up and even run within a hour of being born.",
		"sodaG Fact #8: Giraffes only need to drink once every few days. Most of their water comes from all the plants they eat.",
		"sodaG Fact #9: Baby giraffes are cute.",
		"sodaG Fact #10: Giraffes' tongues can be up to 20 inches long and are darkly colored, which is thought to help protect them during frequent sun-exposure.",
		"sodaG Fact #11: An erect male giraffes pp is 122cm (4ft) long.",
		"sodaG Fact #12: Giraffes usually stay upright while sleeping and if they do settle into a vulnerable position on the ground, it's just for a quick six-minute nap.",
		"sodaG Fact #13: Giraffes are usually yellow.",
		"sodaG Fact #14: Giraffes require over 75 pounds of food a day-and with a diet of leaves, this means they spend most of their time eating.",
		"sodaG Fact #15: Both male and female giraffes have two distinct, hair-covered horns called ossicones. Male giraffes use their horns to sometimes fight with other males.",
		"sodaG Fact #16: The giraffe's scientific name, Giraffa camelopardalis, comes from the ancient Greeks' belief that it looked like a camel wearing a leopard's coat.",
		"sodaG Fact #17: Despite their characteristic long necks, giraffes actually have the same number of neck vertebrae as humans-just seven. Each individual vertebra is super-sized, measuring up to ten inches long.",
		"sodaG Fact #18: Because of their unusual shape, giraffes have a highly-specialized cardiovascular system that starts with an enormous heart. It's two feet long and weighs up to 25 pounds.",
		"sodaG Fact #19: Giraffes are currently an endangered species.",
		"sodaG Fact #20: Male giraffes engage in a ritualized display of dominance called 'necking' that involves head-butting each other's bodies.",
		"sodaG Fact #21: Giraffes walk by moving both legs on the same side of their body together.",
		"sodaG Fact #22: A swift kick from a giraffes long legs can do serious damage to, or even kill, an unlucky lion.",
		"sodaG Fact #23: Male giraffes will test a female's fertility by tasting her urine.",
		"sodaG Fact #24: June 21, 2014 was be the first ever World Giraffe Day.",
		"sodaG Fact #25: Giraffes live roughly 25 years in the wild.",
		"sodaG Fact #26: Lorenzo de' Medici was gifted a giraffe by the sultan of Egypt. Giraffes had not been seen in Italy since antiquity and it caused quite the sensation, wandering the streets of Florence and accepting treats offered out of second-story windows.",
		"sodaG Fact #27: Female giraffes can become pregnant at 5 years old.",
		"sodaG Fact #28: The average height of a giraffe is around 5m (16-18ft)",
		"sodaG Fact #29: Giraffes sleep less than 2 hours per day.",
		"sodaG Fact #30: Over short distances, giraffes can run at speeds up to 35 mph.",
		"sodaG Fact #31: Giraffes are evolved horses.",
		"sodaG Fact #32: 90% of male giraffes mate with other males. KappaPride Clap",
	}

	msgCount   map[int]int    = make(map[int]int)
	lastFact   map[int]string = make(map[int]string)
	lastSender map[int]int    = make(map[int]int)
	tiers      map[int]int    = make(map[int]int)
	mx         sync.Mutex
)

func getPyramidData(roomID int) (string, int, int, int) {
	mx.Lock()
	defer mx.Unlock()
	return lastFact[roomID], lastSender[roomID], tiers[roomID], msgCount[roomID]
}

func setPyramidData(roomID int, fact string, sender, tierCount, count int) {
	mx.Lock()
	defer mx.Unlock()
	lastFact[roomID] = fact
	lastSender[roomID] = sender
	tiers[roomID] = tierCount
	msgCount[roomID] = count
}

// This is a terrible implementation. I'll rewrite later.
func handlePyramids(msg irc.ChatMessage) {
	fact, lastSender, tiers, count := getPyramidData(msg.ChannelID)
	if lastSender != msg.Sender.UserID {
		lastSender = msg.Sender.UserID
		tiers = 0
		count = 0
	}
	i := len(strings.Split(spaces.ReplaceAllString(msg.Message, " "), " "))
	count++

	if !msg.Sender.IsModerator && count > 3 && tiers > 2 && tiers-1 == i {
		fact = getGiraffeFact(fact)
		bot.Send(msg.Channel, fact)
	}
	tiers = i
	setPyramidData(msg.ChannelID, fact, lastSender, tiers, count)
}

func getGiraffeFact(last string) string {
	fact := facts[rand.Intn(len(facts))]
	if fact == last {
		return getGiraffeFact(last)
	}
	return fact
}
