package bs

import (
	"math"
	"testing"
)

func BS(call bool, spot, strike, timeToExpiration, interestRate, dividend, iv float64,
	calculateGreeks bool) BSOptionPrice {
	o := BSOption{Call: call, Spot: spot, Strike: strike,
		TimeToExpiration: timeToExpiration, InterestRate: interestRate,
		Dividend: dividend, IV: iv}
	return o.Calculate(calculateGreeks)
}

func TestBS(t *testing.T) {
	shouldClose(t, BS(true, 100, 95, 0.00273972602739726, 0.000751040922831883,
		0, 0.2, false).Price, 4.99998980469552)
	shouldClose(t, BS(true, 92.45, 107.5, 0.0876712328767123, 0.00192960198828152,
		0, 0.3, false).Price, 0.162619795863781)
	shouldClose(t, BS(true, 93.0766666666667, 107.75, 0.164383561643836,
		0.00266390125346286, 0, 0.2878, false).Price, 0.584588840095316)
	shouldClose(t, BS(true, 93.5333333333333, 107.75, 0.249315068493151,
		0.00319934651984034, 0, 0.2907, false).Price, 1.27026849732877)
	shouldClose(t, BS(true, 93.8733333333333, 107.75, 0.331506849315069,
		0.00350934592318849, 0, 0.2929, false).Price, 1.97015685523537)
	shouldClose(t, BS(true, 94.1166666666667, 107.75, 0.416438356164384,
		0.00367360967852615, 0, 0.2919, false).Price, 2.61731599547608)
	shouldClose(t, BS(false, 94.2666666666667, 107.75, 0.498630136986301,
		0.00372609838856132, 0, 0.2888, false).Price, 16.6074587545269)
	shouldClose(t, BS(false, 94.3666666666667, 107.75, 0.583561643835616,
		0.00370681407974257, 0, 0.2923, false).Price, 17.1686196701434)
	shouldClose(t, BS(false, 94.44, 107.75, 0.668493150684932,
		0.00364163303865433, 0, 0.2908, false).Price, 17.6038273793172)
	shouldClose(t, BS(false, 94.4933333333333, 107.75, 0.750684931506849,
		0.00355604221290591, 0, 0.2919, false).Price, 18.0870982577296)
	shouldClose(t, BS(false, 94.49, 107.75, 0.835616438356164,
		0.00346100468320478, 0, 0.2901, false).Price, 18.5149895730975)
	shouldClose(t, BS(false, 94.39, 107.75, 0.917808219178082,
		0.00337464630758452, 0, 0.2876, false).Price, 18.9397688539483)
	shouldClose(t, BS(true, 100, 95, 1, 1, 0, 1, false).Price, 14.6711476484)
	shouldClose(t, BS(false, 100, 95, 1, 1, 0, 1, false).Price, 12.8317504425)
	shouldClose(t, BS(true, 100, 100, 0.00396825396825397, 0.000771332656950173,
		0, 0.15, false).Price, 0.376962465712609)
	shouldClose(t, BS(false, 100, 100, 0.00396825396825397, 0.000771332656950173,
		0, 0.15, false).Price, 0.376962465712609)
	shouldClose(t, BS(true, 100, 100, 100, 0.042033868311581,
		0, 0.15, false).Price, 0.817104022604705)
	shouldClose(t, BS(false, 100, 100, 100, 0.042033868311581,
		0, 0.15, false).Price, 0.817104022604705)
	shouldClose(t, BS(true, 100, 0.01, 1, 0.00330252458693489,
		0, 0.15, false).Price, 99.660325245681)
	shouldClose(t, BS(false, 100, 0.01, 1, 0.00330252458693489,
		0, 0.15, false).Price, 0)
	shouldClose(t, BS(true, 100, 2147483248, 1, 0.00330252458693489,
		0, 0.15, false).Price, 0)
	shouldClose(t, BS(false, 100, 2147483248, 1, 0.00330252458693489,
		0, 0.15, false).Price, 2140402730.16601)
	shouldClose(t, BS(true, 0.01, 100, 1, 0.00330252458693489,
		0, 0.15, false).Price, 0)
	shouldClose(t, BS(false, 0.01, 100, 1, 0.00330252458693489,
		0, 0.15, false).Price, 99.660325245681)
	shouldClose(t, BS(true, 2147483248, 100, 1, 0.00330252458693489,
		0, 0.15, false).Price, 2140402730.16601)
	shouldClose(t, BS(false, 2147483248, 100, 1, 0.00330252458693489,
		0, 0.15, false).Price, 0)
	shouldClose(t, BS(true, 100, 100, 1, 0.05, -1, 0.15, false).Price,
		1.62505648981223e-11)
	shouldClose(t, BS(false, 100, 100, 1, 0.05, -1, 0.15, false).Price,
		60.1291675389721)
	shouldClose(t, BS(true, 100, 100, 1, 0.05, 1, 0.15, false).Price,
		163.448023481557)
	shouldClose(t, BS(false, 100, 100, 1, 0.05, 1, 0.15, false).Price,
		4.4173615264761e-11)
	shouldClose(t, BS(true, 100, 100, 1, -1, 0, 0.15, false).Price,
		16.2513262267156)
	shouldClose(t, BS(false, 100, 100, 1, -1, 0, 0.15, false).Price,
		16.2513262267156)
	shouldClose(t, BS(true, 100, 100, 1, 1, 0, 0.15, false).Price,
		2.19937783786316)
	shouldClose(t, BS(false, 100, 100, 1, 1, 0, 0.15, false).Price,
		2.19937783786316)
	shouldClose(t, BS(true, 100, 100, 1, 0.05, 0, 0.005, false).Price,
		0.189742620249)
	shouldClose(t, BS(false, 100, 100, 1, 0.05, 0, 0.005, false).Price,
		0.189742620249)
	shouldClose(t, BS(true, 100, 100, 1, 0.05, 0, 1, false).Price,
		36.424945370234)
	shouldClose(t, BS(false, 100, 100, 1, 0.05, 0, 1, false).Price,
		36.424945370234)
	shouldClose(t, BS(true, 100, 100, 1, 0.05, 0, 0.15, false).Price,
		5.68695251984796)

	option := BSOption{true, 100, 100, 1, 0.05, 0, 0.15}
	op := option.Calculate(true)
	shouldClose(t, op.Delta, 0.50404947485)
	shouldClose(t, op.Gamma, 0.025227988795588)
	shouldClose(t, op.Theta, -2.55380111351125)
	shouldClose(t, op.Vega, 37.84198319338195)
	shouldClose(t, op.Rho, 44.7179949651117)

	option.Call = false
	op = option.Calculate(true)
	shouldClose(t, op.Price, 5.68695251984796)
	shouldClose(t, op.Delta, -0.447179949651)
	shouldClose(t, op.Gamma, 0.025227988795588)
	shouldClose(t, op.Theta, -2.55380111351125)
	shouldClose(t, op.Vega, 37.84198319338195)
	shouldClose(t, op.Rho, -50.4049474849597)

	option = BSOption{true, 100, 100, 2, 0.05, 0.05, 0.25}
	op = option.Calculate(true)
	shouldClose(t, op.Vega, 50.7636345571413)

	option.Call = false
	op = option.Calculate(true)
	shouldClose(t, op.Vega, 50.7636345571413)
}

func shouldClose(t *testing.T, x, expected float64) {
	diff := math.Abs(x - expected)
	if diff > .000005 || math.IsNaN(x) {
		t.Helper()
		t.Fatalf("%f is not same as %f (diff: %f)", x, expected, diff)
	}
}

func IV(call bool, spot, strike, timeToExpiration, interestRate, dividend,
	optionPrice float64) float64 {
	o := BSOption{Call: call, Spot: spot, Strike: strike,
		TimeToExpiration: timeToExpiration, InterestRate: interestRate,
		Dividend: dividend}
	return o.CalculateIV(optionPrice, 0)
}

func TestIV(t *testing.T) {
	shouldClose(t, IV(true, 92.45, 107.5, 0.0876712328767123,
		0.00192960198828152, 0, 0.162619795863781), 0.3)
	shouldClose(t, IV(true, 93.0766666666667, 107.75, 0.164383561643836,
		0.00266390125346286, 0, 0.584588840095316), 0.2878)
	shouldClose(t, IV(true, 93.5333333333333, 107.75, 0.249315068493151,
		0.00319934651984034, 0, 1.27026849732877), 0.2907)
	shouldClose(t, IV(true, 93.8733333333333, 107.75, 0.331506849315069,
		0.00350934592318849, 0, 1.97015685523537), 0.2929)
	shouldClose(t, IV(true, 94.1166666666667, 107.75, 0.416438356164384,
		0.00367360967852615, 0, 2.61731599547608), 0.2919)
	shouldClose(t, IV(false, 94.2666666666667, 107.75, 0.498630136986301,
		0.00372609838856132, 0, 16.6074587545269), 0.2888)
	shouldClose(t, IV(false, 94.3666666666667, 107.75, 0.583561643835616,
		0.00370681407974257, 0, 17.1686196701434), 0.2923)
	shouldClose(t, IV(false, 94.44, 107.75, 0.668493150684932,
		0.00364163303865433, 0, 17.6038273793172), 0.2908)
	shouldClose(t, IV(false, 94.4933333333333, 107.75, 0.750684931506849,
		0.00355604221290591, 0, 18.0870982577296), 0.2919)
	shouldClose(t, IV(false, 94.39, 107.75, 0.917808219178082,
		0.00337464630758452, 0, 18.9397688539483), 0.2876)
	shouldClose(t, IV(true, 100, 95, 1, 1, 0, 14.6711476484), 1)
	shouldClose(t, IV(false, 100, 95, 1, 1, 0, 12.8317504425), 1)
}
