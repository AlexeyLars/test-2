package main

import "fmt"

type IGun interface {
	setName(name string)
	setPower(power int)
	getName() string
	getPower() int
}

type Gun struct {
	name  string
	power int
}

func (g *Gun) setName(name string) { g.name = name }
func (g *Gun) getName() string     { return g.name }
func (g *Gun) setPower(power int)  { g.power = power }
func (g *Gun) getPower() int       { return g.power }

type MagicWand struct {
	Gun // Встраивание (композиция)
}

func newMagicWand() IGun {
	return &MagicWand{
		Gun: Gun{
			name:  "Magic Wand",
			power: 4,
		},
	}
}

type Sword struct {
	Gun
}

func newSword() IGun {
	return &Sword{
		Gun: Gun{
			name:  "Creeper Sword",
			power: 1,
		},
	}
}

func getGun(gunType string) (IGun, error) {
	if gunType == "magic_wand" {
		return newMagicWand(), nil
	}
	if gunType == "sword" {
		return newSword(), nil
	}
	return nil, fmt.Errorf("Wrong gun type passed")
}

func main() {
	sword, _ := getGun("sword")
	wand, _ := getGun("magic_wand")

	fmt.Printf("Gun: %s Power: %d\n", sword.getName(), sword.getPower())
	fmt.Printf("Gun: %s Power: %d\n", wand.getName(), wand.getPower())
}
