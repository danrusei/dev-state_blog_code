pub mod car;
pub mod car2;

#[derive(Debug, Default)]
pub enum Equiped {
    #[default]
    EquipedSilver,
    EquipedGold,
    EquipedPlatinum,
}

#[derive(Debug)]
pub enum Engine {
    Petrol(Petrol),
    Diesel(Diesel),
    Electric(Electric),
}

#[derive(Debug)]
pub enum Petrol {
    Petrol150HP,
    Petrol190HP,
    Petrol225HP,
}

#[derive(Debug)]
pub enum Diesel {
    Diesel120HP,
    Diesel200HP,
    Diesel250HP,
}

#[derive(Debug)]
pub enum Electric {
    Electric350kW,
    Electric420kW,
    Electric588kW,
}

impl Default for Engine {
    fn default() -> Self {
        Engine::Petrol(Petrol::Petrol150HP)
    }
}
