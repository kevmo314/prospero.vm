# prospero.vm

My submission for [Matt Keeter's Prospero challenge](https://www.mattkeeter.com/projects/prospero/).

Runs in 500us/frame for 4096x4096 on an RTX 4090. That's sub 1ms per frame.

![prospero](output.png)

Matt's post mentions

> The description should call out any particularly promising results, clever ideas, or interesting tools!

When building my solution, I am reminded of something I tell myself often: there are thousands of people
much smarter than me who have invested more time than me to solve problems more difficult than mine.

Therefore, I took the ["bitter lesson" approach](http://www.incompleteideas.net/IncIdeas/BitterLesson.html)
and converted the `.vm` file to CUDA code which I then compile. This results in a nice 100-line solution
without needing to invent anything particularly novel.

## Usage

Run the codegen with:

```
go run main.go
```

Compile the CUDA kernel:

```
nvcc -use_fast_math -O3 -prec-sqrt=false prospero.cu
```

Run the kernel, noting that the program runs 1000 iterations, therefore the timing is for 1000 frames.

```
$ time ./a.out > out.ppm

real    0m0.524s
user    0m0.410s
sys     0m0.098s
```

## Notes

If you are curious, because I suppose some might consider hiding the execution in `nvcc` cheating,
`nvcc` takes ~2 seconds to run:

```
$ time nvcc -use_fast_math -O3 -prec-sqrt=false prospero.cu

real    0m2.153s
user    0m1.700s
sys     0m0.203s
```

I don't think it's cheating but I also recognize that a logical extension of that is compiling
the program into `#define`'s which probably would yield more of a speedup. If one wanted to make
my solution more dynamic, I would recommend building a `.vm` to ptx compiler and letting ptxas
do its thing, then manually do the `cudaLaunchKernel()` dance and ship the ptx to the GPU dynamically.

This is a cute trick we do in [SCUDA](https://github.com/kevmo314/scuda).
