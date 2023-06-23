# MD5-collision-generator
Multi-threaded md5 collision generator using rainbow tables on Go

## Idea
The main idea is to generate extra-long hash chains and make a rainbow table. After generating a large enough database of hash chains, it would theoretically be possible to find a collision of the form "7215ee9c7d9dc229d2921a40e899ec5f" for any value. Because of the birthday paradox, there is no need to compute all 2^128 variants, and the rainbow table allows you to keep the data in a compact form.

## How to build
    go build RainbowTableGenerator.go

## How to run

### Adding new entries to the rainbow table using 16 threads
    ./RainbowTableGenerator generate 16
### Trying to find a collision
    ./RainbowTableGenerator find candidate1 candidate2 [...candidateN]
