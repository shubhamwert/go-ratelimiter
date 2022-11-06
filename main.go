package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type bucketType struct {
	apiKey string
	bucket *RLRequestBucket
}

var STATIC_API_KEYS map[string]struct{}
var buckets map[string]bucketType = make(map[string]bucketType)

func main() {
	STATIC_API_KEYS := make(map[string]struct{})

	for _, i := range []string{"1", "2", "3"} {

		b, err := CreateNewBucket(3)

		if err != nil {
			fmt.Println("Failed to create bucket number ", i)
			continue
		}
		buckets[i] = bucketType{apiKey: i, bucket: b}
		fmt.Println("Starting rate limiters ", i)
		STATIC_API_KEYS[i] = struct{}{}

	}
	router := mux.NewRouter()
	i := 0
	router.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// go HandelRateLimiter(w, r, i)
		// func HandelRateLimiter(w http.ResponseWriter, r *http.Request, i int) {
		fmt.Println("#### ", i)
		defer fmt.Println("#### ", i)
		i += 1
		fmt.Println(time.Now(), "Request by ", r.Header)
		values := r.URL.Query()
		ApiKey := values.Get("api-key")
		if len(ApiKey) == 0 {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(401)
			return
		}
		_, ok := STATIC_API_KEYS[ApiKey]

		if !ok {
			// fmt.Println(STATIC_API_KEYS.])
			for k, _ := range STATIC_API_KEYS {
				println(k, " ++")
			}
			fmt.Println("API-Key==> ", ApiKey)
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(401)
			return
		}
		fmt.Println("Processing for api-key ==> ", buckets[ApiKey].bucket.currentCapacity)
		if bucket, ok := buckets[ApiKey]; ok {
			if _, ok := bucket.bucket.Request(); ok == nil {

				// create a new url from the raw ruestURI sent by the client
				url := "https://google.com"
				proxyReq, err := http.NewRequest(r.Method, url, r.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				// We may want to filter some headers, otherwise we could just use a shallow copy
				// proxyReq.Header = req.Header
				proxyReq.Header.Set("Host", r.Host)
				proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

				for header, values := range r.Header {
					for _, value := range values {
						proxyReq.Header.Add(header, value)
					}
				}
				client := &http.Client{}
				resp, err := client.Do(proxyReq)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer resp.Body.Close()
				copyHeader(w.Header(), resp.Header)
				w.WriteHeader(resp.StatusCode)
				io.Copy(w, resp.Body)
				resp.Body.Close()
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Println("Too many requests slow down")
				return
			}

			// To create a delay
			fmt.Println("Sleeping for ---->", (10)*int(time.Second))
			time.Sleep(time.Duration((20) * int(time.Second)))

		} else {
			fmt.Println("Key Not Found")
			fmt.Println(buckets["ApiKey"])

		}

		// }
	})
	router.HandleFunc("/parent", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Forwarded")
		fmt.Println(r.Header)
	})
	http.Handle("/", router)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

// https://pkg.go.dev/golang.org/x/time/rate
