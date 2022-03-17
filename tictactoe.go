package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	PlayerOneColor = "\033[1;34m%s\033[0m"
	PlayerTwoColor = "\033[1;33m%s\033[0m"
	InfoColor      = "\033[1;36m%s\033[0m"
	ErrorColor     = "\033[1;31m%s\033[0m"
)

type tictacboard [3][3]rune

var currentPlayer string

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	color.Blue("*UNLTIMATE TIC TAC TOE*\n enter 1 to login: , Enter 2 to see rules: , Enter 3 to create login: ")
	scanner.Scan()
	menuinput := scanner.Text()
	switch menuinput {
	case "1":
		//do you think that this should be recatored even more or does this seem fluid?
		color.Blue("Please enter your username: ")
		scanner.Scan()
		username := scanner.Text()
		_, p, err := getUserByUserName(username)
		if err != nil {
			color.Red("failed to get username: %v", err)
			return
		}
		color.Blue("please enter your password: ")
		scanner.Scan()
		password := scanner.Text()
		if password != p {
			color.Red("incorrect password: %v", err)
			return
		}
		//color.Blue("Enter 1 to start game: \nEnter 2 to look at game log: ")

	case "2":
		color.Blue("Please enter a new username: \n")
		scanner.Scan()
		newusername := scanner.Text()
		_, _, err := getUserByUserName(newusername)
		if err != nil {
			if errors.Is(err, io.EOF) {
				color.Blue("Please enter a new password: \n")
				scanner.Scan()
				newpassword := scanner.Text()
				createNewUser(newusername, newpassword)
				return
			}
			color.Red("server error: %v", err)
			return
		}

	case "3":
		// Question??: is this something that is worth refactoring down where it is already so small or is it still a good idea for readability? if so how?
		rules, err := ioutil.ReadFile("rules.txt")
		if err != nil {
			color.Red("File Read Error ")
			return
		}
		fmt.Println(string(rules))
		return

	default:
		color.Red("NOT A VALID INPUT")
		return
	}

	rand.Seed(time.Now().UnixNano())

	var playerMove bool
	var whoWon string
	var win bool

	var board tictacboard

	color.Blue("Strating Game: Board Empty\n")

	board.displayBoard()

	if rand.Intn(2) == 0 {
		playerMove = true
	} else {
		playerMove = false
	}

	for i := 0; i < 9; i++ {
		if playerMove {
			fmt.Println("Player Move: ", i+1)
			time.Sleep(time.Second)
			board.player()
			playerMove = false
		} else {
			fmt.Println("Computer Move: \n", i+1)
			time.Sleep(time.Second)
			board.computer()
			playerMove = true
		}

		if whoWon, win = board.check(); win {
			break
		}
		board.displayBoard()
	}

	//WinMessage := fmt.Sprintf
	color.HiGreen("*****%v won*****\nFinal Board View:\n", whoWon)
	//fmt.Printf(InfoColor, WinMessage)

}

func (t *tictacboard) displayBoard() {
	fmt.Print("-------------")
	for i := 0; i < 3; i++ {
		fmt.Print(`\n|`)
		for j := 0; j < 3; j++ {
			fmt.Printf(" %c |", t[i][j])
		}
		fmt.Printf("\n-------------")
	}
	fmt.Print("\n")
}

func (t *tictacboard) player() {
	var x, y int

	color.Blue("Enter the Row(1-3 and the Column(1-3 postions: ")
	if _, err := fmt.Scan(&x, &y); err == nil {
		x--
		y--
		if (x >= 0 && x <= 2) && (y >= 0 && y <= 2) && (t[x][y] == 0) {
			t[x][y] = 'x'
		} else {
			color.Red("Invalid input or position not empty. Try Again\n")
			t.player()
		}
	} else {
		color.Red("Invalid input or position not empty. Try Again\n")
		t.player()
	}
}

func (t *tictacboard) computer() {
	var x, y int
	for {
		x = rand.Intn(3)
		y = rand.Intn(3)
		if t[x][y] == 0 {
			t[x][y] = '0'
			break
		}
	}
}

func (t *tictacboard) check() (string, bool) {
	for i := 0; i < 3; i++ {
		if (rune(t[i][0]) == 'x') && (t[i][0] == t[i][1] && t[i][0] == t[i][2]) {
			err := recordGame(currentPlayer, true)
			if err != nil {
				panic(err)
			}
			return currentPlayer, true
		} else if (rune(t[i][0]) == '0') && (t[i][0] == t[i][1]) && (t[i][0] == t[i][2]) {
			err := recordGame(currentPlayer, false)
			if err != nil {
				panic(err)
			}
			return "Computer", true
		}

	}

	for i := 0; i < 3; i++ {
		if (rune(t[0][i]) == 'x') && (t[0][i] == t[1][i]) && (t[0][i] == t[2][i]) {
			err := recordGame(currentPlayer, true)
			if err != nil {
				panic(err)
			}
			return currentPlayer, true
		} else if (rune(t[0][i]) == '0') && (t[0][i] == t[1][i]) && (t[0][i] == t[2][i]) {
			err := recordGame(currentPlayer, false)
			if err != nil {
				panic(err)
			}
			return "Computer", true
		}
	}

	if ((rune(t[0][0]) == 'x') && (t[0][0] == t[1][1] && t[1][1] == t[2][2])) || ((rune(t[0][2]) == 'x') && (t[0][2] == t[1][1]) && (t[1][1] == t[2][0])) {
		err := recordGame(currentPlayer, true)
		if err != nil {
			panic(err)
		}
		return currentPlayer, true
	} else if ((rune(t[0][0]) == '0') && (t[0][0] == t[1][1]) && (t[1][1] == t[2][2])) || ((rune(t[0][2]) == '0') && (t[0][2] == t[1][1]) && (t[1][1] == t[2][0])) {
		err := recordGame(currentPlayer, false)
		if err != nil {
			panic(err)
		}
		return "Computer", true
	}

	//TODO: handle draws.
	return "No one", false

}

func recordGame(player string, isWinner bool) error {
	f, err := os.OpenFile("./playerdb.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		OpenFileError := fmt.Errorf("file did not open: %w", err)
		return fmt.Errorf(ErrorColor, OpenFileError)
	}
	w := csv.NewWriter(f)
	record := []string{player, fmt.Sprintf("%t", isWinner)}
	log.Printf("game played %v", record)
	err = w.Write(record)
	if err != nil {
		GameLogError := fmt.Errorf("game was unable to be recorded. %w", err)
		return fmt.Errorf(ErrorColor, GameLogError)
	}
	w.Flush()
	return w.Error()
}

func getUserByUserName(value string) (string, string, error) {
	f, err := os.OpenFile("./usernamedb.csv", os.O_APPEND|os.O_RDONLY, os.ModeAppend)
	if err != nil {
		OpenFileError := fmt.Errorf("file did not open: %w", err)
		return "", "", fmt.Errorf(ErrorColor, OpenFileError)
	}
	defer f.Close()
	r := csv.NewReader(f)
	for {
		rec, err := r.Read()
		if err != nil {
			return "", "", fmt.Errorf("failed to read row: %w", err)
		}
		if rec[0] == value {
			return rec[0], rec[1], nil
		}
	}
}

func createNewUser(user string, pass string) (string, string, error) {
	f, err := os.OpenFile("./usernamedb.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		openFileError := fmt.Errorf("file did not open: %v", err)
		return "", "", openFileError
	}
	w := csv.NewWriter(f)
	ts := time.Now().String()

	record := []string{user, pass, ts}
	err = w.Write(record)
	if err != nil {
		UserLogError := fmt.Errorf("user was unable to be recorded. %w", err)
		return "", "", UserLogError
	}
	w.Flush()
	err = w.Error()
	if err != nil {
		FlushError := fmt.Errorf("error flushing the file, %w", err)
		fmt.Printf(ErrorColor, FlushError)
		return "", "", err
	}
	log.Printf("Username created %v", user)
	currentPlayer = user
	return "", "", err
}
