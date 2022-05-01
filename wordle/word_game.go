package wordle

import (
	"fmt"
	"strings"

	"github.com/p2p-games/wordle/model"
	"github.com/pkg/errors"
)

const maxAttempts int = 5

type WordGame struct {
	PeerId string

	// to verify if the guess is correct
	Target *model.Word
	Salts  []string

	StateIdx int
	NextWord string

	AttemptedWords []string
	isCorrect      map[string][]bool

	guessMsgC chan guess
}

type guess struct {
	Guess    string
	Proposal string
}

// generate new game session
func NewWordGame(target *model.Word, sendMsgC chan guess) *WordGame {
	salts := GetSaltsFromWord(target)
	return &WordGame{
		Target:         target,
		Salts:          salts,
		StateIdx:       0, // start requesting the word
		AttemptedWords: make([]string, 0),
		isCorrect:      make(map[string][]bool),
		guessMsgC:      sendMsgC,
	}
}

func (w *WordGame) ComposeStateUI() string {
	var s string
	switch w.StateIdx {
	case 0:
		s += fmt.Sprintf("Introduce your word proposal for the next word to guess:\n")
	case 1:
		s = "Guess which is the current Word:\n"
		for _, guessedWord := range w.AttemptedWords {
			if guessedWord != "" {
				// check wheather the word was correct or not
				correct := "x"
				comp, err := model.VerifyString(guessedWord, w.Target)
				if err != nil {
					continue
				}
				if IsGuessSuccess(comp) {
					correct = "v"
				}
				// compose the color strings with color chars

				s += fmt.Sprintf("\t[%s] %s\n", correct, ComposeWordleVisualWord(guessedWord, w.Target))
			}
		}
		s += fmt.Sprintf("\nAttempts left %d\n", maxAttempts-len(w.AttemptedWords))
	case 2:
		s = "\n\tCongrats, you guessed the word!\nWait untill someone guesses it to play again\n"
	case 3:
		s = "\n\tNo more attempts left for this word!\nWait untill someone guesses it to play again\n"
	default:
		s = "unrecognized state to generate the UI\n"
	}
	return s
}

func (w *WordGame) NewStdinInput(input string) error {
	// check if non alphanumeric character
	input = strings.ToLower(input)
	if !IsLetter(input) {
		return nil
	}
	// check in which state do we are
	switch w.StateIdx {
	case 0:
		err := w.addNextTarget(input, w.Salts)
		if err != nil {
			return err
		}
	case 1:
		err := w.addNewGuess(input)
		if err != nil {
			return err
		}
	default:
		// nothing
		return errors.New("unable to add any new guess")
	}
	return nil
}

func (w *WordGame) addNextTarget(nextWord string, salts []string) error {
	// check if we are in state 0
	if w.StateIdx != 0 {
		return errors.New("unable to add next target, not in state 0")
	}

	w.NextWord = nextWord
	// go to state 1
	w.StateIdx = 1
	return nil
}

func (w *WordGame) WasGuessed() bool {
	// check if we have already guessed 5 times or if guess correct
	for _, word := range w.AttemptedWords {
		comp, err := model.VerifyString(word, w.Target)
		if err != nil {
			continue
		}
		if IsGuessSuccess(comp) {
			return true
		}
	}
	return false
}

func (w *WordGame) addNewGuess(guessedWord string) error {
	// check if we are in state 1
	if w.StateIdx != 1 {
		return errors.New("unable to add next target, not in state 1")
	}

	// add the new word to the list of Attempted, Verify if its the correct word, add to the map the result
	w.AttemptedWords = append(w.AttemptedWords, guessedWord)

	correct, err := model.VerifyString(guessedWord, w.Target)
	if err != nil {
		return nil
	}
	w.isCorrect[guessedWord] = correct

	comp, err := model.VerifyString(guessedWord, w.Target)
	if err != nil {
		return err
	}
	if IsGuessSuccess(comp) {
		w.StateIdx = 2 // Congrats, wait untill someone guesses your word
	}

	// check if we did all the attempts
	if len(w.AttemptedWords) == maxAttempts && !IsGuessSuccess(comp) {
		w.StateIdx = 3 // Wait untill you can play again
	}

	// send the msg over the channel to notify the service
	currentGuess := guess{
		guessedWord,
		w.NextWord,
	}
	w.guessMsgC <- currentGuess

	return nil
}
