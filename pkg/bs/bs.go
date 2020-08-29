package bs

import (
	"math"
)

// BS Inputs:
//  ot = optionType
//  fs = price of underlying
//   x = strike
//   t = time to expiration
//   r = risk free rate
//   b = cost of carry
//   v = implied volatility
func BS(ot OptionType, fs, x, t, r, b, v float64) *BSModel {
	tSqrt := math.Sqrt(t)
	d1 := (math.Log(fs/x) + (b+(v*v)/2)*t) / (v * tSqrt)
	d2 := d1 - v*tSqrt

	var (
		value, delta, gamma, theta, vega, rho float64
	)
	if ot == CALL {
		value = fs*math.Exp((b-r)*t)*norm.cdf(d1) - x*math.Exp(-r*t)*norm.cdf(d2)
		delta = math.Exp((b-r)*t) * norm.cdf(d1)
		gamma = math.Exp((b-r)*t) * norm.pdf(d1) / (fs * v * tSqrt)
		theta = -(fs*v*math.Exp((b-r)*t)*norm.pdf(d1))/(2*tSqrt) - (b-r)*fs*math.Exp((b-r)*t)*norm.cdf(d1) - r*x*math.Exp(-r*t)*norm.cdf(d2)
		vega = math.Exp((b-r)*t) * fs * tSqrt * norm.pdf(d1)
		rho = x * t * math.Exp(-r*t) * norm.cdf(d2)
	}
	if ot == PUT {
		value = x*math.Exp(-r*t)*norm.cdf(-d2) - (fs * math.Exp((b-r)*t) * norm.cdf(-d1))
		delta = -math.Exp((b-r)*t) * norm.cdf(-d1)
		gamma = math.Exp((b-r)*t) * norm.pdf(d1) / (fs * v * tSqrt)
		theta = -(fs*v*math.Exp((b-r)*t)*norm.pdf(d1))/(2*tSqrt) + (b-r)*fs*math.Exp((b-r)*t)*norm.cdf(-d1) + r*x*math.Exp(-r*t)*norm.cdf(-d2)
		vega = math.Exp((b-r)*t) * fs * tSqrt * norm.pdf(d1)
		rho = -x * t * math.Exp(-r*t) * norm.cdf(-d2)
	}
	return &BSModel{
		value, delta, gamma, theta, vega, rho,
	}
}

// IV Inputs:
//   ot = optionType
//   fs = price of underlying
//    x = strike
//    t = time to expiration
//    r = risk free rate
//    b = cost of carry
//   cp = Call or Put price
func IV(ot OptionType, fs, x, t, r, b, cp float64) float64 {
	v := approxImpliedVol(ot, fs, x, t, r, b, cp)
	v = math.Max(v, minIV)
	v = math.Min(v, maxIV)
	bs := BS(ot, fs, x, t, r, b, v)
	minDiff := math.Abs(cp - bs.Price)
	countr := 0
	for 0.00001 <= math.Abs(cp-bs.Price) && math.Abs(cp-bs.Price) <= minDiff && countr < maxIter {
		v = v - (bs.Price-cp)/bs.Vega
		if (v > maxIV) || (v < minIV) {
			break
		}
		bs = BS(ot, fs, x, t, r, b, v)
		minDiff = math.Min(math.Abs(cp-bs.Price), minDiff)
		countr++
	}
	if math.Abs(cp-bs.Price) < 0.00001 {
		return v
	}
	return bisectionImpliedVol(ot, fs, x, t, r, b, cp, 0.00001, maxIter)
}

func bisectionImpliedVol(ot OptionType, fs, x, t, r, b, cp, precision float64, maxSteps int) float64 {
	vMid := approxImpliedVol(ot, fs, x, t, r, b, cp)
	var (
		vLow, vHigh float64
	)
	if (vMid <= minIV) || (vMid >= maxIV) {
		vLow = minIV
		vHigh = maxIV
		vMid = (vLow + vHigh) / 2
	} else {
		vLow = math.Max(minIV, vMid*.5)
		vHigh = math.Min(maxIV, vMid*1.5)
	}

	cpMid := BS(ot, fs, x, t, r, b, vMid).Price

	currentStep := 0
	diff := math.Abs(cp - cpMid)

	for (diff > precision) && (currentStep < maxSteps) {
		currentStep++

		if cpMid < cp {
			vLow = vMid
		} else {
			vHigh = vMid
		}

		cpLow := BS(ot, fs, x, t, r, b, vLow).Price
		cpHigh := BS(ot, fs, x, t, r, b, vHigh).Price

		vMid = vLow + (cp-cpLow)*(vHigh-vLow)/(cpHigh-cpLow)
		vMid = math.Max(minIV, vMid)
		vMid = math.Min(maxIV, vMid)

		cpMid = BS(ot, fs, x, t, r, b, vMid).Price
		diff = math.Abs(cp - cpMid)
	}

	if math.Abs(cp-cpMid) < precision {
		return vMid
	}
	return math.NaN()
}

func approxImpliedVol(ot OptionType, fs, x, t, r, b, cp float64) float64 {

	ebrt := math.Exp((b - r) * t)
	ert := math.Exp(-r * t)

	a := math.Sqrt(2*math.Pi) / (fs*ebrt + x*ert)

	var payoff float64
	if ot == CALL {
		payoff = fs*ebrt - x*ert
	}
	if ot == PUT {
		payoff = x*ert - fs*ebrt
	}

	b = cp - payoff/2
	c := math.Pow(payoff, 2) / math.Pi

	v := (a * (b + math.Sqrt(math.Pow(b, 2)+c))) / math.Sqrt(t)

	return v
}

type OptionType int

const (
	PUT OptionType = iota
	CALL
)

type BSModel struct {
	Price, Delta, Gamma, Theta, Vega, Rho float64
}

const (
	maxIV   = 2.0
	minIV   = 0.01
	maxIter = 100
)

type normalDist struct {
	Mu, Sigma float64
}

var norm = normalDist{0, 1}

const invSqrt2Pi = 0.39894228040143267793994605993438186847585863116493465766592583

func (n normalDist) pdf(x float64) float64 {
	z := x - n.Mu
	return math.Exp(-z*z/(2*n.Sigma*n.Sigma)) * invSqrt2Pi / n.Sigma
}

func (n normalDist) cdf(x float64) float64 {
	return math.Erfc(-(x-n.Mu)/(n.Sigma*math.Sqrt2)) / 2
}
