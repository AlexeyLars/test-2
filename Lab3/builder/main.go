package main

import "fmt"

type House struct {
	windowType string
	doorType   string
	floor      int
}

type IBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() House
}

func getBuilder(builderType string) IBuilder {
	if builderType == "normal" {
		return &NormalBuilder{}
	}
	if builderType == "igloo" {
		return &IglooBuilder{}
	}
	return nil
}

type NormalBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func (b *NormalBuilder) setWindowType() { b.windowType = "Деревянное окно" }
func (b *NormalBuilder) setDoorType()   { b.doorType = "Деревянная дверь" }
func (b *NormalBuilder) setNumFloor()   { b.floor = 2 }
func (b *NormalBuilder) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType:   b.doorType,
		floor:      b.floor,
	}
}

type IglooBuilder struct {
	windowType string
	doorType   string
	floor      int
}

func (b *IglooBuilder) setWindowType() { b.windowType = "Ледяное окно" }
func (b *IglooBuilder) setDoorType()   { b.doorType = "Ледяная дверь" }
func (b *IglooBuilder) setNumFloor()   { b.floor = 1 }
func (b *IglooBuilder) getHouse() House {
	return House{
		windowType: b.windowType,
		doorType:   b.doorType,
		floor:      b.floor,
	}
}

type Director struct {
	builder IBuilder
}

func newDirector(b IBuilder) *Director {
	return &Director{builder: b}
}

func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) buildHouse() House {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}

func main() {
	normalBuilder := getBuilder("normal")
	iglooBuilder := getBuilder("igloo")

	director := newDirector(normalBuilder)
	normalHouse := director.buildHouse()

	fmt.Printf("Normal House: Door: %s, Window: %s, Floors: %d\n",
		normalHouse.doorType, normalHouse.windowType, normalHouse.floor)

	director.setBuilder(iglooBuilder)
	iglooHouse := director.buildHouse()

	fmt.Printf("Igloo House: Door: %s, Window: %s, Floors: %d\n",
		iglooHouse.doorType, iglooHouse.windowType, iglooHouse.floor)
}
