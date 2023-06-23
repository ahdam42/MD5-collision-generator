package main

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/csv"
    "fmt"
    "strconv"
    "io"
    "os"
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

func GetRambowTableElement(text int, c chan RambowTableElement) {
    hash := GetMD5Hash(strconv.Itoa(text))
    
    for i := 1; i < CHAIN_LENGTH; i++ {
        hash = GetMD5Hash(hash)
    }

    c <- RambowTableElement{initialValue:text, finalHash:hash}
}

func ReadCollisionDatabase() ([]RambowTableElement, map[string]bool) {
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

func AddColToCollisionDatabase(column []string)  {
    f, err := os.OpenFile(CSV_DB_FILE_NAME, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

    defer f.Close()

    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    w := csv.NewWriter(f)
    w.Write(column)
    w.Flush()
}

func Generate(threads int) {
    c := make(chan RambowTableElement)
    initialValue := 0
    rambowTableElements, _ := ReadCollisionDatabase()

    for _, rambowTableElement := range rambowTableElements {
        if rambowTableElement.initialValue > initialValue {
            initialValue = rambowTableElement.initialValue
        }
    }

    initialValue++

    fmt.Printf("Generating rainbow table. Start position: %d\n", initialValue)

    for i := 0; i < threads; i++ {
        go GetRambowTableElement(initialValue, c)
        initialValue++
    }

    for randowTableElement := range c {
        AddColToCollisionDatabase([]string{strconv.Itoa(randowTableElement.initialValue), randowTableElement.finalHash})
        go GetRambowTableElement(initialValue, c)
        initialValue++
    }
}

func SearchCollision(text string, isFinishedChain chan bool) {
    hash := GetMD5Hash(text)
    step := 1
    rambowTableElements, hashset := ReadCollisionDatabase()
    
    for {
        if (hashset[hash]) {
            initialValue := 0
            for _, rambowTableElement := range rambowTableElements {
                if rambowTableElement.finalHash == hash {
                    initialValue = rambowTableElement.initialValue
                    break
                }
            }
            fmt.Printf("Some collision has been founded!\n Initial Value: %ds. Step: %d\n", initialValue, step)
            isFinishedChain <- true
        }
        hash = GetMD5Hash(hash)
        step++
    }
}

func main() {
    if os.Args[ARG_OPERATION_INDEX] == "find" {
        args := os.Args[ARG_PARAMETER_INDEX:]
        isFinishedChain := make(chan bool)
        for _, collisionCandidate := range args {
            go SearchCollision(collisionCandidate, isFinishedChain)
            fmt.Printf("Searcher for '%s' has been started\n", collisionCandidate)
        }
        <- isFinishedChain
        
    } else {
        threads, _ := strconv.Atoi(os.Args[ARG_PARAMETER_INDEX])
        Generate(threads)
    }
}