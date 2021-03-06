package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
	"net/http"
	"fmt"
	"log"
	"encoding/json"
	"sync"
)

// ComputeResult stores the result on palindrome calculation and the time it took
type ComputeResult struct {
	Number int `json:"Number"`
	Binary string `json:"Binary"`
	Time   int `json:"Time"`
}

var (
	Palindrome_temp *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "Palindrome",
		Help:    "help",
		Buckets: prometheus.DefBuckets,
	}, []string{"base"})
)

// MinerSingle executes the palindrome computations one at a time
func MinerSingle(number int) (ComputeResult, error) {
	// output := make(map[int]string)
	var output ComputeResult

	if isPalindrome(number) {
		if isBinaryPalindrome(number) {
			output.Number = number
			output.Binary = convertToBinary(number)
			return output, nil
		}
	}
	return output, errors.New("Not a palindrome")
}

func isPalindrome(number int) bool {
	// converts int into string and then check if the string is a palindrome

	defer func(begin time.Time) {

		s := time.Since(begin).Seconds()
		ms := s * 1e3
		Palindrome_temp.WithLabelValues("Ten").Observe(ms)

	}(time.Now())

	forwardString := strconv.Itoa(number)
	reversedString := reverse(forwardString)

	return forwardString == reversedString
}

func isBinaryPalindrome(number int) bool {
	// converts int into a string of binary and then check if the string is a palindrome

	defer func(begin time.Time) {

		s := time.Since(begin).Seconds()
		ms := s * 1e3
		Palindrome_temp.WithLabelValues("Two").Observe(ms)

	}(time.Now())

	forwardBinaryString := convertToBinary(number)
	reversedBinaryString := reverse(forwardBinaryString)

	return forwardBinaryString == reversedBinaryString
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func convertToBinary(i int) string {
	i64 := int64(i)
	return strconv.FormatInt(i64, 2) // base 2 for binary
}

type NumNum struct {
    Block int `json:"block"`
}

func handle_single_miner(w http.ResponseWriter, r *http.Request){
	test_num := NumNum{}

	err := json.NewDecoder(r.Body).Decode(&test_num)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

	start := time.Now()
	result, err := MinerSingle(test_num.Block)
	elapsed := time.Since(start)
	result.Time = int(elapsed.Nanoseconds())

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	profile := ComputeResult{result.Number, result.Binary, result.Time}
	js, _ := json.Marshal(profile)

	w.Header().Set("Content-Type", "application/json")
  	w.Write(js)
}

type NumNumNum struct {
    StartBlock int `json:"startBlock"`
	EndBlock int `json:"EndBlock"`
}

func handle_block_miner(w http.ResponseWriter, r *http.Request){
	test_num := NumNumNum{}

	err := json.NewDecoder(r.Body).Decode(&test_num)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

	var allResults []ComputeResult
	var mu = &sync.Mutex{}
	
	for i := test_num.StartBlock; i < test_num.EndBlock; i++{
		start := time.Now()
		result, err := MinerSingle(i)
		elapsed := time.Since(start)
		result.Time = int(elapsed.Nanoseconds())

		if err != nil {
			continue
		}
		timed_result := ComputeResult{result.Number, result.Binary, result.Time}
		mu.Lock()
		allResults = append(allResults, timed_result)
		mu.Unlock()
	}

	js, _ := json.Marshal(allResults)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main(){

	// Create default route handler
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/MinerSingle", handle_single_miner)
	http.HandleFunc("/MinerBlock", handle_block_miner)
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8082", nil))
}