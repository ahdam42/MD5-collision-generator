# MD5-collision-generator
Multi-threaded md5 collision generator using rainbow tables on Go

## How to build
    go build RainbowTableGenerator.go

## How to run

### Adding new entries to the rainbow table using 16 threads
    ./RainbowTableGenerator generate 16
### Trying to find a collision
    ./RainbowTableGenerator find candidate1 candidate2 [...candidateN]
