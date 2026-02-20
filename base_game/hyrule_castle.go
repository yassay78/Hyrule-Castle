package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Couleurs interface
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

type Stats struct {
	ID int; Name string; Strengths []int; Weaknesses []int
}

type Character struct {
	Name string `json:"name"`; Hp int `json:"hp"`; MaxHp int
	Str  int    `json:"str"`; Def int `json:"def"`; Spd int `json:"spd"`
	Race int    `json:"race"`; Class int `json:"class"`
}

var (
	Classes = make(map[int]Stats)
	Races   = make(map[int]Stats)
	Players, Enemies, Bosses []Character
)

func loadFiles() {
	read := func(f string, v interface{}) {
		d, _ := os.ReadFile("mods/" + f)
		json.Unmarshal(d, v)
	}
	var cL, rL []Stats
	read("classes.json", &cL); read("races.json", &rL)
	for _, c := range cL { Classes[c.ID] = c }; for _, r := range rL { Races[r.ID] = r }
	read("players.json", &Players); read("enemies.json", &Enemies); read("bosses.json", &Bosses)
	for i := range Players { Players[i].MaxHp = Players[i].Hp }
}

// barre de vie
func drawBar(hp, max int) string {
	percent := float64(hp) / float64(max)
	color := Green
	if percent < 0.5 { color = Yellow }
	if percent < 0.2 { color = Red }

	fill := (hp * 20) / max
	if fill < 0 { fill = 0 }
	bar := color + strings.Repeat("█", fill) + Gray + strings.Repeat("░", 20-fill) + Reset
	return fmt.Sprintf("%s %d/%d HP", bar, hp, max)
}

func attack(atk, def *Character) {
	dmg := atk.Str - def.Def
	if dmg < 1 { dmg = 1 }


	def.Hp -= dmg
	fmt.Printf("%s%s%s attaque %s !\n 💥 %s-%d HP%s\n", Cyan, atk.Name, Reset, def.Name, Red, dmg, Reset)
	time.Sleep(200 * time.Millisecond)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	loadFiles()

	p := Players[0]
	fmt.Printf("\n%s╔═══════════════════════════════════════╗\n║        🏰 CHATEAU D'HYRULE 🏰            ║\n╚═══════════════════════════════════════╝%s\n", Yellow, Reset)
	fmt.Printf(" Héros : %s%s%s\n", Green, p.Name, Reset)

	for floor := 1; floor <= 10; floor++ {
		e := Enemies[rand.Intn(len(Enemies))]
		if floor == 10 { e = Bosses[0]; fmt.Printf("\n%s🔥 BOSS FINAL : %s 🔥%s\n", Red, e.Name, Reset) }
		e.MaxHp = e.Hp

		fmt.Printf("\n%s--- ÉTAGE %d ---%s\n", Blue, floor, Reset)
		fmt.Printf("⚠️  %s%s%s sauvage apparait !\n", Red, e.Name, Reset)

		for p.Hp > 0 && e.Hp > 0 {
			fmt.Printf("\n%-15s %s\n", p.Name, drawBar(p.Hp, p.MaxHp))
			fmt.Printf("%-15s %s\n", e.Name, drawBar(e.Hp, e.MaxHp))
			
			fmt.Printf("\n%s[1]%s Attaquer  %s[2]%s Soin : ", Cyan, Reset, Cyan, Reset)
			var choice int
			fmt.Scan(&choice)

			fmt.Println(Gray + "---------------------------------------" + Reset)
			
			if p.Spd >= e.Spd {
				if choice == 1 { attack(&p, &e) } else { p.Hp += p.MaxHp/3; fmt.Println("💚 Soin !") }
				if e.Hp > 0 { attack(&e, &p) }
			} else {
				attack(&e, &p)
				if p.Hp > 0 {
					if choice == 1 { attack(&p, &e) } else { p.Hp += p.MaxHp/3; fmt.Println("💚 Soin !") }
				}
			}
			if p.Hp > p.MaxHp { p.Hp = p.MaxHp }
		}

		if p.Hp <= 0 { 
			fmt.Printf("\n%s💀 GAME OVER... %s s'est effondré.%s\n", Red, p.Name, Reset)
			return 
		}
		fmt.Printf("\n%s✅ Énnemi vaincu !%s\n", Green, Reset)
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\n%s🌟 VICTOIRE ! TU EST NOTRE HERO ! 🌟%s\n", Yellow, Reset)
}