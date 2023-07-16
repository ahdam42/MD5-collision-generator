# MD5-collision-generator
Multi-threaded md5 collision generator using rainbow tables on Go

## Idea
The main idea is to generate extra-long hash chains and make a rainbow table. After generating a large enough database of hash chains, it would theoretically be possible to find a collision of the form "7215ee9c7d9dc229d2921a40e899ec5f" for any value. Because of the birthday paradox, there is no need to compute all 2^128 variants, and the rainbow table allows you to keep the data in a compact form. An alternative search variant using Floyd's algorithm is also implemented. 

## How to build
    go build CollisionChecker.go

## How to run

### Collision search using a pre-generic table using 16 threads. Due to the low generation speed on CPU, the table is generated on GPU, but more data is still needed.
    ./CollisionChecker rainbowTableSearcher 16 
### Collision search using Floyd's algorithm using 16 threads
    ./CollisionChecker floydSearcher 16
