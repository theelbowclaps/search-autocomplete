package main

import (
	"errors"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Declare a pool variable to hold the pool of Redis connections.
var pool *redis.Pool

// ErrNoWordFound - triggered when no hits
var ErrNoWordFound = errors.New("no words found, updating search database")

// FindWord looks for keywords registered in redis
func FindWord(searchRequest SearchRequest) ([]string, error) {
	conn := pool.Get()
	defer conn.Close()

	keyword := searchRequest.Prefixed
	confirmed := searchRequest.Confirmed
	if confirmed {
		InsertWord(conn, keyword)
	}
	var results []string
	var grab int64 = 10
	count := 50

	rank, err := redis.Int64(conn.Do("ZRANK", "mylist", keyword))
	if err != nil {
		log.Println("No keywords found")
		return nil, nil
	}

	for len(results) != count {
		found, err := redis.Strings(conn.Do("ZRANGE", "mylist", rank, rank+grab-1))
		rank = rank + grab
		if err != nil {
			log.Println("ERROR while getting range")
			return nil, err
		}
		if len(found) == 0 {
			break
		}
		for _, entry := range found {
			var minLen int
			if len(entry) < len(keyword) {
				minLen = len(entry)
			} else {
				minLen = len(keyword)
			}
			if entry[:minLen] != keyword[:minLen] {
				count = len(results)
				break
			}
			if string(entry[len(entry)-1]) == "%" && len(results) != count {
				results = append(results, entry[0:len(entry)-1])
			}
		}

	}
	return results, nil
}

// InsertWord when new word detected
func InsertWord(conn redis.Conn, keyword string) error {
	s := keyword

	for i := range s[:len(s)] {
		newWord := s[:i+1]
		if i == len(s)-1 {
			newWord = newWord + "%"
		}
		_, err := conn.Do("ZADD", "mylist", 0, newWord)
		if err != nil {
			return err
		}
	}

	return nil
}
