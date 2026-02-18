package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

// structure pour lire les classes et les races depuis les fichiers json
type Stats struct {
	ID int; Name string; Strengths []int; Weaknesses []int; AttackType string
}

// structure pour les persos (joueur, ennemis, boss)
// j'utilise les tags json pour que ça corresponde aux fichiers
type Character struct {
	Name  string `json:"name"`
	Hp    int    `json:"hp"`
	MaxHp int    // pas dans le json, je le calcule au chargement pour la barre de vie
	Mp    int    `json:"mp"`
	Str   int    `json:"str"`; Int int `json:"int"`; Def int `json:"def"`; Res int `json:"res"`
	Spd   int    `json:"spd"`; Luck int `json:"luck"`
	Race  int    `json:"race"`; Class int `json:"class"`
}

// j'utilise des maps pour retrouver facile les infos avec l'ID
var Classes = make(map[int]Stats)
var Races   = make(map[int]Stats)

// listes pour stocker tout le monde
var Players []Character
var Enemies []Character
var Bosses  []Character

// fonction pratique pour lire un fichier json et le mettre dans une variable
func load(file string, v interface{}) {
	// on lit le fichier dans le dossier mods
	data, err := ioutil.ReadFile("mods/" + file)
	if err != nil { panic(fmt.Sprintf("Erreur lecture %s: %v", file, err)) }
	
	// on transforme le json en objet go
	if err := json.Unmarshal(data, v); err != nil { panic(fmt.Sprintf("Erreur JSON %s: %v", file, err)) }
}

// ici on charge tout au lancement du jeu
func initGame() {
	var clsList, rcList []Stats
	
	// chargement des fichiers de config
	load("classes.json", &clsList); load("races.json", &rcList)
	
	// je remplis les maps pour pouvoir chercher par ID plus tard
	for _, c := range clsList { Classes[c.ID] = c }
	for _, r := range rcList  { Races[r.ID] = r }

	// chargement des persos
	load("players.json", &Players); load("enemies.json", &Enemies); load("bosses.json", &Bosses)

	// important: on fixe les HP max au début car c'est égal aux HP actuels
	for i := range Players { Players[i].MaxHp = Players[i].Hp }
	for i := range Enemies { Enemies[i].MaxHp = Enemies[i].Hp }
	for i := range Bosses  { Bosses[i].MaxHp = Bosses[i].Hp }
}

// petite fonction pour voir si un ID est dans une liste (pour les forces/faiblesses)
func contains(list []int, id int) bool {
	for _, v := range list { if v == id { return true } }
	return false
}

// gestion du combat
func attack(atk, def *Character) {
	// on recupere les infos de classe et race
	clsAtk, rcAtk := Classes[atk.Class], Races[atk.Race]
	clsDef, rcDef := Classes[def.Class], Races[def.Race]

	// calcul des dégats de base selon le type d'attaque
	dmg, defense := atk.Str, def.Def
	if clsAtk.AttackType == "magical" { dmg, defense = atk.Int, def.Res }
	
	base := dmg - defense
	if base < 0 { base = 0 } // pour pas soigner l'ennemi si on tape pas assez fort

	// 2. gestion des multiplicateurs (x2 ou /2)
	mult := 1.0
	// on vérifie les forces
	if contains(clsAtk.Strengths, clsDef.ID) || contains(rcAtk.Strengths, rcDef.ID) { mult *= 2 }
	// on vérifie les faiblesses
	if contains(clsAtk.Weaknesses, clsDef.ID) || contains(rcAtk.Weaknesses, rcDef.ID) { mult *= 0.5 }

	// on applique les dégats
	finalDmg := int(float64(base) * mult)
	if finalDmg < 1 { finalDmg = 1 } // minimum 1 dégat
	def.Hp -= finalDmg

	// petit message pour le joueur
	msg := ""
	if mult > 1 { msg = " [Coup Critique!]" } else if mult < 1 { msg = " [Peu Efficace...]" }
	fmt.Printf("%s attaque %s !%s Dégâts: %d\n", atk.Name, def.Name, msg, finalDmg)
}

// fonction pour dessiner la barre de vie [====    ]
func bar(cur, max int) string {
	if cur < 0 { cur = 0 }
	if cur > max { cur = max }
	// produit en croix pour une barre de 20 char
	fill := (cur * 20) / max
	return "[" + strings.Repeat("=", fill) + strings.Repeat(" ", 20-fill) + "]"
}

func main() {
	// pour que l'aléatoire change à chaque lancement
	rand.Seed(time.Now().UnixNano())
	
	initGame() // on charge tout

	p := Players[0] // on prend Link par défaut
	fmt.Printf("=== HYRULE CASTLE ===\nHéros: %s\n", p.Name)

	// boucle des 10 étages
	for floor := 1; floor <= 10; floor++ {
		e := Enemies[rand.Intn(len(Enemies))] // ennemi au pif
		
		if floor == 10 { 
			fmt.Printf("\n--- BOSS FINAL ---\n")
			e = Bosses[rand.Intn(len(Bosses))] // boss au pif
		} else {
			fmt.Printf("\n--- ÉTAGE %d ---\n", floor)
		}
		
		fmt.Printf("%s apparaît !\n", e.Name)

		// boucle de combat : tant que personne n'est mort
		for p.Hp > 0 && e.Hp > 0 {
			// affichage des stats
			fmt.Printf("\n%s %s %d/%d\n", p.Name, bar(p.Hp, p.MaxHp), p.Hp, p.MaxHp)
			fmt.Printf("%s %s %d/%d\n", e.Name, bar(e.Hp, e.MaxHp), e.Hp, e.MaxHp)
			
			var ch int
			fmt.Print("1.Attaquer 2.Soin : ")
			fmt.Scan(&ch)

			// on vérifie la vitesse pour savoir qui tape en premier
			if p.Spd >= e.Spd {
				// tour du joueur
				if ch == 1 { 
					attack(&p, &e) 
				} else { 
					// soin simple (remet la moitié de la vie max)
					p.Hp += p.MaxHp/2
					if p.Hp > p.MaxHp { p.Hp = p.MaxHp }
					fmt.Println("Soin !") 
				}
				
				// riposte de l'ennemi s'il est vivant
				if e.Hp > 0 { attack(&e, &p) }
			} else {
				// l'ennemi tape d'abord
				attack(&e, &p)
				
				// si le joueur survit, il joue
				if p.Hp > 0 {
					if ch == 1 { 
						attack(&p, &e) 
					} else { 
						p.Hp += p.MaxHp/2
						if p.Hp > p.MaxHp { p.Hp = p.MaxHp }
						fmt.Println("Soin !") 
					}
				}
			}
		}

		// fin du combat
		if p.Hp <= 0 { fmt.Println("GAME OVER"); os.Exit(0) }
		fmt.Println("Ennemi vaincu ! Entrée pour continuer...")
		fmt.Scanln(); fmt.Scanln()
	}
	fmt.Println("VICTOIRE !")
}