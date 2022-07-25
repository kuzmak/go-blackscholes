kpackage bs

import (
	"math"
)

const (
	maxIV                = 2.0
	minIV                = 0.01
	defaultMaxIterations = 100
	defaultPrecision     = 0.00001
	invSqrt2Pi           = 0.39894228040143267793994605993438186847585863116493465766592583
)

type normalDist struct {
	mu    float64
	sigma float64
}

var norm = normalDist{0.0, 1.0}
var sqrt2Pi float64 = math.Sqrt(math.Pi * 2)

type BSOption struct {
	Call             bool
	Spot             float64
	Strike           float64
	TimeToExpiration float64
	InterestRate     float64
	Dividend         float64
	IV               float64
}

type BSOptionPrice struct {
	Price float64
	Delta float64
	Gamma float64
	Theta float64
	Vega  float64
	Rho   float64
}

func (o BSOption) calculate(calculateVega bool, calculateGreeks bool) BSOptionPrice {
	sign := -1.0
	if o.Call {
		sign = 1.0
	}

	op := BSOptionPrice{}
	if o.TimeToExpiration <= 0 {
		op.Price = math.Max(0, math.Abs(o.Strike-o.Spot))
		return op
	}

	tSqrt := math.Sqrt(o.TimeToExpiration)
	exp1 := math.Exp(-o.Dividend * o.TimeToExpiration)
	exp2 := math.Exp(-o.InterestRate * o.TimeToExpiration)
	vt := o.IV * math.Sqrt(o.TimeToExpiration)

	d1 := math.Log(o.Spot/o.Strike) +
		(o.TimeToExpiration * (o.InterestRate - o.Dividend + ((o.IV * o.IV) * 0.5)))
	d1 = d1 / vt
	pdfD1 := norm.pdf(d1)

	d2 := sign * (d1 - o.IV*tSqrt)
	d1 = sign * d1

	nd1, nd2 := norm.cdf(d1), norm.cdf(d2)

	op.Price = sign * ((o.Spot * exp1 * nd1) - (o.Strike * exp2 * nd2))

	if calculateVega || calculateGreeks {
		op.Vega = exp1 * o.Spot * tSqrt * pdfD1 * 0.01
	}

	if calculateGreeks {
		op.Delta = sign * exp1 * nd1
		op.Gamma = exp1 * pdfD1 / (o.Spot * o.IV * tSqrt)
		op.Rho = sign * o.Strike * o.TimeToExpiration * exp2 * nd2 * 0.01

		d1pdf := math.Exp(-(d1*d1)*0.5) / sqrt2Pi
		p1 := ((o.Spot * o.IV * exp1) / (2 * tSqrt)) * d1pdf
		p2 := sign * o.InterestRate * o.Strike * exp2 * nd2
		p3 := sign * o.Dividend * o.Spot * exp1 * nd1
		op.Theta = (p3 - p1 - p2) / 365
	}

	return op
}

func (o BSOption) calculateWithVega() BSOptionPrice {
	return o.calculate(true, false)
}

func (o BSOption) calculateWithIV(iv float64) BSOptionPrice {
	o.IV = iv
	return o.calculate(false, false)
}

func (o BSOption) Calculate(calculateGreeks bool) BSOptionPrice {
	return o.calculate(false, true)
}

func (o BSOption) getApproximateIV(optionPrice float64) float64 {
	ebrt := math.Exp((o.Dividend - o.InterestRate) * o.TimeToExpiration)
	ert := math.Exp(-o.InterestRate * o.TimeToExpiration)

	a := math.Sqrt(2*math.Pi) / (o.Spot*ebrt + o.Strike*ert)

	var payoff float64
	if o.Call {
		payoff = o.Spot*ebrt - o.Strike*ert
	} else {
		payoff = o.Strike*ert - o.Spot*ebrt
	}

	d := optionPrice - payoff/2 // dividend
	c := math.Pow(payoff, 2) / math.Pi

	return (a * (d + math.Sqrt(math.Pow(d, 2)+c))) / math.Sqrt(o.TimeToExpiration)
}

func (o BSOption) getBisectionIV(optionPrice float64, maxSteps int) float64 {
	middle := o.getApproximateIV(optionPrice)
	var low, high float64
	if (middle <= minIV) || (middle >= maxIV) {
		low = minIV
		high = maxIV
		middle = (low + high) / 2
	} else {
		low = math.Max(minIV, middle*.5)
		high = math.Min(maxIV, middle*1.5)
	}

	cpMid := o.calculateWithIV(middle).Price

	currentStep := 0
	diff := math.Abs(optionPrice - cpMid)

	for (diff > defaultPrecision) && (currentStep < maxSteps) {
		currentStep++

		if cpMid < optionPrice {
			low = middle
		} else {
			high = middle
		}

		cpLow := o.calculateWithIV(low).Price
		cpHigh := o.calculateWithIV(high).Price

		middle = low + (optionPrice-cpLow)*(high-low)/(cpHigh-cpLow)
		middle = math.Max(minIV, middle)
		middle = math.Min(maxIV, middle)

		cpMid = o.calculateWithIV(middle).Price
		diff = math.Abs(optionPrice - cpMid)
	}

	if math.Abs(optionPrice-cpMid) < defaultPrecision {
		return middle
	}
	return math.NaN()
}

func (o BSOption) CalculateIV(optionPrice float64, maxIterations int) float64 {
	v := o.getApproximateIV(optionPrice)
	v = math.Max(v, minIV)
	v = math.Min(v, maxIV)

	if maxIterations <= 0 {
		maxIterations = defaultMaxIterations
	}

	bs := o.calculateWithVega()
	minDiff := math.Abs(optionPrice - bs.Price)
	count := 0
	for defaultPrecision <= math.Abs(optionPrice-bs.Price) &&
		math.Abs(optionPrice-bs.Price) <= minDiff && count < maxIterations {
		v = v - (bs.Price-optionPrice)/bs.Vega
		if (v > maxIV) || (v < minIV) {
			break
		}

		op := o.calculateWithIV(v)
		minDiff = math.Min(math.Abs(optionPrice-op.Price), minDiff)
		count++
	}
	if math.Abs(optionPrice-bs.Price) < defaultPrecision {
		return v
	}
	return o.getBisectionIV(optionPrice, maxIterations)
}

func (n normalDist) pdf(x float64) float64 {
	z := x - n.mu
	return math.Exp(-z*z/(2*n.sigma*n.sigma)) * invSqrt2Pi / n.sigma
}

func (n normalDist) cdf(x float64) float64 {
	return math.Erfc(-(x-n.mu)/(n.sigma*math.Sqrt2)) / 2
}

