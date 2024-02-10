package car1

import "fmt"

// String returns a string representation of the selected Engine
func (c Car) String() string {
	return fmt.Sprintf("You selected the Car equiped with %s, having a %s engine and %s color", c.equip, c.carEngine, c.color)
}

// String returns a string representation of the selected Engine
func (e Engine) String() string {
	switch {
	case e.petrol != 0:
		return fmt.Sprintf("%s", e.petrol)
	case e.diesel != 0:
		return fmt.Sprintf("%s", e.diesel)
	case e.electric != 0:
		return fmt.Sprintf("%s", e.electric)
	default:
		return "Unknown Engine Type"
	}
}

func (p Petrol) String() string {
	switch p {
	case Petrol150HP:
		return "Petrol 150HP"
	case Petrol190HP:
		return "Petrol 190HP"
	case Petrol225HP:
		return "Petrol 225HP"
	default:
		return "Unknown Petrol"
	}
}

func (d Diesel) String() string {
	switch d {
	case Diesel120HP:
		return "Diesel 120HP"
	case Diesel200HP:
		return "Diesel 200HP"
	case Diesel250HP:
		return "Diesel 250HP"
	default:
		return "Unknown Diesel"
	}
}

func (e Electric) String() string {
	switch e {
	case Electric350kW:
		return "Electric 350kW"
	case Electric420kW:
		return "Electric 420kW"
	case Electric588kW:
		return "Electric 588kW"
	default:
		return "Unknown Electric"
	}
}

func (e Equiped) String() string {
	switch e {
	case EquipedSilver:
		return "Silver Package"
	case EquipedGold:
		return "Gold Package"
	case EquipedPlatinum:
		return "Platinum Package"
	default:
		return "Unknown Package"
	}
}
