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

	fmt.Println("Bienvenue dans le chateau d'Hyrule !")

	// Boucle. on va de l'etage 1 a 10
	for floor := 1; floor <= 10; floor++ {
		
		// def de l'ennemi pour cet étage
		var enemy Character

		if floor == 10 {
			// Si derniere etage alors = boss
			fmt.Printf("\n--- ÉTAGE %d : BOSS FINAL ---\n", floor)
			enemy = Character{Name: "Ganon", Hp: 150, MaxHp: 150, Str: 20}
		} else {
			// Sinon c'est un monstre normal 
			fmt.Printf("\n--- ÉTAGE %d ---\n", floor)
			enemy = Character{Name: "Bokoblin", Hp: 30, MaxHp: 30, Str: 5}
		}

		fmt.Printf("Un %s sauvage apparaît avec %d HP !\n", enemy.Name, enemy.Hp)

		// La boucle de combat
		// Tant que Link et l'ennemi sont vivants, on continue le combat
		for link.Hp > 0 && enemy.Hp > 0 {
			
			// Affichage des stats
			fmt.Printf("\n%s (HP: %d/%d) vs %s (HP: %d/%d)\n", link.Name, link.Hp, link.MaxHp, enemy.Name, enemy.Hp, enemy.MaxHp)
			fmt.Println("Actions : 1. Attaquer | 2. Se soigner")
			fmt.Print("Votre choix : ")

			// Lecture du choix du joueur
			var choice int
			fmt.Scan(&choice)

			// Tour du joueur
			switch choice {
			case 1: // Attaque
				fmt.Printf("%s attaque et inflige %d dégâts !\n", link.Name, link.Str)
				enemy.Hp -= link.Str // On retire des PV a l'ennemi
			case 2: // Soin
				healAmount := link.MaxHp / 2
				// On verifie de ne pas depasser le MaxHp
				if link.Hp+healAmount > link.MaxHp {
					link.Hp = link.MaxHp
				} else {
					link.Hp += healAmount
				}
				fmt.Printf("%s se soigne et récupère %d HP.\n", link.Name, healAmount)
			default:
				fmt.Println("Action inconnue, vous passez votre tour.. ")
			}

			// Verif l'ennemi est il mort ?
			if enemy.Hp <= 0 {
				fmt.Printf("Le %s est vaincu !\n", enemy.Name)
				break // On sort de la boucle de combat on passe à l'etage suivant
			}

			// Tour de l'enemie
			fmt.Printf("Le %s riposte et inflige %d dégâts !\n", enemy.Name, enemy.Str)
			link.Hp -= enemy.Str

			// verif si link est mort
			if link.Hp <= 0 {
				fmt.Println("\n=== GAME OVER ===")
				fmt.Println("Vous etes mort...")
				os.Exit(0) // On arrête le programme
			}
		}
	}

	// Si on sort vivant du chateau
	fmt.Println("\n=== FÉLICITATIONS ===")
	fmt.Println("Vous avez vaincu Ganon!")
}