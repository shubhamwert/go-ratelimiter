package main

import "errors"

type RLRequestBucket struct {
	netCapacity     int
	currentCapacity int
}

func CreateNewBucket(netCapacity int) (*RLRequestBucket, error) {
	return &RLRequestBucket{netCapacity: netCapacity, currentCapacity: netCapacity}, nil
}

func (bucket *RLRequestBucket) Request() (int, error) {
	if bucket.currentCapacity <= 0 {
		return bucket.currentCapacity, errors.New("bucket is already full")
	}
	bucket.currentCapacity -= 1

	return bucket.currentCapacity, nil
}
func (bucket *RLRequestBucket) Completed() (int, error) {
	if bucket.currentCapacity > bucket.netCapacity {
		return bucket.currentCapacity, errors.New("something is fishy! ")
	}
	bucket.currentCapacity += 1

	return bucket.currentCapacity, nil
}
