package main

import "fmt"

// Интерфейсы продуктов
type IFavicon interface {
	setLogo(logo string)
	getLogo() string
}

type IBrand interface {
	setLogo(logo string)
	getLogo() string
}

// Базовые реализации (для уменьшения дублирования кода)
type Favicon struct{ logo string }

func (f *Favicon) setLogo(logo string) { f.logo = logo }
func (f *Favicon) getLogo() string     { return f.logo }

type Brand struct{ logo string }

func (b *Brand) setLogo(logo string) { b.logo = logo }
func (b *Brand) getLogo() string     { return b.logo }

// Конкретные продукты
type AdidasBrand struct{ Brand }
type AdidasFavicon struct{ Favicon }

type NikeBrand struct{ Brand }
type NikeFavicon struct{ Favicon }

type ILogoFactory interface {
	makeFavicon() IFavicon
	makeBrand() IBrand
}

// --- Фабрика Adidas ---
type Adidas struct{}

func (a *Adidas) makeFavicon() IFavicon {
	return &AdidasFavicon{Favicon: Favicon{logo: "Adidas"}}
}

func (a *Adidas) makeBrand() IBrand {
	return &AdidasBrand{Brand: Brand{logo: "Adidas"}}
}

// --- Фабрика Nike ---
type Nike struct{}

func (n *Nike) makeFavicon() IFavicon {
	return &NikeFavicon{Favicon: Favicon{logo: "Nike"}}
}

func (n *Nike) makeBrand() IBrand {
	return &NikeBrand{Brand: Brand{logo: "Nike"}}
}

// --- Функция получения фабрики ---
func getLogoFactory(brand string) (ILogoFactory, error) {
	if brand == "adidas" {
		return &Adidas{}, nil
	}
	if brand == "nike" {
		return &Nike{}, nil
	}
	return nil, fmt.Errorf("Wrong brand passed")
}

func main() {
	factory, _ := getLogoFactory("nike")
	favicon := factory.makeFavicon()
	brand := factory.makeBrand()

	fmt.Printf("Logo (favicon): %s\n", favicon.getLogo())
	fmt.Printf("Logo (brand): %s\n", brand.getLogo())
}
