# Wordle
> An async Wordle game solely on p2p.

## What
A fully distributed p2p implementation of the popular Wordle game with a few amendments to the rules. Instead of 
globally defined(by centralized authority) set of words that change every 24 hours, we implement a simple asynchronous 
'first come, first serve' consensus where a guesser becomes a new leader and chooses the next word for the whole network
to guess.

## Why
This project is a [hackaton](https://p2p.paris/en/event/hackathon-1/) submission written in two days. The original idea
was to implement something based on the [Lazy Client for Lazy Blockchains paper](https://arxiv.org/abs/2203.15968), and 
we decided to implement a game with global state transitions, like Wordle.

## Play
* Clone it
* `make build` it
* `./build/wordle light start` it
* Wait until you discover peers
* Play it

## Team
* [@Wondertan](https://github.com/Wondertan)
* [@cortze](https://github.com/cortze)
* [@ajnavarro](https://github.com/ajnavarro)
