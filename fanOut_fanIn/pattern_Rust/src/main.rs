use rand::seq::SliceRandom;
use rand::thread_rng;
use std::io;
use std::thread;
use std::sync::mpsc::{channel, Receiver};
use threadpool::ThreadPool;

fn generator(iters: u32) -> io::Result<Receiver<u32>> {
    let (gen_sender, gen_receiver) = channel();
    let mut vec: Vec<u32> = (1..100).collect();
    vec.shuffle(&mut thread_rng());
    thread::spawn(move || {
        for i in 0..iters {
            if gen_sender.send(vec[i as usize].clone()).is_err() {
                break;
            }
        }
    });

    Ok(gen_receiver)
}

fn fan_out(rx_gen: Receiver<u32>, pool: ThreadPool, n_jobs: u32) -> io::Result<Receiver<u32>>{
    let (tx, rx) = channel();
    for _ in 0..n_jobs {
        let tx = tx.clone();
        let n = rx_gen.recv().unwrap();
        pool.execute(move || {
            tx.send(n)
                .expect("channel will be there waiting for the pool");
        });
    }
    Ok(rx)
}

fn fan_in(rx_fan_out: Receiver<u32>, n_jobs: u32) -> io::Result<()> {

    let stats: Vec<u32> = rx_fan_out.iter().take(n_jobs as usize).collect();
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
