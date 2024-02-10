package car2

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

// Engine represents the type of car engine
type Engine struct {
	Petrol   Petrol
	Diesel   Diesel
	Electric Electric
}

// Car represents a car with specific attributes
type Car struct {
	equip              Equiped
	color              string
	carEngine          Engine
	firstPrivateField  string
	secondPrivatefield string
}

// Option is a functional option for configuring Car instances
type Option func(*Car)

// WithColor sets the color of the car
func WithColor(color string) Option {
	return func(c *Car) {
		c.color = color
	}
}

// WithPetrolEngine sets the petrol engine type
func WithPetrolEngine(power Petrol) Option {
	return func(c *Car) {
		c.carEngine.Petrol = power
	}
}

// WithDieselEngine sets the diesel engine type
func WithDieselEngine(power Diesel) Option {
	return func(c *Car) {
		c.carEngine.Diesel = power
	}
}

// WithElectricEngine sets the electric engine type
func WithElectricEngine(power Electric) Option {
	return func(c *Car) {
		c.carEngine.Electric = power
	}
}

// WithEquipment sets the equipment type
func WithEquipment(equip Equiped) Option {
	return func(c *Car) {
		c.equip = equip
	}
}

// NewCar creates a new Car instance with specified options
func NewCar(options ...Option) *Car {
	car := &Car{}
	for _, option := range options {
		option(car)
	}
	return car
}
