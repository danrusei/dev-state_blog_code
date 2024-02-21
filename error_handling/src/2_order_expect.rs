use chrono::NaiveDate;
use std::{env, fmt, fs};

#[derive(Debug)]
struct UserCommand {
    product: String,
    quantity: u32,
    delivery_date: NaiveDate,
}

fn main() {
    let file_name = env::args()
        .nth(1)
        .expect("Please provide a file name as a command-line argument.");
    let content = fs::read_to_string(&file_name).expect("Error reading the file");

    let mut commands: Vec<UserCommand> = Vec::new();

    for line in content.lines() {
        let mut parts = line.split_whitespace();
        let product = parts
            .next()
            .expect("Missing product information")
            .to_string();

        let quant = parts.next().expect("Missing quantity information");
        let quantity = quant
            .trim()
            .parse::<u32>()
            .expect("Invalid quantity format, expecting integer");

        let d_date = parts.next().expect("Missing delivery date information");
        let delivery_date = NaiveDate::parse_from_str(d_date.trim(), "%d.%m.%Y")
            .expect("Invalid date format, should be %d.%m.%Y");

        let command = UserCommand {
            product,
            quantity,
            delivery_date,
        };

        commands.push(command);
    }

    println!("\nYour command was processed and it is ready for delivery. The ordered items:\n");

    commands.iter().for_each(|cmd| println!(" * {} ", cmd));
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
