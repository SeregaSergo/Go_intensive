package present

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetNCoolestPresents(t *testing.T) {
	presents1 := []Present{
		{4, 5},
		{3, 1},
		{5, 2},
		{5, 1},
	}
	presents2 := []Present{
		{4, 5},
		{3, 1},
		{5, 2},
		{5, 1},
		{8, 3},
		{8, 3},
	}
	fmt.Println("Testing GetNCoolestPresents...")
	{
		best, _ := GetNCoolestPresents(presents1, 2)
		if !reflect.DeepEqual(best, []Present{{5, 1}, {5, 2}}) {
			t.Error("Present1:", best)
		}
	}
	{
		best, _ := GetNCoolestPresents(presents2, 3)
		if !reflect.DeepEqual(best, []Present{{8, 3}, {8, 3}, {5, 1}}) {
			t.Error("Present2:", best)
		}
	}
	{
		_, err := GetNCoolestPresents(presents2, 10)
		if err == nil {
			t.Error("Present2:", err.Error())
		}
	}
}

func TestGrabPresents(t *testing.T) {
	presents1 := []Present{
		{4, 5},
		{3, 1},
		{5, 2},
		{5, 1},
	}
	presents2 := []Present{
		{4, 5},
		{3, 1},
		{5, 2},
		{5, 1},
		{8, 3},
		{8, 3},
	}
	fmt.Println("Testing GrabPresents...")
	{
		best := GrabPresents(presents1, 2)
		sum := 0
		for _, v := range best {
			sum += v.Value
		}
		if sum != 8 {
			t.Errorf("Present1 sum: %d\nNeed: %d\n", best, 8)
		}
	}
	{
		best := GrabPresents(presents2, 4)
		sum := 0
		for _, v := range best {
			sum += v.Value
		}
		if sum != 13 {
			t.Errorf("Present2 sum: %d\nNeed: %d\n", best, 13)
		}
	}
	{
		best := GrabPresents(presents2, 3)
		sum := 0
		for _, v := range best {
			sum += v.Value
		}
		if sum != 10 {
			t.Errorf("Present2 sum: %d\nNeed: %d\n", best, 10)
		}
	}
	{
		best := GrabPresents(presents2, 35)
		sum := 0
		for _, v := range best {
			sum += v.Value
		}
		if sum != 33 {
			t.Errorf("Present2 sum: %d\nNeed: %d\n", best, 33)
		}
	}
}
