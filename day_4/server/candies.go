package calculator

import "errors"

type Candy struct {
	Name         string
	Abbreviation string
	Price        int64
}

type Candies []Candy

func (c Candies) Contains(candyType string) bool {
	for _, a := range c {
		if a.Abbreviation == candyType {
			return true
		}
	}
	return false
}

func (c Candies) GetPrice(candyType string) int64 {
	for _, a := range c {
		if a.Abbreviation == candyType {
			return a.Price
		}
	}
	return 0
}

var candies = Candies{
	Candy{"Cool Eskimo", "CE", 10},
	Candy{"Apricot Aardvark", "AA", 15},
	Candy{"Natural Tiger", "NT", 17},
	Candy{"Dazzling Elderberry", "DE", 21},
	Candy{"Yellow Rambutan", "YR", 23},
}

func GetChange(candyType string, total int64, count int64) (int64, error) {
	if !candies.Contains(candyType) || total <= 0 || count <= 0 {
		return 0, errors.New("")
	}
	change := total - count*candies.GetPrice(candyType)
	return change, nil
}
