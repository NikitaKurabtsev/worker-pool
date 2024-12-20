package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	phone    string
}

type UserOutput struct {
	User
	err error
}

func main() {
	// fetch data from json files
	users, phones, err := openDataset()
	if err != nil {
		log.Fatal("cannot open testing dataset: %w", err)
	}

	const workerCount = 3

	inputCh := make(chan User)
	outputCh := make(chan UserOutput)
	wg := sync.WaitGroup{}

	// here we "produce" data
	go func() {
		defer close(inputCh)

		for user := range users {
			inputCh <- users[user]
		}
	}()

	// here we process and send data from input channel
	// to the output channel
	wg.Add(workerCount)
	for w := 1; w <= workerCount; w++ {
		go processUsers(&wg, inputCh, outputCh, phones)
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	outputUsers := make([]User, 0)

	// here we "consume" data to the []User slice
	for res := range outputCh {
		outputUsers = append(outputUsers, res.User)
	}

	for i := range outputUsers {
		fmt.Println(outputUsers[i])
	}

}

// addPhone returns user after mapping the username to
// the phones.json keys to fetch user phone number
func addPhone(user User, phones map[string]string) User {
	time.Sleep(1 * time.Second)
	user.phone = phones[user.Username]

	return user
}

// processUsers use addPhone to add phone number
// to the user struct and send result to the channel
func processUsers(
	wg *sync.WaitGroup,
	inputCh <-chan User,
	outputCh chan<- UserOutput,
	phones map[string]string,
) {
	defer wg.Done()

	for user := range inputCh {
		user = addPhone(user, phones)

		outputCh <- UserOutput{
			User: user,
			err:  nil,
		}
	}
}
