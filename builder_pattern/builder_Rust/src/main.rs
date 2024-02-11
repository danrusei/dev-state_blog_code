use anyhow::Result;
use builder_rust::car::CarBuilder;
use builder_rust::car2::Car2Builder;
use builder_rust::*;

fn main() -> Result<()> {
    // Create the Car using consuming Pattern
    let equip = Equiped::EquipedGold;

    let car = CarBuilder::new(equip)
        .set_color("Blue")
        .set_engine(Engine::Electric(Electric::Electric350kW))
        .build()?;

    println!("Consuming Builder Pattern:");
    println!("You selected a {}", car);

    // Create the Car using non-consuming Pattern

    let equip = Equiped::EquipedSilver;

    let mut car_builder = Car2Builder::new(equip);

    let car1 = car_builder
        .set_color("Black")
        .set_engine(Engine::Petrol(Petrol::Petrol150HP))
        .build()?;

    let car2 = car_builder
        .set_color("Red")
        .set_engine(Engine::Diesel(Diesel::Diesel250HP))
        .build()?;

    println!("Non-Consuming Builder Pattern:");
    println!("You selected a {}", car1);
    println!("You selected a {}", car2);

    Ok(())
}
