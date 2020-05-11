mod parse;
use parse::parse;

fn main() {
    println!("{}",parse(10).unwrap())
}
