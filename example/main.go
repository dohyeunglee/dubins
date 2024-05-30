package main

import (
	"fmt"
	"math"

	"github.com/dohyeunglee/dubins"
)

func main() {
	start := dubins.State{
		X:   0,
		Y:   0,
		Yaw: 0,
	}

	goal := dubins.State{
		X:   1,
		Y:   1,
		Yaw: math.Pi,
	}

	turningRadius := 1.0

	path, ok := dubins.MinLengthPath(start, goal, turningRadius)
	if !ok {
		fmt.Println("There is no available DubinsPath.")
		return
	}

	fmt.Printf("[MinLength DubinsPath from %s to %s]\n", start, goal)
	fmt.Println("- PathType: ", path.PathType())
	fmt.Println("- Length: ", path.Length())
	fmt.Println("- Segments: ", path.Segments())
	fmt.Println("- Interpolated: ", path.Interpolate(0.1))
}
