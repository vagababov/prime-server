/*
Copyright 2018 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Based upon: https://github.com/vagababov/maxprime.
*/
package main

import (
	"math"
	"math/big"
)

func calcPrime(maxNumber int64) int64 {
	// If prime, increment to keep the execution time
	// the result will be the same but the run time for
	// primes won't be artificially short.
	if big.NewInt(int64(maxNumber)).ProbablyPrime(0) {
		maxNumber++
	}

	var x, y, n int64
	nsqrt := math.Sqrt(float64(maxNumber))

	isPrime := make([]bool, maxNumber)
	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= maxNumber && (n%12 == 1 || n%12 == 5) {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) + y*y
			if n <= maxNumber && n%12 == 7 {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= maxNumber && n%12 == 11 {
				isPrime[n] = !isPrime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if isPrime[n] {
			for y = n * n; y < maxNumber; y += n * n {
				isPrime[y] = false
			}
		}
	}

	isPrime[2] = true
	isPrime[3] = true

	primes := make([]int64, 0, 1270606)
	for i := 0; i < len(isPrime)-1; i++ {
		if isPrime[i] {
			primes = append(primes, int64(i))
		}
	}
	return primes[len(primes)-1]
}
