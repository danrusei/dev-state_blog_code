mod parse;
use parse::{parse, Error};

//use async_std::prelude::*;
//use async_std::stream;
use async_std::task;
use futures::channel::mpsc::{channel, Receiver};
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
                .expect("Could not send the generated number over the channel")
        }
    });
    Ok(rx)
}

async fn fan_out(mut rx_gen: Receiver<u32>) -> Result<Receiver<String>, Error> {
    let (tx, rx) = channel(0);

    let mut tasks = vec![];
    while let Some(num) = rx_gen.try_next()? {
        let mut tx_num = tx.clone();
        let task = task::spawn(async move {
            let rep = parse(num).await.unwrap();
            tx_num
                .try_send(rep)
                .expect("Could not send the parsed string over the channel");
        });
        tasks.push(task);
    }

    for task in tasks.into_iter() {
        task.await;
    }

    Ok(rx)
}

async fn fan_in(mut rx_fan_out: Receiver<String>) -> Result<Receiver<String>, Error> {
    let (mut tx, rx) = channel(0);
    task::spawn(async move {
        while let Ok(value) = rx_fan_out.try_next() {
            let processed_value: String = match value {
                Some(n) => format!("{} _ Processed", n),
                None => format!("Could not find the num"),
            };
            tx.try_send(processed_value)
                .expect("Could not send the processed string over the channel");
        }
    });

    Ok(rx)
}

#[async_std::main]
async fn main() -> Result<(), Error> {
    let n_jobs = 8;
    let mut rx_fan_in = fan_in(fan_out(generator(n_jobs).await?).await?).await?;
    loop {
        let eu: String = match rx_fan_in.try_next()? {
            Some(n) => format!("{} _ Processed", n),
            None => break,
        };
        println!("{}", eu);
    }
    /*
        let rx_gen = generator(n_jobs).await?;
        let rx_fan_out = fan_out(rx_gen).await?;
        let mut rx_fan_in = fan_in(rx_fan_out).await?;
        let eu: String = rx_fan_in.try_next()?.unwrap();
        println!("{}", eu);
    */
    /*
         while let Ok(item) = rx_fan_in.try_next() {
         //   println!("{}",item.try_next())
        }
    */
    Ok(())
}
