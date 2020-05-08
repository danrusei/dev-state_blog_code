use std::io;
use std::thread;
use rand::thread_rng;
use rand::seq::SliceRandom;
use crossbeam_channel::{bounded, Sender, Receiver};

fn generator(iters: u32) -> io::Result<Receiver<u32>> {
    let (gen_sender, gen_receiver) = bounded(0);
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

fn main() {
    
    let receiver = generator(10).unwrap();
    loop {
        if let new = receiver.recv().is_err() {
            break;
        }
        println!("I received {}",new)
    }
}
