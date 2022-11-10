## Benchmarks

We follow [the method](https://openscoring.io/blog/2021/08/04/benchmarking_sklearn_jpmml_evaluator/) outlined by Openscoring for evaluation.

The best reported numbers for JPMML-Evaluator are `4.190 — 4.770` microseconds (or 4190 — 4770 nanoseconds, as 1 microsecond = 1000 nanoseconds).

Using an identical model from the created flow, Pummel is about 3x faster, clocking around ~1550 ns per iteration.

```
pkg: github.com/stillmatic/pummel/testdata
BenchmarkAuditLR-10                	  774922	      1553 ns/op	     797 B/op	      19 allocs/op
BenchmarkAuditLRConcurrently-10    	  934758	      1184 ns/op	     853 B/op	      21 allocs/op
BenchmarkAuditRF-10                	    9090	    124334 ns/op	   34987 B/op	     516 allocs/op
BenchmarkAuditRFConcurrently-10    	   56708	     20299 ns/op	   27763 B/op	     361 allocs/op
PASS
```

Alternatively, we can consider a benchmark where we instead process 10k rows in a row (to control for warming up effect). That yields a similar 1556 ns/row. In particular, it means we do not have cold start problems (assuming the model is kept in memory).

```
BenchmarkAuditLR-10    	      75	  15597657 ns/op	 7977965 B/op	  192148 allocs/op
```

This is also likely comparable to the Python SKLearn timings, though there can be a difference in the benchmarking. The Python end-to-end pipeline takes ~1.7 microseconds in the reported benchmarks, and this Go code similarly takes ~1.55 nanoseconds. Our data processing has higher overhead than the Python code, which does not need type changes. 

The Random Forest shows a slightly slower performance than the JPMML numbers (120 microseconds vs 85 microseconds) but makes up for it with concurrency, being able to drop to 20 microseconds, about 4x faster.

