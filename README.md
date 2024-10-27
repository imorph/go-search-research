# Linear vs binary search for small arrays in Go language

## Problem statement

Binary search is typically faster than linear search in sorted arrays, but is that always true? It might depend on the array length. For shorter arrays, linear search could be faster due to factors like caching, prefetching, out-of-order execution, or compiler-level optimizations like auto-vectorization. I want to investigate if there’s a clear threshold where one approach becomes better than the other.

### Approach
* Language: Go
* Data: Sorted slices of float64
* Benchmark: Go’s SearchFloat64s (stdlib) which already uses binary search
* Hypothesis: For arrays below a certain length, linear search should outperform binary search. The goal is to pinpoint this threshold and check if it’s consistent across different hardware or if it fluctuates.

### Steps
* Write and optimize a linear search implementation for float64, pushing for maximum performance.
* Set up benchmarks to compare the optimized linear search with Go’s binary search.
* Test a simple linear search against a more optimized one for additional insights.
* Explore how each implementation performs in a highly parallel environment.

Additional Questions:
* How stable is the threshold? Does it depend heavily on hardware specifics?
* Can we generalize any findings to other data types or languages?

## Final Conclusions

### Linear vs. Binary Search (Serial Execution)

#### Small Arrays (n ≤ 30-40)

* **Linear Search Wins**: For small arrays, linear search consistently beats binary search across all positions -- whether the target is at the beginning, middle, end, or not in the array at all.
* **Speed**: Linear search has lower nanoseconds per operation (ns/op) compared to binary, often by a noticeable margin.
* **Best for Early Targets**: The benefit is clearest when the target is near the beginning.

#### Medium Arrays (40 < n ≤ 500)

* **Gap Narrows**: As arrays get bigger, the difference between linear and binary search shrinks.
* **Depends on Target Position**: Linear search still does better if the target is near the start, but binary starts to take over when the target’s in the middle or at the end.

#### Large Arrays (n ≥ 500)

* **Binary Search Dominates**: For big arrays, binary search generally wins, especially when the target isn’t near the beginning.
* **Efficiency Kicks In**: The advantage of binary’s logarithmic time complexity becomes clear as the array size increases.

### Optimized vs. Basic Linear Search

#### Almost No Difference

* **Barely Any Gain**: The optimized linear search (with loop unrolling) shows only tiny improvements over the basic linear search.
* **Small Arrays**: In some cases, the basic linear search even outperforms the optimized version for small arrays and targets near the start.

#### Code Complexity vs. Benefit

* **Keep It Simple**: The basic linear search is simpler to read and maintain.
* **Verdict**: Given the tiny gains, the extra complexity of the optimized version isn’t really worth it.

### Impact of Target Position

#### Target Near the Start

* **Linear Search Shines**: When the target is near the start, linear search is super-efficient, way faster than binary search, no matter the array size.
* **Binary Search Overhead**: Binary doesn’t benefit from early positions in the same way.

#### Middle or End of the Array

* **Linear Slows Down**: Linear search gets slower as the target position moves toward the middle or end.
* **Binary Stays Steady**: Binary search keeps a pretty consistent performance regardless of position.

#### Target Not Found

* **Surprise Advantage**: Interestingly, linear search is fastest when the target isn’t in the array at all, likely because it can end early due to the sorted order.
* **Binary Search Consistency**: Binary search is consistent here but still slower than linear for small arrays in this case.

### Parallel Execution with Higher Concurrency

* **Both Slow Down**: Linear and binary search both see increased execution time per operation (ns/op) as concurrency levels go up (1 to 16 goroutines).

### Picking the Right Algorithm

* **Small Arrays and Early Targets**: Basic linear search is best for small arrays, especially when targets are often near the beginning.
* **Larger Arrays or Uncertain Targets**: For larger arrays or random target positions, binary search is the way to go.
* **Simplicity Wins**: Stick with simple code unless profiling shows a clear performance issue.
* **Concurrency Implications**: In highly concurrent environments, the difference between linear and binary search has less impact on performance.

Testing was done on an Intel Xeon Platinum 8488C CPU, and while hardware can vary, these trends likely apply to other modern CPUs.

## Raw benchmark data

```
cpu: Intel(R) Xeon(R) Platinum 8488C
BenchmarkSearchFunctions/Linear/n=10/pos=beginning-8         	547200334	        2.209 ns/op
BenchmarkSearchFunctions/Binary/n=10/pos=beginning-8         	100000000	       10.67 ns/op
BenchmarkSearchFunctions/Linear/n=10/pos=middle-8            	300476539	        3.974 ns/op
BenchmarkSearchFunctions/Binary/n=10/pos=middle-8            	127543168	        9.442 ns/op
BenchmarkSearchFunctions/Linear/n=10/pos=end-8               	209386896	        5.717 ns/op
BenchmarkSearchFunctions/Binary/n=10/pos=end-8               	126691474	        9.449 ns/op
BenchmarkSearchFunctions/Linear/n=10/pos=notfound-8          	615317810	        1.952 ns/op
BenchmarkSearchFunctions/Binary/n=10/pos=notfound-8          	122716440	        9.706 ns/op
BenchmarkSearchFunctions/Linear/n=20/pos=beginning-8         	480896254	        2.488 ns/op
BenchmarkSearchFunctions/Binary/n=20/pos=beginning-8         	100000000	       10.67 ns/op
BenchmarkSearchFunctions/Linear/n=20/pos=middle-8            	185427771	        6.445 ns/op
BenchmarkSearchFunctions/Binary/n=20/pos=middle-8            	98809566	       11.97 ns/op
BenchmarkSearchFunctions/Linear/n=20/pos=end-8               	100000000	       10.07 ns/op
BenchmarkSearchFunctions/Binary/n=20/pos=end-8               	83167263	       14.41 ns/op
BenchmarkSearchFunctions/Linear/n=20/pos=notfound-8          	612644331	        1.904 ns/op
BenchmarkSearchFunctions/Binary/n=20/pos=notfound-8          	100000000	       11.78 ns/op
BenchmarkSearchFunctions/Linear/n=30/pos=beginning-8         	420927338	        2.836 ns/op
BenchmarkSearchFunctions/Binary/n=30/pos=beginning-8         	86549244	       13.80 ns/op
BenchmarkSearchFunctions/Linear/n=30/pos=middle-8            	133975551	        8.959 ns/op
BenchmarkSearchFunctions/Binary/n=30/pos=middle-8            	80062736	       14.83 ns/op
BenchmarkSearchFunctions/Linear/n=30/pos=end-8               	81175699	       14.64 ns/op
BenchmarkSearchFunctions/Binary/n=30/pos=end-8               	82513237	       14.43 ns/op
BenchmarkSearchFunctions/Linear/n=30/pos=notfound-8          	619382228	        1.906 ns/op
BenchmarkSearchFunctions/Binary/n=30/pos=notfound-8          	100000000	       11.79 ns/op
BenchmarkSearchFunctions/Linear/n=50/pos=beginning-8         	302574694	        3.983 ns/op
BenchmarkSearchFunctions/Binary/n=50/pos=beginning-8         	76876615	       15.65 ns/op
BenchmarkSearchFunctions/Linear/n=50/pos=middle-8            	84785664	       14.14 ns/op
BenchmarkSearchFunctions/Binary/n=50/pos=middle-8            	81054819	       14.75 ns/op
BenchmarkSearchFunctions/Linear/n=50/pos=end-8               	49784884	       23.95 ns/op
BenchmarkSearchFunctions/Binary/n=50/pos=end-8               	82908088	       14.38 ns/op
BenchmarkSearchFunctions/Linear/n=50/pos=notfound-8          	611169012	        1.926 ns/op
BenchmarkSearchFunctions/Binary/n=50/pos=notfound-8          	90476612	       13.42 ns/op
BenchmarkSearchFunctions/Linear/n=100/pos=beginning-8        	186533698	        6.439 ns/op
BenchmarkSearchFunctions/Binary/n=100/pos=beginning-8        	68561656	       17.51 ns/op
BenchmarkSearchFunctions/Linear/n=100/pos=middle-8           	44603016	       26.90 ns/op
BenchmarkSearchFunctions/Binary/n=100/pos=middle-8           	68395600	       17.48 ns/op
BenchmarkSearchFunctions/Linear/n=100/pos=end-8              	25502226	       47.05 ns/op
BenchmarkSearchFunctions/Binary/n=100/pos=end-8              	64132982	       18.75 ns/op
BenchmarkSearchFunctions/Linear/n=100/pos=notfound-8         	617374749	        1.937 ns/op
BenchmarkSearchFunctions/Binary/n=100/pos=notfound-8         	79222620	       15.05 ns/op
BenchmarkSearchFunctions/Linear/n=200/pos=beginning-8        	100000000	       11.65 ns/op
BenchmarkSearchFunctions/Binary/n=200/pos=beginning-8        	61490961	       19.36 ns/op
BenchmarkSearchFunctions/Linear/n=200/pos=middle-8           	22332049	       59.51 ns/op
BenchmarkSearchFunctions/Binary/n=200/pos=middle-8           	60781104	       19.88 ns/op
BenchmarkSearchFunctions/Linear/n=200/pos=end-8              	12888156	       92.87 ns/op
BenchmarkSearchFunctions/Binary/n=200/pos=end-8              	54642687	       21.89 ns/op
BenchmarkSearchFunctions/Linear/n=200/pos=notfound-8         	613246635	        1.941 ns/op
BenchmarkSearchFunctions/Binary/n=200/pos=notfound-8         	72052268	       17.06 ns/op
BenchmarkSearchFunctions/Linear/n=500/pos=beginning-8        	44595696	       26.93 ns/op
BenchmarkSearchFunctions/Binary/n=500/pos=beginning-8        	52778299	       22.90 ns/op
BenchmarkSearchFunctions/Linear/n=500/pos=middle-8           	9241406	      129.2 ns/op
BenchmarkSearchFunctions/Binary/n=500/pos=middle-8           	52700864	       22.40 ns/op
BenchmarkSearchFunctions/Linear/n=500/pos=end-8              	5194936	      231.0 ns/op
BenchmarkSearchFunctions/Binary/n=500/pos=end-8              	49014163	       24.36 ns/op
BenchmarkSearchFunctions/Linear/n=500/pos=notfound-8         	613993623	        1.940 ns/op
BenchmarkSearchFunctions/Binary/n=500/pos=notfound-8         	67244574	       18.63 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=beginning/conc=1-8 	2489484	      484.7 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=beginning/conc=1-8 	2463798	      492.6 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=beginning/conc=2-8 	1623634	      733.9 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=beginning/conc=2-8 	1591402	      755.0 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=beginning/conc=4-8 	 892212	     1249 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=beginning/conc=4-8 	 884155	     1291 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=beginning/conc=8-8 	 489027	     2269 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=beginning/conc=8-8 	 494328	     2320 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=beginning/conc=16-8         	 275502	     4224 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=beginning/conc=16-8         	 263853	     4401 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=middle/conc=1-8             	2473491	      483.2 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=middle/conc=1-8             	2451406	      489.0 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=middle/conc=2-8             	1632957	      739.4 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=middle/conc=2-8             	1602012	      748.3 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=middle/conc=4-8             	 906898	     1252 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=middle/conc=4-8             	 953086	     1266 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=middle/conc=8-8             	 492516	     2264 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=middle/conc=8-8             	 511858	     2330 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=middle/conc=16-8            	 277383	     4242 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=middle/conc=16-8            	 269838	     4316 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=end/conc=1-8                	2489995	      484.0 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=end/conc=1-8                	2455897	      485.2 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=end/conc=2-8                	1638370	      737.6 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=end/conc=2-8                	1603390	      750.0 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=end/conc=4-8                	 881224	     1249 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=end/conc=4-8                	 880746	     1284 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=end/conc=8-8                	 507787	     2253 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=end/conc=8-8                	 488414	     2294 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=end/conc=16-8               	 271400	     4219 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=end/conc=16-8               	 265344	     4364 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=notfound/conc=1-8           	2459294	      484.5 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=notfound/conc=1-8           	2453449	      492.6 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=notfound/conc=2-8           	1627190	      733.5 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=notfound/conc=2-8           	1572800	      754.3 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=notfound/conc=4-8           	 916003	     1240 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=notfound/conc=4-8           	 911833	     1270 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=notfound/conc=8-8           	 513288	     2235 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=notfound/conc=8-8           	 495081	     2328 ns/op
BenchmarkParallelSearches/Linear/n=10/pos=notfound/conc=16-8          	 280911	     4177 ns/op
BenchmarkParallelSearches/Binary/n=10/pos=notfound/conc=16-8          	 270358	     4379 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=beginning/conc=1-8          	2487918	      480.7 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=beginning/conc=1-8          	2454782	      488.7 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=beginning/conc=2-8          	1636348	      744.5 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=beginning/conc=2-8          	1594114	      756.7 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=beginning/conc=4-8          	 947815	     1255 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=beginning/conc=4-8          	 884289	     1290 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=beginning/conc=8-8          	 534913	     2232 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=beginning/conc=8-8          	 475543	     2322 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=beginning/conc=16-8         	 268023	     4271 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=beginning/conc=16-8         	 272548	     4336 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=middle/conc=1-8             	2454813	      483.9 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=middle/conc=1-8             	2417991	      492.5 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=middle/conc=2-8             	1614637	      744.5 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=middle/conc=2-8             	1584153	      760.1 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=middle/conc=4-8             	 909252	     1258 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=middle/conc=4-8             	 880182	     1293 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=middle/conc=8-8             	 488300	     2293 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=middle/conc=8-8             	 476864	     2334 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=middle/conc=16-8            	 272996	     4272 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=middle/conc=16-8            	 272316	     4435 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=end/conc=1-8                	2458165	      489.2 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=end/conc=1-8                	2427418	      494.1 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=end/conc=2-8                	1610809	      749.0 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=end/conc=2-8                	1574347	      757.0 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=end/conc=4-8                	 942980	     1268 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=end/conc=4-8                	 896056	     1293 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=end/conc=8-8                	 486682	     2238 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=end/conc=8-8                	 497738	     2362 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=end/conc=16-8               	 263796	     4324 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=end/conc=16-8               	 270639	     4420 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=notfound/conc=1-8           	2474901	      481.7 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=notfound/conc=1-8           	2448562	      494.5 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=notfound/conc=2-8           	1622848	      735.0 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=notfound/conc=2-8           	1583972	      756.4 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=notfound/conc=4-8           	 927669	     1231 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=notfound/conc=4-8           	 884816	     1291 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=notfound/conc=8-8           	 514773	     2197 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=notfound/conc=8-8           	 454634	     2328 ns/op
BenchmarkParallelSearches/Linear/n=20/pos=notfound/conc=16-8          	 269982	     4257 ns/op
BenchmarkParallelSearches/Binary/n=20/pos=notfound/conc=16-8          	 255736	     4486 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=beginning/conc=1-8          	2467827	      486.2 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=beginning/conc=1-8          	2421727	      496.3 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=beginning/conc=2-8          	1614805	      741.7 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=beginning/conc=2-8          	1595362	      759.0 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=beginning/conc=4-8          	 943063	     1243 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=beginning/conc=4-8          	 889010	     1291 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=beginning/conc=8-8          	 486642	     2240 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=beginning/conc=8-8          	 476979	     2366 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=beginning/conc=16-8         	 278638	     4269 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=beginning/conc=16-8         	 266744	     4446 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=middle/conc=1-8             	2453551	      489.4 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=middle/conc=1-8             	2443285	      494.3 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=middle/conc=2-8             	1612738	      745.9 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=middle/conc=2-8             	1582786	      757.8 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=middle/conc=4-8             	 951526	     1269 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=middle/conc=4-8             	 862189	     1299 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=middle/conc=8-8             	 476094	     2308 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=middle/conc=8-8             	 499282	     2341 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=middle/conc=16-8            	 272356	     4351 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=middle/conc=16-8            	 268182	     4398 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=end/conc=1-8                	2438353	      494.9 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=end/conc=1-8                	2416608	      493.9 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=end/conc=2-8                	1584787	      752.4 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=end/conc=2-8                	1582602	      762.2 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=end/conc=4-8                	 895069	     1291 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=end/conc=4-8                	 917103	     1297 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=end/conc=8-8                	 473794	     2335 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=end/conc=8-8                	 499837	     2354 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=end/conc=16-8               	 260761	     4340 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=end/conc=16-8               	 263390	     4502 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=notfound/conc=1-8           	2463843	      482.2 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=notfound/conc=1-8           	2440993	      493.7 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=notfound/conc=2-8           	1625569	      736.7 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=notfound/conc=2-8           	1582402	      757.5 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=notfound/conc=4-8           	 946311	     1254 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=notfound/conc=4-8           	 917758	     1297 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=notfound/conc=8-8           	 538237	     2224 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=notfound/conc=8-8           	 480237	     2363 ns/op
BenchmarkParallelSearches/Linear/n=30/pos=notfound/conc=16-8          	 276534	     4218 ns/op
BenchmarkParallelSearches/Binary/n=30/pos=notfound/conc=16-8          	 259818	     4381 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=beginning/conc=1-8          	2480029	      483.9 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=beginning/conc=1-8          	2424061	      490.9 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=beginning/conc=2-8          	1639674	      735.6 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=beginning/conc=2-8          	1581430	      766.0 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=beginning/conc=4-8          	 897920	     1248 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=beginning/conc=4-8          	 909123	     1301 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=beginning/conc=8-8          	 530450	     2253 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=beginning/conc=8-8          	 465669	     2334 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=beginning/conc=16-8         	 267950	     4274 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=beginning/conc=16-8         	 260677	     4475 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=middle/conc=1-8             	2436864	      484.6 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=middle/conc=1-8             	2414734	      491.3 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=middle/conc=2-8             	1609281	      749.9 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=middle/conc=2-8             	1579686	      764.7 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=middle/conc=4-8             	 944019	     1275 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=middle/conc=4-8             	 895980	     1289 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=middle/conc=8-8             	 516668	     2301 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=middle/conc=8-8             	 487382	     2358 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=middle/conc=16-8            	 269998	     4304 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=middle/conc=16-8            	 263402	     4377 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=end/conc=1-8                	2424722	      495.9 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=end/conc=1-8                	2431542	      493.9 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=end/conc=2-8                	1549137	      763.2 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=end/conc=2-8                	1583527	      754.8 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=end/conc=4-8                	 898783	     1322 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=end/conc=4-8                	 904975	     1283 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=end/conc=8-8                	 473643	     2377 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=end/conc=8-8                	 492014	     2329 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=end/conc=16-8               	 260604	     4488 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=end/conc=16-8               	 263364	     4375 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=notfound/conc=1-8           	2480556	      481.0 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=notfound/conc=1-8           	2435042	      498.4 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=notfound/conc=2-8           	1638276	      726.8 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=notfound/conc=2-8           	1581224	      759.5 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=notfound/conc=4-8           	 944505	     1238 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=notfound/conc=4-8           	 854727	     1306 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=notfound/conc=8-8           	 497780	     2239 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=notfound/conc=8-8           	 472665	     2335 ns/op
BenchmarkParallelSearches/Linear/n=40/pos=notfound/conc=16-8          	 272766	     4215 ns/op
BenchmarkParallelSearches/Binary/n=40/pos=notfound/conc=16-8          	 258199	     4482 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=beginning/conc=1-8          	2442326	      485.5 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=beginning/conc=1-8          	2410635	      496.5 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=beginning/conc=2-8          	1622883	      735.1 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=beginning/conc=2-8          	1568824	      768.4 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=beginning/conc=4-8          	 905330	     1259 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=beginning/conc=4-8          	 844674	     1295 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=beginning/conc=8-8          	 494224	     2261 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=beginning/conc=8-8          	 475452	     2373 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=beginning/conc=16-8         	 270141	     4290 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=beginning/conc=16-8         	 264853	     4443 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=middle/conc=1-8             	2424339	      493.6 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=middle/conc=1-8             	2427396	      494.6 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=middle/conc=2-8             	1603596	      754.3 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=middle/conc=2-8             	1578926	      759.5 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=middle/conc=4-8             	 872170	     1280 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=middle/conc=4-8             	 885121	     1293 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=middle/conc=8-8             	 494847	     2324 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=middle/conc=8-8             	 487958	     2308 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=middle/conc=16-8            	 267386	     4427 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=middle/conc=16-8            	 268192	     4384 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=end/conc=1-8                	2395981	      501.0 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=end/conc=1-8                	2430667	      493.0 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=end/conc=2-8                	1541818	      783.1 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=end/conc=2-8                	1582297	      755.5 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=end/conc=4-8                	 888759	     1339 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=end/conc=4-8                	 883659	     1293 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=end/conc=8-8                	 463831	     2458 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=end/conc=8-8                	 493922	     2342 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=end/conc=16-8               	 253958	     4555 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=end/conc=16-8               	 263998	     4444 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=notfound/conc=1-8           	2492644	      483.6 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=notfound/conc=1-8           	2427432	      495.9 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=notfound/conc=2-8           	1620042	      738.9 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=notfound/conc=2-8           	1568023	      758.7 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=notfound/conc=4-8           	 914716	     1243 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=notfound/conc=4-8           	 860804	     1291 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=notfound/conc=8-8           	 510148	     2200 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=notfound/conc=8-8           	 483796	     2372 ns/op
BenchmarkParallelSearches/Linear/n=50/pos=notfound/conc=16-8          	 268788	     4220 ns/op
BenchmarkParallelSearches/Binary/n=50/pos=notfound/conc=16-8          	 262425	     4442 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=10/pos=beginning-8         	515687382	        2.304 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=10/pos=beginning-8             	718238126	        1.620 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=10/pos=beginning-8            	100000000	       10.63 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=10/pos=middle-8            	299125765	        4.003 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=10/pos=middle-8                	262081933	        4.420 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=10/pos=middle-8               	127216148	        9.456 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=10/pos=end-8               	209739775	        5.744 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=10/pos=end-8                   	220924279	        5.388 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=10/pos=end-8                  	127434408	        9.405 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=10/pos=notfound-8          	584832828	        2.046 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=10/pos=notfound-8              	686669466	        1.744 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=10/pos=notfound-8             	123916125	        9.679 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=20/pos=beginning-8         	453475995	        2.656 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=20/pos=beginning-8             	417808647	        2.929 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=20/pos=beginning-8            	100000000	       10.72 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=20/pos=middle-8            	185765174	        6.444 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=20/pos=middle-8                	196453834	        6.044 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=20/pos=middle-8               	100000000	       12.00 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=20/pos=end-8               	100000000	       10.10 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=20/pos=end-8                   	132825579	        9.029 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=20/pos=end-8                  	81776875	       14.46 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=20/pos=notfound-8          	593005045	        2.016 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=20/pos=notfound-8              	746257814	        1.519 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=20/pos=notfound-8             	100000000	       11.89 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=30/pos=beginning-8         	404196260	        2.954 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=30/pos=beginning-8             	322860038	        3.792 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=30/pos=beginning-8            	87959095	       13.72 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=30/pos=middle-8            	133928763	        8.973 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=30/pos=middle-8                	148528254	        8.063 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=30/pos=middle-8               	80040516	       14.83 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=30/pos=end-8               	82398397	       14.62 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=30/pos=end-8                   	89753682	       13.29 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=30/pos=end-8                  	84615033	       14.43 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=30/pos=notfound-8          	598048832	        2.084 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=30/pos=notfound-8              	742503250	        1.731 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=30/pos=notfound-8             	100000000	       11.59 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=40/pos=beginning-8         	330284442	        3.593 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=40/pos=beginning-8             	253683132	        4.858 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=40/pos=beginning-8            	77249092	       15.64 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=40/pos=middle-8            	100000000	       11.63 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=40/pos=middle-8                	100000000	       10.44 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=40/pos=middle-8               	79686306	       14.83 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=40/pos=end-8               	61778659	       19.15 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=40/pos=end-8                   	67275985	       17.57 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=40/pos=end-8                  	80791954	       14.89 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=40/pos=notfound-8          	582033614	        2.098 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=40/pos=notfound-8              	707713222	        1.867 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=40/pos=notfound-8             	87756196	       13.54 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=50/pos=beginning-8         	300006421	        3.998 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=50/pos=beginning-8             	291814297	        4.425 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=50/pos=beginning-8            	77527852	       15.57 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=50/pos=middle-8            	84818540	       14.16 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=50/pos=middle-8                	93544887	       12.84 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=50/pos=middle-8               	79923805	       14.87 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=50/pos=end-8               	50033209	       23.96 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=50/pos=end-8                   	54785232	       21.88 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=50/pos=end-8                  	79579575	       14.89 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=50/pos=notfound-8          	588881270	        2.063 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=50/pos=notfound-8              	714136119	        1.694 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=50/pos=notfound-8             	88312768	       13.44 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=60/pos=beginning-8         	272010813	        4.423 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=60/pos=beginning-8             	255304492	        4.972 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=60/pos=beginning-8            	76217299	       15.69 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=60/pos=middle-8            	71573947	       16.76 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=60/pos=middle-8                	78376677	       15.25 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=60/pos=middle-8               	79257802	       14.89 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=60/pos=end-8               	41746408	       28.57 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=60/pos=end-8                   	45667957	       26.15 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=60/pos=end-8                  	70865548	       16.88 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=60/pos=notfound-8          	570769028	        2.004 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=60/pos=notfound-8              	734117944	        1.714 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=60/pos=notfound-8             	89720463	       13.52 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=100/pos=beginning-8        	185679102	        6.474 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=100/pos=beginning-8            	198400063	        6.138 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=100/pos=beginning-8           	67686448	       17.42 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=100/pos=middle-8           	44524358	       26.92 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=100/pos=middle-8               	48821330	       24.74 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=100/pos=middle-8              	68553262	       17.40 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=100/pos=end-8              	25332984	       47.24 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=100/pos=end-8                  	27416332	       43.50 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=100/pos=end-8                 	63430437	       18.78 ns/op
BenchmarkLinearSearchImplementations/OptimizedLinear/n=100/pos=notfound-8         	594794934	        2.033 ns/op
BenchmarkLinearSearchImplementations/BasicLinear/n=100/pos=notfound-8             	733113880	        1.801 ns/op
BenchmarkLinearSearchImplementations/BinarySearch/n=100/pos=notfound-8            	83311723	       14.92 ns/op

```
