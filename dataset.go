package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func openDataset() (users []User, phones map[string]string, err error) {
	usersBytes, err := os.ReadFile("users.json")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read users: %w", err)
	}

	usersPhones, err := os.ReadFile("phones.json")
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read phones: %w", err)
	}

	err = json.Unmarshal(usersBytes, &users)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot umarshal users: %w", err)
	}

	phones = map[string]string{}
	err = json.Unmarshal(usersPhones, &phones)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot unmarshal phones: %w", err)
	}

	return users, phones, nil
}
