package dubins

import (
	"fmt"
	"math"
)

// PathCourseType represents DubinsPath segment course type: L, S, R.
type PathCourseType string

const (
	CourseTypeLeft     PathCourseType = "L"
	CourseTypeStraight PathCourseType = "S"
	CourseTypeRight    PathCourseType = "R"
)

// PathType represents 6 DubinsPath type.
type PathType string

const (
	LSL PathType = "LSL"
	LSR PathType = "LSR"
	RSL PathType = "RSL"
	RSR PathType = "RSR"
	RLR PathType = "RLR"
	LRL PathType = "LRL"
)

type PathSegment struct {
	Length     float64
	CourseType PathCourseType
}

// Path corresponds to DubinsPath.
type Path struct {
	t        float64
	p        float64
	q        float64
	pathType PathType

	totalLength   float64
	segments      []PathSegment
	turningRadius float64
	origin        State
}

// Length returns the total length of DubinsPath.
func (p *Path) Length() float64 {
	return p.totalLength
}

// PathType returns the type of DubinsPath.
func (p *Path) PathType() PathType {
	return p.pathType
}

// Segments returns all `PathSegment` information about DubinsPath.
func (p *Path) Segments() []PathSegment {
	return p.segments
}

// Interpolate interpolates DubinsPath by `stepSize`, returning `State` list.
func (p *Path) Interpolate(stepSize float64) []State {
	states := make([]State, 0)

	for distance := 0.0; distance <= p.Length(); distance += stepSize {
		states = append(states, p.sample(distance))
	}

	return states
}

// sample returns the state which is `distance` far from `origin` state.
func (p *Path) sample(distance float64) State {
	if distance <= 0 {
		return p.origin
	}

	distance = min(distance, p.Length())

	start := p.origin

	for _, segment := range p.segments {
		deltaDistance := min(distance, segment.Length)
		distance -= deltaDistance

		start = stateAtDistance(start, deltaDistance, segment.CourseType, p.turningRadius)

		if distance <= 0 {
			break
		}
	}

	return start
}

func (p *Path) setTuringRadius(turingRadius float64) {
	p.turningRadius = turingRadius
	p.t *= p.turningRadius
	p.p *= p.turningRadius
	p.q *= p.turningRadius

	lengths := [3]float64{p.t, p.p, p.q}
	courseTypes := pathTypeCourseTypesMap[p.pathType]
	p.segments = make([]PathSegment, 0, len(lengths))

	for i, length := range lengths {
		p.totalLength += length
		p.segments = append(p.segments, PathSegment{Length: length, CourseType: courseTypes[i]})
	}
}

// State represents the 3D state, x, y coordinate and yaw angle.
type State struct {
	X   float64
	Y   float64
	Yaw float64
}

func (s State) String() string {
	return fmt.Sprintf("(%v, %v, %v)", s.X, s.Y, s.Yaw)
}

// PathByType returns DubinsPath that can reach the `goal` state from the `start` state, given a `turningRadius` and `pathType`
func PathByType(start State, goal State, turingRadius float64, pathType PathType) (Path, bool) {
	pp := toPathParams(start, goal, turingRadius)
	if pp.alpha == 0 && pp.beta == 0 && pp.d == 0 {
		return Path{}, false
	}

	f, ok := pathTypeFuncMap[pathType]
	if !ok {
		panic("unknown path type")
	}

	path, ok := f(pp)
	if !ok {
		return Path{}, false
	}

	path.setTuringRadius(turingRadius)
	if path.Length() <= 0 {
		return Path{}, false
	}

	path.origin = start
	return path, true
}

// AvailablePaths returns all possible DubinsPath that can reach the `goal` state from the `start` state, given a `turningRadius`.
func AvailablePaths(start State, goal State, turingRadius float64) []Path {
	pp := toPathParams(start, goal, turingRadius)
	if pp.alpha == 0 && pp.beta == 0 && pp.d == 0 {
		return []Path{}
	}

	paths := make([]Path, 0, len(pathTypeFuncMap))

	for _, f := range pathTypeFuncMap {
		path, ok := f(pp)
		if !ok {
			continue
		}

		path.setTuringRadius(turingRadius)
		if path.Length() <= 0 {
			continue
		}

		path.origin = start
		paths = append(paths, path)
	}

	return paths
}

// MinLengthPath returns the shortest DubinsPath among all possible paths from the `start` state to the `goal` state, given a `turningRadius`.
// In case there are multiple paths with the same length, one of them is randomly selected and returned.
// The shortest DubinsPath length can be obtained using the `Length` function of the returned DubinsPath.
func MinLengthPath(start State, goal State, turningRadius float64) (Path, bool) {
	paths := AvailablePaths(start, goal, turningRadius)
	if len(paths) == 0 {
		return Path{}, false
	}

	minLength := math.MaxFloat64
	var bestPath Path

	for _, path := range paths {
		length := path.Length()
		if length < minLength {
			minLength = length
			bestPath = path
		}
	}

	return bestPath, true
}

func toPathParams(start State, goal State, turningRadius float64) pathParams {
	if turningRadius <= 0 {
		panic("turningRadius must be greater than zero")
	}

	dx := goal.X - start.X
	dy := goal.Y - start.Y
	D := math.Sqrt(dx*dx + dy*dy)
	theta := 0.0
	if D > 0 {
		theta = mod2pi(math.Atan2(dy, dx))
	}

	alpha := mod2pi(start.Yaw - theta)
	beta := mod2pi(goal.Yaw - theta)
	d := D / turningRadius

	return pathParams{
		alpha:             alpha,
		beta:              beta,
		d:                 d,
		sinAlpha:          math.Sin(alpha),
		sinBeta:           math.Sin(beta),
		cosAlpha:          math.Cos(alpha),
		cosBeta:           math.Cos(beta),
		cosAlphaMinusBeta: math.Cos(alpha - beta),
		dSquare:           d * d,
	}
}

func stateAtDistance(start State, deltaDistance float64, courseType PathCourseType, turningRadius float64) State {
	if deltaDistance == 0 {
		return start
	}

	deltaX := 0.0
	deltaY := 0.0
	deltaYaw := 0.0

	phi := deltaDistance / turningRadius

	switch courseType {
	case CourseTypeLeft:
		deltaX = turningRadius*math.Sin(start.Yaw+phi) - turningRadius*math.Sin(start.Yaw)
		deltaY = turningRadius*-math.Cos(start.Yaw+phi) + turningRadius*math.Cos(start.Yaw)
		deltaYaw = phi
	case CourseTypeRight:
		deltaX = turningRadius*-math.Sin(start.Yaw-phi) + turningRadius*math.Sin(start.Yaw)
		deltaY = turningRadius*math.Cos(start.Yaw-phi) - turningRadius*math.Cos(start.Yaw)
		deltaYaw = -phi
	case CourseTypeStraight:
		deltaX = deltaDistance * math.Cos(start.Yaw)
		deltaY = deltaDistance * math.Sin(start.Yaw)
	default:
		panic("unknown course type")
	}

	return State{
		X:   start.X + deltaX,
		Y:   start.Y + deltaY,
		Yaw: start.Yaw + deltaYaw,
	}
}
