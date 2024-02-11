pub mod car;
pub mod car2;

#[derive(Debug, Default, Clone)]
pub enum Equiped {
    #[default]
    EquipedSilver,
    EquipedGold,
    EquipedPlatinum,
}

#[derive(Debug, Clone)]
pub enum Engine {
    Petrol(Petrol),
    Diesel(Diesel),
    Electric(Electric),
}

#[derive(Debug, Clone)]
pub enum Petrol {
    Petrol150HP,
    Petrol190HP,
    Petrol225HP,
}

#[derive(Debug, Clone)]
pub enum Diesel {
    Diesel120HP,
    Diesel200HP,
    Diesel250HP,
}

#[derive(Debug, Clone)]
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
