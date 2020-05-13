mod parse;
use parse::{parse, Error};
use rand::Rng;

use async_std::task;
use futures::channel::mpsc::{channel, Receiver};
use futures::stream::StreamExt;
use futures::sink::SinkExt;

async fn generator(n_jobs: u32) -> Result<Receiver<u32>, Error> {
    let (mut tx, rx) = channel(0);
    let mut rng = rand::thread_rng();
    let nums: Vec<u32> = (0..n_jobs).map(|_| rng.gen_range(1, 100)).collect();
    task::spawn(async move {
        for num in nums {
            println!("I'm sending from generator");
            tx.send(num).await
                .expect("Could not send the generated number over the channel")
        }
    });
    Ok(rx)
}

async fn fan_out(mut rx_gen: Receiver<u32>) -> Result<Receiver<String>, Error> {
    let (tx, rx) = channel(0);

    let mut handles = Vec::new();
    loop {
        match rx_gen.next().await {
            Some(num) => {
                let mut tx_num = tx.clone();
                let task = task::spawn(async move {
                    println!("I'm within a task from fanout");
                    let rep = parse(num).await.unwrap();
                    tx_num
                        .send(rep).await
                        .expect("Could not send the parsed string over the channel");
                });
                handles.push(task);
            }
            None => break,
        }
    }

    for handle in handles {
        handle.await;
    }

    Ok(rx)
}

async fn fan_in(mut rx_fan_out: Receiver<String>) -> Result<Receiver<String>, Error> {
    let (mut tx, rx) = channel(0);
    task::spawn(async move {
        loop {
            match rx_fan_out.next().await {
                Some(value) => {
                    let processed_value = format!("{} _ Processed", value);

                    println!("I'm collecting in FanIn");

                    tx.send(processed_value).await
                        .expect("Could not send the processed string over the channel");
                }
                None => break,
            }
        }
    });

    Ok(rx)
}

#[async_std::main]
async fn main() -> Result<(), Error> {
    let n_jobs = 8;
    let mut rx_fan_in = fan_in(fan_out(generator(n_jobs).await?).await?).await?;
    loop {
        match rx_fan_in.next().await {
            Some(value) => {
                println!("{}", value);
            }
            None => break,
        }
    }
    Ok(())
}
