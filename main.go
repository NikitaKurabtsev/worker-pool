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

		for _, user := range users {
			inputCh <- user
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

	// here we "consume" data to the []User slice
	for user := range outputCh {
		if user.err != nil {
			log.Printf("%v\n", user.err)
			continue
		}
		fmt.Println(user)
	}
}

// addPhone returns user after mapping the username to
// the phones.json keys to fetch user phone number
func addPhone(user User, phones map[string]string) (User, error) {
	time.Sleep(1 * time.Second)

	if userPhone, ok := phones[user.Username]; ok {
		user.phone = userPhone
		return user, nil
	}

	return user, fmt.Errorf("cannot find phone number for user ID: %s", user.ID)
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
		var err error
		user, err = addPhone(user, phones)

		outputCh <- UserOutput{
			User: user,
			err:  err,
		}
	}
}
