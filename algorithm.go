package dubins

import (
	"math"
)

const (
	twoPi = 2 * math.Pi
)

type pathTypeFunc func(pathParams) (Path, bool)

var (
	pathTypeFuncMap = map[PathType]pathTypeFunc{
		LSL: lsl,
		LSR: lsr,
		RSL: rsl,
		RSR: rsr,
		RLR: rlr,
		LRL: lrl,
	}
	pathTypeCourseTypesMap = map[PathType][3]PathCourseType{
		LSL: {CourseTypeLeft, CourseTypeStraight, CourseTypeLeft},
		LSR: {CourseTypeLeft, CourseTypeStraight, CourseTypeRight},
		RSL: {CourseTypeRight, CourseTypeStraight, CourseTypeLeft},
		RSR: {CourseTypeRight, CourseTypeStraight, CourseTypeRight},
		RLR: {CourseTypeRight, CourseTypeLeft, CourseTypeRight},
		LRL: {CourseTypeLeft, CourseTypeRight, CourseTypeLeft},
	}
)

type pathParams struct {
	alpha float64
	beta  float64
	d     float64

	sinAlpha          float64
	sinBeta           float64
	cosAlpha          float64
	cosBeta           float64
	cosAlphaMinusBeta float64
	dSquare           float64
}

func lsl(pp pathParams) (Path, bool) {
	pSquare := 2 + pp.dSquare - 2*pp.cosAlphaMinusBeta + 2*pp.d*(pp.sinAlpha-pp.sinBeta)
	if pSquare < 0 {
		return Path{}, false
	}

	tmp := math.Atan2(pp.cosBeta-pp.cosAlpha, pp.d+pp.sinAlpha-pp.sinBeta)
	return Path{
		t:        mod2pi(tmp - pp.alpha),
		p:        math.Sqrt(pSquare),
		q:        mod2pi(pp.beta - tmp),
		pathType: LSL,
	}, true
}

func rsr(pp pathParams) (Path, bool) {
	pSquare := 2 + pp.dSquare - 2*pp.cosAlphaMinusBeta + 2*pp.d*(pp.sinBeta-pp.sinAlpha)
	if pSquare < 0 {
		return Path{}, false
	}

	tmp := math.Atan2(pp.cosAlpha-pp.cosBeta, pp.d-pp.sinAlpha+pp.sinBeta)
	return Path{
		t:        mod2pi(pp.alpha - tmp),
		p:        math.Sqrt(pSquare),
		q:        mod2pi(tmp - pp.beta),
		pathType: RSR,
	}, true
}

func lsr(pp pathParams) (Path, bool) {
	pSquare := -2 + pp.dSquare + 2*pp.cosAlphaMinusBeta + 2*pp.d*(pp.sinAlpha+pp.sinBeta)
	if pSquare < 0 {
		return Path{}, false
	}

	p := math.Sqrt(pSquare)
	tmp := math.Atan2(-pp.cosAlpha-pp.cosBeta, pp.d+pp.sinAlpha+pp.sinBeta) - math.Atan2(-2, p)
	t := mod2pi(tmp - pp.alpha)
	q := mod2pi(tmp - mod2pi(pp.beta))
	return Path{
		t:        t,
		p:        p,
		q:        q,
		pathType: LSR,
	}, true
}

func rsl(pp pathParams) (Path, bool) {
	pSquare := -2 + pp.dSquare + 2*pp.cosAlphaMinusBeta - 2*pp.d*(pp.sinAlpha+pp.sinBeta)
	if pSquare < 0 {
		return Path{}, false
	}
	p := math.Sqrt(pSquare)
	tmp := math.Atan2(pp.cosAlpha+pp.cosBeta, pp.d-pp.sinAlpha-pp.sinBeta) - math.Atan2(2, p)
	t := mod2pi(pp.alpha - tmp)
	q := mod2pi(pp.beta - tmp)
	return Path{
		t:        t,
		p:        p,
		q:        q,
		pathType: RSL,
	}, true
}

func rlr(pp pathParams) (Path, bool) {
	tmp := (6 - pp.dSquare + 2*pp.cosAlphaMinusBeta + 2*pp.d*(pp.sinAlpha-pp.sinBeta)) / 8
	if math.Abs(tmp) > 1 {
		return Path{}, false
	}
	phi := math.Atan2(pp.cosAlpha-pp.cosBeta, pp.d-pp.sinAlpha+pp.sinBeta)
	p := mod2pi(twoPi - math.Acos(tmp))
	t := mod2pi(pp.alpha - phi + mod2pi(p/2))
	q := mod2pi(pp.alpha - pp.beta - t + mod2pi(p))
	return Path{
		t:        t,
		p:        p,
		q:        q,
		pathType: RLR,
	}, true
}

func lrl(pp pathParams) (Path, bool) {
	tmp := (6 - pp.dSquare + 2*pp.cosAlphaMinusBeta + 2*pp.d*(pp.sinBeta-pp.sinAlpha)) / 8
	if math.Abs(tmp) > 1 {
		return Path{}, false
	}
	phi := math.Atan2(pp.cosAlpha-pp.cosBeta, pp.d+pp.sinAlpha-pp.sinBeta)
	p := mod2pi(twoPi - math.Acos(tmp))
	t := mod2pi(-pp.alpha - phi + p/2)
	q := mod2pi(mod2pi(pp.beta) - pp.alpha - t + mod2pi(p))
	return Path{
		t:        t,
		p:        p,
		q:        q,
		pathType: LRL,
	}, true
}

func fmodr(x float64, y float64) float64 {
	return x - y*math.Floor(x/y)
}

func mod2pi(theta float64) float64 {
	return fmodr(theta, twoPi)
}
