package car1

type Petrol int
type Diesel int
type Electric int
type Equiped int

const (
	Petrol150HP Petrol = iota
	Petrol190HP
	Petrol225HP

	Diesel120HP Diesel = iota
	Diesel200HP
	Diesel250HP

	Electric350kW Electric = iota
	Electric420kW
	Electric588kW

	EquipedSilver Equiped = iota
	EquipedGold
	EquipedPlatinum
)

// A car with specific attributes. This struct will not be exposed to users!!! as it may contain fields that has to remain private.
type Car struct {
	equip              Equiped
	color              string
	carEngine          Engine
	firstPrivateField  string
	secondPrivatefield string
}

// Engine represents the type of car engine
type Engine struct {
	petrol   Petrol
	diesel   Diesel
	electric Electric
}

// CarBuilder is a builder for creating Car instances
type CarBuilder struct {
	car Car
}

// NewCarBuilder creates a new CarBuilder instance
func NewCarBuilder(equip Equiped) *CarBuilder {
	return &CarBuilder{
		car: Car{
			equip: equip,
		},
	}
}

// SetColor sets the color of the car
func (cb *CarBuilder) SetColor(color string) *CarBuilder {
	cb.car.color = color
	return cb
}

// SetPetrolEngine sets the petrol engine type
func (cb *CarBuilder) SetPetrolEngine(power Petrol) *CarBuilder {
	cb.car.carEngine.petrol = power
	return cb
}

// SetDieselEngine sets the diesel engine type
func (cb *CarBuilder) SetDieselEngine(power Diesel) *CarBuilder {
	cb.car.carEngine.diesel = power
	return cb
}

// SetElectricEngine sets the electric engine type
func (cb *CarBuilder) SetElectricEngine(power Electric) *CarBuilder {
	cb.car.carEngine.electric = power
	return cb
}

// Build constructs and returns the final Car instance
func (cb *CarBuilder) Build() Car {
	return cb.car
}
