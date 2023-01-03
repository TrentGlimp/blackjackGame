package main

import (
	"deck"
	"fmt"
	"strings"
)

// Simple minimum function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Hand declaration and functions
type Hand []deck.Card

// How to print a hand
func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

// How to print the dealers hand
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// How to score a hand
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

// Takes into account that aces can be 1 or 11 points
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Draws a card from the deck into your hand
func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// Defines a State and the values it can have
type State int8

const (
	StatePlayerTurn = iota // StatePlayerTurn = 0, StateDealerTurn = 1, StateHandOver = 2
	StateDealerTurn
	StateHandOver
)

// Defines GameStateStruct
type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

// Determines what phase of the game it is
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("It isn't currently any player's turn")
	}
}

// makes a clone of the current game (used to pass GameStates by value)
func clone(gs GameState) GameState {
	cloned := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(cloned.Deck, gs.Deck)
	copy(cloned.Player, gs.Player)
	copy(cloned.Dealer, gs.Dealer)
	return cloned
}

// Shuffles a deck of cards
func Shuffle(gs GameState) GameState {
	shuffled := clone(gs)
	shuffled.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return shuffled
}

// Deals out the hands to the player and dealer
func Deal(gs GameState) GameState {
	handsDealt := clone(gs)
	handsDealt.Player = make(Hand, 0, 5)
	handsDealt.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, handsDealt.Deck = draw(handsDealt.Deck)
		handsDealt.Player = append(handsDealt.Player, card)
		card, handsDealt.Deck = draw(handsDealt.Deck)
		handsDealt.Dealer = append(handsDealt.Dealer, card)
	}
	handsDealt.State = StatePlayerTurn
	return handsDealt
}

// Ends player's or dealer's turn
func Stand(gs GameState) GameState {
	playerTurnOver := clone(gs)
	playerTurnOver.State++
	return playerTurnOver
}

// Draws card from deck and adds it to player's or dealer's hand + checks if busted
func Hit(gs GameState) GameState {
	gainedCard := clone(gs)
	hand := gainedCard.CurrentPlayer()
	var card deck.Card
	card, gainedCard.Deck = draw(gainedCard.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(gainedCard)
	}
	return gainedCard
}

// Ends a hand and checks to see who won
func EndHand(gs GameState) GameState {
	finishedHands := clone(gs)
	pScore, dScore := finishedHands.Player.Score(), finishedHands.Dealer.Score()
	var message string
	switch {
	case pScore > 21:
		message = "You busted"
	case dScore > 21:
		message = "Dealer busted"
	case pScore > dScore:
		message = "You win!"
	case dScore > pScore:
		message = "You lose"
	case pScore == dScore:
		message = "Draw"
	default:
		message = "********ERROR********"
	}

	fmt.Println("\n=====FINAL HANDS=====")
	fmt.Println("Player:", finishedHands.Player, "\nScore:", pScore)
	fmt.Println("Dealer:", finishedHands.Dealer, "\nScore:", dScore)
	fmt.Println(message)
	fmt.Println()

	finishedHands.Player = nil
	finishedHands.Dealer = nil
	return finishedHands
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player:", gs.Player)
		fmt.Println("Dealer:", gs.Dealer.DealerString())
		fmt.Println("What will you do?\n[hit/stand]")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "hit":
			gs = Hit(gs)
		case "stand":
			gs = Stand(gs)
		default:
			fmt.Println("Options are 'hit' or 'stand'\nPlease pick one of those two options")
		}
	}

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}
