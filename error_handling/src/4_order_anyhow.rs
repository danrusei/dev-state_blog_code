use anyhow::{anyhow, Context};
use chrono::NaiveDate;
use std::{env, fmt, fs};

#[derive(Debug)]
struct UserCommand {
    product: String,
    quantity: u32,
    delivery_date: NaiveDate,
}

fn main() -> anyhow::Result<()> {
    let file_name = env::args()
        .nth(1)
        .ok_or_else(|| anyhow!("Please provide a file name as a command-line argument."))?;
    let content = fs::read_to_string(&file_name)
        .with_context(|| format!("Error reading the file: {}", &file_name))?;

    let mut commands: Vec<UserCommand> = Vec::new();

    for (line_num, line) in content.lines().enumerate() {
        let mut parts = line.split_whitespace();
        let product = parts
            .next()
            .ok_or_else(|| anyhow!("Missing product information at line {}", line_num + 1))?
            .to_string();

        let quant = parts
            .next()
            .ok_or_else(|| anyhow!("Missing quantity information at line {}", line_num + 1))?;
        let quantity = quant
            .trim()
            .parse::<u32>()
            .with_context(|| format!("Invalid quantity format at line {}", line_num + 1))?;

        let d_date = parts
            .next()
            .ok_or_else(|| anyhow!("Missing delivery date information at line {}", line_num + 1))?;
        let delivery_date = NaiveDate::parse_from_str(d_date.trim(), "%d.%m.%Y")
            .with_context(|| format!("Invalid date format at line {}", line_num + 1))?;

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
