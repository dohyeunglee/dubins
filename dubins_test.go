package dubins_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dohyeunglee/dubins"
)

// TestMinLengthPath compares the result of MinLengthPath and `dubins_shortest_path` function of Dubins-Curves.
func TestMinLengthPath(t *testing.T) {
	turningRadius := 1.0
	type args struct {
		start dubins.State
		goal  dubins.State
	}
	tests := []struct {
		name       string
		args       args
		want       float64
		wantNoPath bool
	}{
		{
			name: "if start and goal state are the same",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
			},
			wantNoPath: true,
		},
		{
			name: "(0, 0, 0) -> (4, 4, 3.142)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   4,
					Y:   4,
					Yaw: 3.142,
				},
			},
			want: 7.61450033750528,
		},
		{
			name: "(0, 0, 0) -> (-4, 0, 0)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   -4,
					Y:   0,
					Yaw: 0,
				},
			},
			want: 10.283185307179586,
			// 3.141592653589793 4.0 3.141592653589793
			// LSL
		},
		{
			name: "(0, 0, 0) -> (4, 4, 0)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   4,
					Y:   4,
					Yaw: 0,
				},
			},
			want: 5.854590436003225,
			// 0.9272952180016123 4.0 0.9272952180016123
			// LSR
		},
		{
			name: "(0, 0, 0) -> (4, -4, 0)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   4,
					Y:   -4,
					Yaw: 0,
				},
			},
			want: 5.8545904360032255,
			// 0.9272952180016121 4.000000000000001 0.9272952180016121
			// RSL
		},
		{
			name: "(0, 0, 0) -> (-4, 4, 0)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   -4,
					Y:   4,
					Yaw: 0,
				},
			},
			want: 10.283185307179586,
			// 3.141592653589793 4.000000000000001 3.141592653589793
			// LSR
		},
		{
			name: "(0, 0, 0) -> (-4, -4, 0)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   -4,
					Y:   -4,
					Yaw: 0,
				},
			},
			want: 10.283185307179586,
			// 3.141592653589793 4.0 3.141592653589793
			// RSL
		},
		{
			name: "(0, 0, 0) -> (4, 4, pi)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   4,
					Y:   4,
					Yaw: math.Pi,
				},
			},
			want: 7.613728608589373,
			// 0.46364760900080615 4.47213595499958 2.677945044588987
			// LSL
		},
		{
			name: "(0, 0, 0) -> (4, -4, pi)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   4,
					Y:   -4,
					Yaw: math.Pi,
				},
			},
			want: 7.613728608589374,
			// 0.46364760900080604 4.4721359549995805 2.677945044588987
			// RSR
		},
		{
			name: "(0, 0, 0) -> (0.5, 0, pi)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: 0,
				},
				goal: dubins.State{
					X:   0.5,
					Y:   0,
					Yaw: math.Pi,
				},
			},
			want: 7.258935602260172,
			// 1.2743144002944586 5.200264127924982 0.7843570740407309
			// RLR
		},
		{
			name: "(0, 0, pi/4) -> (4, 4, pi/4)",
			args: args{
				start: dubins.State{
					X:   0,
					Y:   0,
					Yaw: math.Pi / 4,
				},
				goal: dubins.State{
					X:   4,
					Y:   4,
					Yaw: math.Pi / 4,
				},
			},
			want: 5.656854249492381,
			// 0.0 5.656854249492381 0.0
			// LSL
		},
		{
			name: "(4, 4, pi/4) -> (0, 0, pi/4)",
			args: args{
				start: dubins.State{
					X:   4,
					Y:   4,
					Yaw: math.Pi / 4,
				},
				goal: dubins.State{
					X:   0,
					Y:   0,
					Yaw: math.Pi / 4,
				},
			},
			want: 11.940039556671966,
			// 3.141592653589793 5.656854249492381 3.141592653589793
			// LSL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			path, ok := dubins.MinLengthPath(tt.args.start, tt.args.goal, turningRadius)

			// THEN
			assert.Equal(t, !tt.wantNoPath, ok)
			if !tt.wantNoPath {
				assert.InDelta(t, tt.want, path.Length(), 1e-9)
			}
		})
	}
}

func BenchmarkMinLengthPath(b *testing.B) {
	turningRadius := 1.0

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		start := dubins.State{
			X:   rand.Float64() * 100,
			Y:   rand.Float64() * 100,
			Yaw: rand.Float64() * 2 * math.Pi,
		}
		goal := dubins.State{
			X:   rand.Float64() * 100,
			Y:   rand.Float64() * 100,
			Yaw: rand.Float64() * 2 * math.Pi,
		}
		b.StartTimer()
		path, _ := dubins.MinLengthPath(start, goal, turningRadius)
		_ = path.Length()
	}
}
