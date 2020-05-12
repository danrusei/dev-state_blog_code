mod parse;
use parse::{parse, Error};

//use async_std::prelude::*;
//use async_std::stream;
use async_std::task;
use futures_channel::mpsc::{channel, Receiver};
use rand::Rng;

/*
//This is another option to asynchronously generate random numbers
async fn generator_other(n_jobs: u32) -> Result<u32, Error> {
    let mut num = 0u32;
    let mut rng = rand::thread_rng();
    let mut s = stream::from_fn(|| {
        num = rng.gen_range(1, 100);
        Some(num)
    });

    Ok(s.next().await.unwrap())
}
*/

async fn generator(n_jobs: u32) -> Result<Receiver<u32>, Error> {
    let (mut tx, rx) = channel(0);
    let mut rng = rand::thread_rng();
    let nums: Vec<u32> = (0..n_jobs).map(|_| rng.gen_range(1, 100)).collect();
    task::spawn(async move {
        for num in nums {
            tx.try_send(num)
                .expect("Could not sent the generated number over the channel")
        }
    });
    Ok(rx)
}

#[async_std::main]
async fn main() -> Result<(), Error> {
    println!("{}", parse(10).await?);
    Ok(())
}
