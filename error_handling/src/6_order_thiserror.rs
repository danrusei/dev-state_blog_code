use chrono::NaiveDate;
use std::{env, fmt, fs};
use thiserror::Error;

#[derive(Debug, Error)]
enum MyError {
    #[error("Please provide a file name as a command-line argument.")]
    CommandLineArgs,

    #[error("Error reading the file: {0}")]
    FileReadError(#[from] std::io::Error),

    #[error("Parsing error: {0}")]
    ParsingError(String),
}

#[derive(Debug)]
struct UserCommand {
    product: String,
    quantity: u32,
    delivery_date: NaiveDate,
}

fn main() -> Result<(), MyError> {
    let file_name = env::args().nth(1).ok_or(MyError::CommandLineArgs)?;
    let content = fs::read_to_string(&file_name).map_err(MyError::FileReadError)?;

    let mut commands: Vec<UserCommand> = Vec::new();

    for (line_num, line) in content.lines().enumerate() {
        let mut parts = line.split_whitespace();

        let product = parts
            .next()
            .ok_or(MyError::ParsingError(format!(
                "Missing product information in line {}",
                line_num + 1
            )))?
            .to_string();

        let quant = parts.next().ok_or(MyError::ParsingError(format!(
            "Missing quantity information in line {}",
            line_num + 1
        )))?;
        let quantity = quant.trim().parse::<u32>().map_err(|e| {
            MyError::ParsingError(format!(
                "Invalid quantity format in line {}: {}",
                line_num + 1,
                e
            ))
        })?;

        let d_date = parts.next().ok_or(MyError::ParsingError(format!(
            "Missing delivery date information in line {}",
            line_num + 1
        )))?;
        let delivery_date = NaiveDate::parse_from_str(d_date.trim(), "%d.%m.%Y").map_err(|e| {
            MyError::ParsingError(format!(
                "Invalid date format in line {}: {}",
                line_num + 1,
                e
            ))
        })?;

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
