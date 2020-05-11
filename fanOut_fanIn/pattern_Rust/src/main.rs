mod parse;
use parse::parse;

use std::io;
use std::error::Error;
use std::thread;
use std::sync::mpsc::{channel, Receiver};
use threadpool::ThreadPool;
use rand::Rng;

fn generator(n_jobs: u32) -> io::Result<Receiver<u32>> {
    let (gen_sender, gen_receiver) = channel();
    let mut rng = rand::thread_rng();
    let nums: Vec<u32> = (0..n_jobs).map(|_| rng.gen_range(1,100)).collect();
    thread::spawn(move || {
        for num in nums{
            gen_sender.send(num).expect("Could not send the generated number over gen_sender channel")
        }
    });
    Ok(gen_receiver)
}

fn fan_out(rx_gen: Receiver<u32>, pool: ThreadPool, n_jobs: u32) -> Result<Receiver<String>, Box<(dyn Error)>>{
    let (tx, rx) = channel();
    for _ in 0..n_jobs {
        let tx = tx.clone();
        let n = rx_gen.recv().unwrap();
        pool.execute(move || {
            let parse_result = parse(n).unwrap();
            tx.send(parse_result)
                .expect("channel will be there waiting for the pool");
        });
    }
    Ok(rx)
}

fn fan_in(rx_fan_out: Receiver<String>, n_jobs: u32) -> io::Result<()> {

    let stats: Vec<String> = rx_fan_out.iter().take(n_jobs as usize).collect();
    println!("{:#?}", stats);

    Ok(())

    //similar as with take
    //drop(tx);
    //assert_eq!(rx.iter().sum::<usize>(), 8);

    // if would be a struct comming on the channel
    //let result = rx.iter().map(|num| {
    //    if let Some(num) = num {
    //        Ok(num)
    //    }else {
    //        Err("our custom message")
    //    }
    //})
    //.collect::<Result<Vec<u32>, ()>>()
    //.expect("unable to get results");

}

fn main() {
    let n_workers = 4;
    let n_jobs = 8;
    let pool = ThreadPool::new(n_workers);

    let rx_gen = generator(n_jobs).unwrap();
    let rx_fan_out = fan_out(rx_gen, pool, n_jobs).unwrap();
    fan_in(rx_fan_out, n_jobs).unwrap();
}
