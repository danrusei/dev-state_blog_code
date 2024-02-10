package main

import (
	"fmt"

	"github.com/danrusei/dev-state_blog_code/tree/master/builder_pattern/builder_Go/car1"
	"github.com/danrusei/dev-state_blog_code/tree/master/builder_pattern/builder_Go/car2"
)

func main() {
	// Example usage of Builder Pattern:
	builder := car1.NewCarBuilder(car1.EquipedGold).
		SetColor("Silver").
		SetPetrolEngine(car1.Petrol225HP)

	car_ex1 := builder.Build()
	fmt.Printf("Builder Pattern: %v\n", car_ex1)

	// ---------------------------------------------

	// Example usage of Functional Options Pattern
	car_ex2 := car2.NewCar(
		car2.WithEquipment(car2.EquipedGold),
		car2.WithColor("Silver"),
		car2.WithPetrolEngine(car2.Petrol225HP),
	)

	fmt.Printf("Functional Options Pattern: %s\n", car_ex2)

}
