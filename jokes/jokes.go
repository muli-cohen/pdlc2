package jokes

import "math/rand"

var corpus = [15]string{
	"Why do programmers prefer dark mode? Because light attracts bugs.",
	"Why did the scarecrow win an award? Because he was outstanding in his field.",
	"I told my wife she was drawing her eyebrows too high. She looked surprised.",
	"Why don't scientists trust atoms? Because they make up everything.",
	"What do you call a fake noodle? An impasta.",
	"Why did the bicycle fall over? Because it was two-tired.",
	"I used to hate facial hair, but then it grew on me.",
	"What do you call cheese that isn't yours? Nacho cheese.",
	"Why can't you give Elsa a balloon? Because she'll let it go.",
	"What did the ocean say to the beach? Nothing, it just waved.",
	"Why did the math book look so sad? It had too many problems.",
	"What do you call a bear with no teeth? A gummy bear.",
	"Why don't eggs tell jokes? Because they'd crack each other up.",
	"What do you call a sleeping dinosaur? A dino-snore.",
	"Why did the coffee file a police report? It got mugged.",
}

func All() []string {
	return corpus[:]
}

func Random(r *rand.Rand) string {
	return corpus[r.Intn(len(corpus))]
}
