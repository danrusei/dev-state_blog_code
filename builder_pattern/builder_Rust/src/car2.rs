// Non-Consuming Pattern

use crate::*;
use anyhow::{anyhow, Result};
use std::default::Default;
use std::fmt;

#[allow(dead_code)]
#[derive(Debug, Default)]
pub struct Car2 {
    equip: Equiped,
    color: Option<String>,
    car_engine: Option<Engine>,
    first_private_field: Option<String>,
    second_private_field: Option<String>,
}

pub struct Car2Builder {
    car: Car2,
}

impl Car2Builder {
    pub fn new(equip: Equiped) -> Self {
        Car2Builder {
            car: Car2 {
                equip,
                ..Default::default()
            },
        }
    }

    pub fn set_color(&mut self, color: impl Into<String>) -> &mut Self {
        self.car.color = Some(color.into());
        self
    }

    pub fn set_engine(&mut self, engine: Engine) -> &mut Self {
        self.car.car_engine = Some(engine);
        self
    }

    pub fn build(&mut self) -> Result<Car2> {
        let car_engine = self
            .car
            .car_engine
            .take()
            .ok_or_else(|| anyhow!("You need to select an Engine type"))?;

        Ok(Car2 {
            equip: std::mem::take(&mut self.car.equip),
            color: Some(std::mem::take(&mut self.car.color).unwrap_or("White".to_owned())),
            car_engine: Some(car_engine),
            ..Default::default()
        })
    }
}

impl fmt::Display for Car2 {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "Car ")?;
        write!(f, "equiped with {:?}, ", self.equip)?;

        if let Some(ref engine) = self.car_engine {
            write!(f, "having a {:?} engine ,", engine)?;
        } else {
            write!(f, "car_engine: None")?;
        }

        if let Some(ref color) = self.color {
            write!(f, "and {} color. ", color)?;
        } else {
            write!(f, "color: None")?;
        }

        write!(f, "")
    }
}
