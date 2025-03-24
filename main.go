package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("prospero.vm")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	g, err := os.OpenFile("prospero.cu", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer g.Close()
	if _, err := g.WriteString("#include <cuda_runtime.h>\n#include <stdio.h>\n\n__global__ void prospero(int* zs, int n) {\n"); err != nil {
		panic(err)
	}

	g.WriteString(`	int i = blockIdx.x * blockDim.x + threadIdx.x;
	if (i >= n)
		return;
	float x = (((float)(i % 4096)) / (4096.0f * 0.5f)) - 1.0f;
	float y = -(((float)(i / 4096)) / (4096.0f * 0.5f)) + 1.0f;
`)

	scanner := bufio.NewScanner(f)
	last := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.Split(line, " ")
		switch tokens[1] {
		case "var-x":
			g.WriteString(fmt.Sprintf("\tfloat %s = x;\n", tokens[0]))
		case "var-y":
			g.WriteString(fmt.Sprintf("\tfloat %s = y;\n", tokens[0]))
		case "const":
			g.WriteString(fmt.Sprintf("\tfloat %s = %s;\n", tokens[0], tokens[2]))
		case "add":
			g.WriteString(fmt.Sprintf("\tfloat %s = %s + %s;\n", tokens[0], tokens[2], tokens[3]))
		case "sub":
			g.WriteString(fmt.Sprintf("\tfloat %s = %s - %s;\n", tokens[0], tokens[2], tokens[3]))
		case "mul":
			g.WriteString(fmt.Sprintf("\tfloat %s = %s * %s;\n", tokens[0], tokens[2], tokens[3]))
		case "max":
			g.WriteString(fmt.Sprintf("\tfloat %s = fmaxf(%s, %s);\n", tokens[0], tokens[2], tokens[3]))
		case "min":
			g.WriteString(fmt.Sprintf("\tfloat %s = fminf(%s, %s);\n", tokens[0], tokens[2], tokens[3]))
		case "neg":
			g.WriteString(fmt.Sprintf("\tfloat %s = -%s;\n", tokens[0], tokens[2]))
		case "square":
			g.WriteString(fmt.Sprintf("\tfloat %s = %s * %s;\n", tokens[0], tokens[2], tokens[2]))
		case "sqrt":
			g.WriteString(fmt.Sprintf("\tfloat %s = sqrtf(%s);\n", tokens[0], tokens[2]))
		default:
			panic("unknown instruction")
		}
		last = tokens[0]
	}
	g.WriteString("\n")
	g.WriteString("\tzs[i] = " + last + " < 0 ? 0 : 255;\n")
	g.WriteString("}\n")
	g.WriteString(`
int main()
{
	int *z;
	int N = 4096 * 4096;

	cudaMallocManaged(&z, sizeof(int) * N);

	for (int i = 0; i < 1000; i++)
	{
		prospero<<<N / 512, 512>>>(z, N);
		cudaDeviceSynchronize();
	}

	printf("P2\n4096 4096\n255\n");
	for (int i = 0; i < N; i++)
		printf("%d ", z[i]);

	return 0;
}
`)
}
