package main

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/csv"
    "fmt"
    "strconv"
    "io"
    "os"
    "sync"
    "math/rand"
)

type RambowTableElement struct {
    initialValue int
    finalHash string
}

const CHAIN_LENGTH = 4294967297
const CSV_INITIAL_VALUE_INDEX = 0
const CSV_FINAL_HASH_INDEX = 1
const CSV_DB_FILE_NAME = "collision_db.csv"
const ARG_OPERATION_INDEX = 1
const ARG_PARAMETER_INDEX = 2

func GetMD5Hash(text string) string {
   hash := md5.Sum([]byte(text))

   return hex.EncodeToString(hash[:])
}

func ReadRainbowTable() ([]RambowTableElement, map[string]bool) {
    rambowTableElements := []RambowTableElement{};
    hashSet := make(map[string]bool)
    f, _ := os.Open(CSV_DB_FILE_NAME)
    r := csv.NewReader(f)

    defer f.Close()
    
    for {
        record, err := r.Read()

        if err == io.EOF {
            break
        }

        finalHash := record[CSV_FINAL_HASH_INDEX ]
        initiavValue, _ := strconv.Atoi(record[CSV_INITIAL_VALUE_INDEX])
        rambowTableElements = append(rambowTableElements, RambowTableElement{initialValue: initiavValue, finalHash: finalHash})
        hashSet[finalHash] = true
    }

    return rambowTableElements, hashSet
}

func RainbowTableSearcher(text string, rambowTableElements []RambowTableElement, hashset map[string]bool, wg *sync.WaitGroup) {
    defer wg.Done()
    
    hash := GetMD5Hash(text)
    step := 1
    
    for {
        if step > CHAIN_LENGTH {
            fmt.Printf("No collision for '%s' was found.\n", text)
            break;
        }

        if (hashset[hash]) {
            initialValue := 0

            // find initial chain value
            for _, rambowTableElement := range rambowTableElements {
                if rambowTableElement.finalHash == hash {
                    initialValue = rambowTableElement.initialValue
                    break
                }
            }
            fmt.Printf("Some collision has been founded!\n Initial Value: %ds. Step: %d\n", initialValue, step)

            // find collisioned value
            collisionChain := GetMD5Hash(strconv.Itoa(initialValue))
    
            for i := 1; i < CHAIN_LENGTH - step - 1; i++ {
                collisionChain = GetMD5Hash(collisionChain)
            }

            fmt.Printf("Collision for '%s'\n Initial Value: %d. Hash: %s\n", text, hash)

            break;
        }
        hash = GetMD5Hash(hash)
        step++
    }
}


func randSeq(n int) string {
    symbols := []rune("0123456789abcdef")
    randStr := make([]rune, n)
    
    for i := range randStr {
        randStr[i] = symbols[rand.Intn(len(symbols))]
    }
    
    return string(randStr)
}

func FloydCollisionSearcher(initialStr string, wg *sync.WaitGroup) {
    defer wg.Done()

    slowPointer := initialStr
    fastPointer := initialStr
    step := 0;

    for {
        step++
        slowPointer = GetMD5Hash(slowPointer)
        fastPointer = GetMD5Hash(GetMD5Hash(fastPointer))

        if(slowPointer == fastPointer) {
            fmt.Printf("Collision has been founded!\n Meeting point: %s. Step: %d. Initial value: %s. \n", slowPointer, step, initialStr)
            break;
        }
    }

    slowPointer = initialStr

    for {
        prevSlowPointer := slowPointer
        prevFastPointer := fastPointer
        slowPointer = GetMD5Hash(slowPointer)
        fastPointer = GetMD5Hash(fastPointer)

        if(slowPointer == fastPointer) {
            fmt.Printf("Slow pointer: %s. Fast pointer: %s \n", prevSlowPointer, prevFastPointer)
            break;
        }
    }

}

func main() {
    if os.Args[ARG_OPERATION_INDEX] == "floydSearcher" {
        var wg sync.WaitGroup
        threads, _ := strconv.Atoi(os.Args[ARG_PARAMETER_INDEX])

        wg.Add(threads)
        for i := 0; i < threads; i++ {
            randomString := randSeq(32)
            go FloydCollisionSearcher(randomString, &wg)
            fmt.Printf("Searcher for '%s' has been started\n", randomString)
        }
        wg.Wait()
    } else if os.Args[ARG_OPERATION_INDEX] == "rainbowTableSearcher" {
        rambowTableElements, hashset := ReadRainbowTable()
        for {
            var wg sync.WaitGroup
            threads, _ := strconv.Atoi(os.Args[ARG_PARAMETER_INDEX])
            
            wg.Add(threads)
            for i := 0; i < threads; i++ {
                randomString := randSeq(32)
                go RainbowTableSearcher(randomString, rambowTableElements, hashset, &wg)
                fmt.Printf("Searcher for '%s' has been started\n", randomString)
            }
            wg.Wait()
        }
    }
}
