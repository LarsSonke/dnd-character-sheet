package storage

import (
	"encoding/json"
	"io/ioutil"
)

type Character struct {
	Name  string `json:"name"`
	Race  string `json:"race"`
	Class string `json:"class"`
	Level int    `json:"level"`
	Stats Stats  `json:"stats"`
}

type Stats struct {
	Str int `json:"str"`
	Dex int `json:"dex"`
	Con int `json:"con"`
	Int int `json:"int"`
	Wis int `json:"wis"`
	Cha int `json:"cha"`
}

func SaveCharacter(character Character, filename string) error {
	data, err := json.Marshal(character)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadCharacter(filename string) (Character, error) {
	var character Character
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return character, err
	}
	err = json.Unmarshal(data, &character)
	return character, err
}