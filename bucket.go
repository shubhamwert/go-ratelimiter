package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type RLRequestBucket struct {
	mu              sync.Mutex
	netCapacity     int
	currentCapacity int
	isAlive         bool
}

func (bucket *RLRequestBucket) start(ti time.Duration) (int, error) {
	// this process might even starve
	fmt.Println("Starting the bucket filler")
	for bucket.isAlive {
		time.Sleep(time.Duration(ti * time.Second))
		if c, err := bucket.Fill(); err != nil {
			fmt.Println("Filled bucket, new size for bucket is ", c)
		}

	}
	fmt.Println("Terminating the bucket filler")

	return 0, nil
}

func CreateNewBucket(netCapacity int) (*RLRequestBucket, error) {

	bucket := &RLRequestBucket{netCapacity: netCapacity, currentCapacity: netCapacity, isAlive: true}
	go bucket.start(5)
	return bucket, nil
}

func (bucket *RLRequestBucket) Request() (int, error) {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	if bucket.currentCapacity <= 0 {
		return bucket.currentCapacity, errors.New("bucket is already full")
	}
	bucket.currentCapacity -= 1

	return bucket.currentCapacity, nil
}
func (bucket *RLRequestBucket) Fill() (int, error) {
	if bucket.currentCapacity > bucket.netCapacity {
		return bucket.currentCapacity, nil
	}
	bucket.currentCapacity += 1

	return bucket.currentCapacity, nil
}
func (bucket *RLRequestBucket) Completed() (int, error) {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()
	if bucket.currentCapacity > bucket.netCapacity {
		return bucket.currentCapacity, errors.New("something is fishy! ")
	}
	bucket.currentCapacity += 1

	return bucket.currentCapacity, nil
}
