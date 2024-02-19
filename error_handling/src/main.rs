use chrono::NaiveDate;
use std::{env, fs};

#[derive(Debug)]
struct UserCommand {
    product: String,
    quantity: u32,
    delivery_date: NaiveDate,
}

fn main() {
    let file_name = env::args().nth(1).unwrap();
    let content = fs::read_to_string(file_name).unwrap();

    let mut commands: Vec<UserCommand> = Vec::new();

    for line in content.lines() {
        let mut parts = line.split_whitespace();
        let product = parts.next().unwrap().to_string();

        let quant = parts.next().unwrap();
        let quantity = quant.trim().parse::<u32>().unwrap();

        let d_date = parts.next().unwrap();
        let delivery_date = NaiveDate::parse_from_str(d_date.trim(), "%d.%m.%Y").unwrap();

        let command = UserCommand {
            product,
            quantity,
            delivery_date,
        };

        commands.push(command);
    }

    println!(
        "Your command was processed and it is ready for delivery. The ordered items: {:#?}",
        commands
    );
}
