//Consuming Pattern

use crate::*;
use anyhow::{anyhow, Result};
use std::default::Default;
use std::fmt;

#[allow(dead_code)]
#[derive(Debug, Default)]
pub struct Car {
    equip: Equiped,
    color: Option<String>,
    car_engine: Option<Engine>,
    first_private_field: Option<String>,
    second_private_field: Option<String>,
}

pub struct CarBuilder {
    car: Car,
}

impl CarBuilder {
    pub fn new(equip: Equiped) -> Self {
        CarBuilder {
            car: Car {
                equip: equip,
                ..Default::default()
            },
        }
    }
    pub fn set_color(mut self, color: impl Into<String>) -> Self {
        self.car.color = Some(color.into());
        self
    }
    pub fn set_engine(mut self, engine: Engine) -> Self {
        self.car.car_engine = Some(engine);
        self
    }
    pub fn build(self) -> Result<Car> {
        let Some(car_engine) = self.car.car_engine else {
            return Err(anyhow!("You need to select an Engine type"));
        };

        Ok(Car {
            equip: self.car.equip,
            color: Some(self.car.color).unwrap_or(Some("White".to_owned())),
            car_engine: Some(car_engine),
            ..Default::default()
        })
    }
}

impl fmt::Display for Car {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Car ")?;
        write!(f, "equiped with {:?}, ", self.equip)?;

        if let Some(ref engine) = self.car_engine {
            write!(f, "having a {:?} engine ,", engine)?;
        } else {
            write!(f, "car_engine: None")?;
        }

        if let Some(ref color) = self.color {
            write!(f, "and {} color.", color)?;
        } else {
            write!(f, "color: None")?;
        }

        write!(f, "")
    }
}
