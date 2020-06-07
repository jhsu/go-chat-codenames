package codenames

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-chat-bot/bot"
)

type color int

const (
	Blue color = iota
	Neutral
	Black
)

func (c *color) String() string {
	return [...]string{"Blue", "Neutral", "Black"}[*c]
}

type card struct {
	color color
	word  string
}

type game struct {
	spymaster    string
	players      []string
	words        []*card
	correctWords int
	inprogress   bool
}

func (g *game) pickSpyMaster() (spymaster string) {
	g.spymaster = g.players[rand.Intn(len(g.players))]
	var guessers []string
	for _, name := range g.players {
		if name != g.spymaster {
			guessers = append(guessers, name)
		}
	}
	g.players = guessers
	return g.spymaster
}

func (g *game) guess(words []string) (matches []string, err error) {
	for _, word := range words {
		for _, card := range g.words {
			if word == card.word {
				switch color := card.color; color {
				case Blue:
					matches = append(matches, word)
				case Neutral:
					return matches, nil
				case Black:
					// end the game
					g.inprogress = false
					return matches, errors.New("You lose")
				default:
				}
			}
		}
	}
	if len(matches) == g.correctWords {
		g.inprogress = false
	}

	return matches, nil
}

// global games store
var games = make(map[string]*game)

func filterNames(names []string) (validNames []string, err error) {
	valid := make([]string, 0)
	for _, name := range names {
		if strings.HasPrefix(name, "@") {
			valid = append(valid, strings.TrimPrefix(name, "@"))
		}
	}

	if len(valid) == 0 {
		err := errors.New("Please provide at least one person participating using '@name'")
		return nil, err
	}
	return valid, nil
}

func generateWords() (cards []*card, correctCount int) {
	cards = []*card{
		{color: Blue, word: "daisy"},
		{color: Neutral, word: "flower"},
		{color: Black, word: "plant"}}
	correctCount = 0
	for _, card := range cards {
		if card.color == Blue {
			correctCount++
		}
	}
	return cards, correctCount
}

func startGame(cmd *bot.Cmd) (msg string, err error) {
	words, correctWords := generateWords()
	players, _ := filterNames(cmd.Args)
	if len(players) == 0 {
		return "Try starting a game with a list of players", errors.New("not enough players")
	}

	g := &game{players: players, words: words, correctWords: correctWords, inprogress: true}
	games[cmd.Channel] = g
	spymaster := g.pickSpyMaster()
	instructions := "Please give a clue"
	playerList := strings.Join(g.players, ", ")

	hiddenCards := make([]string, 0, len(words))
	for _, card := range words {
		hiddenCards = append(hiddenCards, card.word)
	}

	message := "spymaster: " + spymaster + " , players: " + playerList + "\n" + instructions + "\n" + strings.Join(hiddenCards, ", ")

	return message, nil
}

func contains(players []string, name string) bool {
	for _, player := range players {
		if player == name {
			return true
		}
	}
	return false
}

func makeGuess(cmd *bot.PassiveCmd) (response string, err error) {
	g := games[cmd.Channel]
	// game has not started, do nothing or user is spymaster
	if g == nil || !g.inprogress || g.spymaster == cmd.User.Nick {
		return "", nil
	}

	words := strings.Split(cmd.Raw, " ")
	correct, err := g.guess(words)

	if err != nil {
		return err.Error(), nil
	}

	if !g.inprogress {
		return "Your win! You hav guessed " + fmt.Sprintf("%d", len(correct)) + " words!", nil
	}
	return "You guessed " + fmt.Sprintf("%d", len(correct)) + " out of " + fmt.Sprintf("%d", g.correctWords), nil
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	bot.RegisterCommand("codenames", "Starts a codenames game, picks a random spymaster", "codenames @you @me @other @them", startGame)
	bot.RegisterPassiveCommand("guess", makeGuess)
}

var wordList = `
Acne
Acre
Addendum
Advertise
Aircraft
Aisle
Alligator
Alphabetize
America
Ankle
Apathy
Applause
Applesauc
Application
Archaeologist
Aristocrat
Arm
Armada
Asleep
Astronaut
Athlete
Atlantis
Aunt
Avocado
Baby-Sitter
Backbone
Bag
Baguette
Bald
Balloon
Banana
Banister
Baseball
Baseboards
Basketball
Bat
Battery
Beach
Beanstalk
Bedbug
Beer
Beethoven
Belt
Bib
Bicycle
Big
Bike
Billboard
Bird
Birthday
Bite
Blacksmith
Blanket
Bleach
Blimp
Blossom
Blueprint
Blunt
Blur
Boa
Boat
Bob
Bobsled
Body
Bomb
Bonnet
Book
Booth
Bowtie
Box
Boy
Brainstorm
Brand
Brave
Bride
Bridge
Broccoli
Broken
Broom
Bruise
Brunette
Bubble
Buddy
Buffalo
Bulb
Bunny
Bus
Buy
Cabin
Cafeteria
Cake
Calculator
Campsite
Can
Canada
Candle
Candy
Cape
Capitalism
Car
Cardboard
Cartography
Cat
Cd
Ceiling
Cell
Century
Chair
Chalk
Champion
Charger
Cheerleader
Chef
Chess
Chew
Chicken
Chime
China
Chocolate
Church
Circus
Clay
Cliff
Cloak
Clockwork
Clown
Clue
Coach
Coal
Coaster
Cog
Cold
College
Comfort
Computer
Cone
Constrictor
Continuum
Conversation
Cook
Coop
Cord
Corduroy
Cot
Cough
Cow
Cowboy
Crayon
Cream
Crisp
Criticize
Crow
Cruise
Crumb
Crust
Cuff
Curtain
Cuticle
Czar
Dad
Dart
Dawn
Day
Deep
Defect
Dent
Dentist
Desk
Dictionary
Dimple
Dirty
Dismantle
Ditch
Diver
Doctor
Dog
Doghouse
Doll
Dominoes
Door
Dot
Drain
Draw
Dream
Dress
Drink
Drip
Drums
Dryer
Duck
Dump
Dunk
Dust
Ear
Eat
Ebony
Elbow
Electricity
Elephant
Elevator
Elf
Elm
Engine
England
Ergonomic
Escalator
Eureka
Europe
Evolution
Extension
Eyebrow
Fan
Fancy
Fast
Feast
Fence
Feudalism
Fiddle
Figment
Finger
Fire
First
Fishing
Fix
Fizz
Flagpole
Flannel
Flashlight
Flock
Flotsam
Flower
Flu
Flush
Flutter
Fog
Foil
Football
Forehead
Forever
Fortnight
France
Freckle
Freight
Fringe
Frog
Frown
Gallop
Game
Garbage
Garden
Gasoline
Gem
Ginger
Gingerbread
Girl
Glasses
Goblin
Gold
Goodbye
Grandpa
Grape
Grass
Gratitude
Gray
Green
Guitar
Gum
Gumball
Hair
Half
Handle
Handwriting
Hang
Happy
Hat
Hatch
Headache
Heart
Hedge
Helicopter
Hem
Hide
Hill
Hockey
Homework
Honk
Hopscotch
Horse
Hose
Hot
House
Houseboat
Hug
Humidifier
Hungry
Hurdle
Hurt
Hut
Ice
Implode
Inn
Inquisition
Intern
Internet
Invitation
Ironic
Ivory
Ivy
Jade
Japan
Jeans
Jelly
Jet
Jig
Jog
Journal
Jump
Key
Killer
Kilogram
King
Kitchen
Kite
Knee
Kneel
Knife
Knight
Koala
Lace
Ladder
Ladybug
Lag
Landfill
Lap
Laugh
Laundry
Law
Lawn
Lawnmower
Leak
Leg
Letter
Level
Lifestyle
Ligament
Light
Lightsaber
Lime
Lion
Lizard
Log
Loiterer
Lollipop
Loveseat
Loyalty
Lunch
Lunchbox
Lyrics
Machine
Macho
Mailbox
Mammoth
Mark
Mars
Mascot
Mast
Matchstick
Mate
Mattress
Mess
Mexico
Midsummer
Mine
Mistake
Modern
Mold
Mom
Monday
Money
Monitor
Monster
Mooch
Moon
Mop
Moth
Motorcycle
Mountain
Mouse
Mower
Mud
Music
Mute
Nature
Negotiate
Neighbor
Nest
Neutron
Niece
Night
Nightmare
Nose
Oar
Observatory
Office
Oil
Old
Olympian
Opaque
Opener
Orbit
Organ
Organize
Outer
Outside
Ovation
Overture
Pail
Paint
Pajamas
Palace
Pants
Paper
Paper
Park
Parody
Party
Password
Pastry
Pawn
Pear
Pen
Pencil
Pendulum
Penis
Penny
Pepper
Personal
Philosopher
Phone
Photograph
Piano
Picnic
Pigpen
Pillow
Pilot
Pinch
Ping
Pinwheel
Pirate
Plaid
Plan
Plank
Plate
Platypus
Playground
Plow
Plumber
Pocket
Poem
Point
Pole
Pomp
Pong
Pool
Popsicle
Population
Portfolio
Positive
Post
Princess
Procrastinate
Protestant
Psychologist
Publisher
Punk
Puppet
Puppy
Push
Puzzle
Quarantine
Queen
Quicksand
Quiet
Race
Radio
Raft
Rag
Rainbow
Rainwater
Random
Ray
Recycle
Red
Regret
Reimbursement
Retaliate
Rib
Riddle
Rim
Rink
Roller
Room
Rose
Round
Roundabout
Rung
Runt
Rut
Sad
Safe
Salmon
Salt
Sandbox
Sandcastle
Sandwich
Sash
Satellite
Scar
Scared
School
Scoundrel
Scramble
Scuff
Seashell
Season
Sentence
Sequins
Set
Shaft
Shallow
Shampoo
Shark
Sheep
Sheets
Sheriff
Shipwreck
Shirt
Shoelace
Short
Shower
Shrink
Sick
Siesta
Silhouette
Singer
Sip
Skate
Skating
Ski
Slam
Sleep
Sling
Slow
Slump
Smith
Sneeze
Snow
Snuggle
Song
Space
Spare
Speakers
Spider
Spit
Sponge
Spool
Spoon
Spring
Sprinkler
Spy
Square
Squint
Stairs
Standing
Star
State
Stick
Stockholder
Stoplight
Stout
Stove
Stowaway
Straw
Stream
Streamline
Stripe
Student
Sun
Sunburn
Sushi
Swamp
Swarm
Sweater
Swimming
Swing
Tachometer
Talk
Taxi
Teacher
Teapot
Teenager
Telephone
Ten
Tennis
Thief
Think
Throne
Through
Thunder
Tide
Tiger
Time
Tinting
Tiptoe
Tiptop
Tired
Tissue
Toast
Toilet
Tool
Toothbrush
Tornado
Tournament
Tractor
Train
Trash
Treasure
Tree
Triangle
Trip
Truck
Tub
Tuba
Tutor
Television
Twang
Twig
Twitterpated
Type
Unemployed
Upgrade
Vest
Vision
Wag
Water
Watermelon
Wax
Wedding
Weed
Welder
Whatever
Wheelchair
Whiplash
Whisk
Whistle
White
Wig
Will
Windmill
Winter
Wish
Wolf
Wool
World
Worm
Wristwatch
Yardstick
Zamboni
Zen
Zero
Zipper
Zone
Zoo
`
