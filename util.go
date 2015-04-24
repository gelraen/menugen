package main

import (
	"math/rand"

	"./data"
)

const (
	kBreakfast = "breakfast"
	kLunch     = "lunch"
	kMeat      = "meat"
	kGarnish   = "garnish"
	kMain      = "main"
)

type Weekday int

const (
	Mon Weekday = iota
	Tue
	Wed
	Thu
	Fri
	Sat
	Sun
)

type DayMenu struct {
	Breakfast data.Dish
	Lunch     data.Dish
	Dinner    []data.Dish
}

type WeekMenu map[Weekday]DayMenu

func makeDishMap(dishes []data.Dish) map[string][]data.Dish {
	r := map[string][]data.Dish{}
	for _, d := range dishes {
		r[d.Type] = append(r[d.Type], d)
	}
	return r
}

func genRandomOfType(d map[string][]data.Dish, t string) data.Dish {
	return d[t][rand.Intn(len(d[t]))]
}

func genBreakfast(d map[string][]data.Dish) data.Dish {
	return genRandomOfType(d, kBreakfast)
}
func genLunch(d map[string][]data.Dish) data.Dish {
	return genRandomOfType(d, kLunch)
}

func genDinner(d map[string][]data.Dish) []data.Dish {
	switch rand.Intn(2) {
	case 0:
		return []data.Dish{genRandomOfType(d, kMeat), genRandomOfType(d, kGarnish)}
	case 1:
		return []data.Dish{genRandomOfType(d, kMain)}
	default:
		panic("rand.Intn(2) > 1")
	}
}
