use chrono::NaiveDate;
use std::{env, fmt, fs};

#[derive(Debug)]
struct UserCommand {
    product: String,
    quantity: u32,
    delivery_date: NaiveDate,
}

fn main() -> Result<(), String> {
    let file_name = env::args()
        .nth(1)
        .ok_or("Please provide a file name as a command-line argument.".to_string())?;
    let content =
        fs::read_to_string(&file_name).map_err(|e| format!("Error reading the file: {}", e))?;

    let mut commands: Vec<UserCommand> = Vec::new();

    for line in content.lines() {
        let mut parts = line.split_whitespace();
        let product = parts
            .next()
            .ok_or("Missing product information".to_string())?
            .to_string();

        let quant = parts
            .next()
            .ok_or("Missing quantity information".to_string())?;
        let quantity = quant
            .trim()
            .parse::<u32>()
            .map_err(|e| format!("Invalid quantity format: {}", e))?;

        let d_date = parts
            .next()
            .ok_or("Missing delivery date information".to_string())?;
        let delivery_date = NaiveDate::parse_from_str(d_date.trim(), "%d.%m.%Y")
            .map_err(|e| format!("Invalid date format: {}", e))?;

        let command = UserCommand {
            product,
            quantity,
            delivery_date,
        };

        commands.push(command);
    }

    println!("\nYour command was processed and it is ready for delivery. The ordered items:\n");

    for cmd in &commands {
        println!(" * {} ", cmd);
    }

    Ok(())
}

impl fmt::Display for UserCommand {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(
            f,
            "{} {} - to be delivered on {}",
            self.quantity, self.product, self.delivery_date
        )
    }
}
