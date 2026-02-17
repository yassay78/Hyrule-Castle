package main

import (
	"fmt"
	"os"
)

// Def de la structure dun personnage
type Character struct {
	Name   string
	Hp     int
	MaxHp  int
	Str    int
}

func main() {
	// Initialisation de Link
	// (j'ai pris les stats du sujet)
	link := Character{
		Name:  "Link",
		Hp:    60,
		MaxHp: 60,
		Str:   15,
	}

