package main

import (
	"fmt"

	"github.com/danrusei/dev-state_blog_code/tree/master/builder_pattern/builder_Go/car1"
	"github.com/danrusei/dev-state_blog_code/tree/master/builder_pattern/builder_Go/car2"
)

func main() {
	// Example usage of Builder Pattern:
	builder := car1.NewCarBuilder(car1.EquipedGold).
		SetColor("Brown").
		SetPetrolEngine(car1.Petrol225HP)

	car_ex1 := builder.Build()
	fmt.Println("===========================================")
	fmt.Printf("Builder Pattern: %v\n", car_ex1)
	fmt.Println("===========================================")

	// ---------------------------------------------

	// Example usage of Functional Options Pattern
	car_ex2 := car2.NewCar(
		car2.WithEquipment(car2.EquipedSilver),
		car2.WithColor("White"),
		car2.WithElectricEngine(car2.Electric420kW),
	)

	fmt.Printf("Functional Options Pattern: %s\n", car_ex2)

}
