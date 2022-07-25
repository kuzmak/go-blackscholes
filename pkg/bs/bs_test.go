package bs

import (
	"testing"
)

func TestBS(t *testing.T) {
	/*
	   Values verified using Excel spreadsheet provided available at http://www-2.rotman.utoronto.ca/~hull/software/index.html
	   DerivaGem 4.00 (http://www-2.rotman.utoronto.ca/~hull/software/DG400.zip)
	   Book: John C. Hull OPTIONS, FUTURES, AND OTHER DERIVATIVES 8th Edition
	*/

	t.Run("check call option price/iv", func(t *testing.T) {
		option := BSOption{Call: true, Spot: 45.0, Strike: 45.0,
			TimeToExpiration: 1, InterestRate: 0.02,
			Dividend: 0, IV: 0.25}
		op := option.Calculate(true)

		e := 4.891751320750913
		if op.Price != e {
			t.Errorf("got %0.2f, want %0.2f", op.Price, e)
		}

		e = 0.5812139374874482
		if op.Delta != e {
			t.Errorf("got %0.2f, want %0.2f", op.Delta, e)
		}

		e = -0.007185333390922677
		if op.Theta != e {
			t.Errorf("got %0.2f, want %0.2f", op.Theta, e)
		}

		e = 0.21262875866184255
		if op.Rho != e {
			t.Errorf("got %0.2f, want %0.2f", op.Rho, e)
		}

		e = 0.034724174544009355
		if op.Gamma != e {
			t.Errorf("got %0.2f, want %0.2f", op.Gamma, e)
		}

		e = 0.17579113362904736
		if op.Vega != e {
			t.Errorf("got %0.2f, want %0.2f", op.Vega, e)
		}

		e = 0.24990038376703688
		iv := option.CalculateIV(4.89, 0)
		if iv != e {
			t.Errorf("got %0.2f, want %0.2f", iv, e)
		}
	})

	t.Run("check put option price/iv", func(t *testing.T) {
		option := BSOption{Call: false, Spot: 45.0, Strike: 45.0,
			TimeToExpiration: 1, InterestRate: 0.02,
			Dividend: 0, IV: 0.25}
		op := option.Calculate(true)

		e := 4.000691619554896
		if op.Price != e {
			t.Errorf("got %0.2f, want %0.2f", op.Price, e)
		}

		e = -0.41878606251255185
		if op.Delta != e {
			t.Errorf("got %0.2f, want %0.2f", op.Delta, e)
		}

		e = -0.004768405155371774
		if op.Theta != e {
			t.Errorf("got %0.2f, want %0.2f", op.Theta, e)
		}

		e = -0.2284606443261973
		if op.Rho != e {
			t.Errorf("got %0.2f, want %0.2f", op.Rho, e)
		}

		e = 0.034724174544009355
		if op.Gamma != e {
			t.Errorf("got %0.2f, want %0.2f", op.Gamma, e)
		}

		e = 0.17579113362904736
		if op.Vega != e {
			t.Errorf("got %0.2f, want %0.2f", op.Vega, e)
		}

		e = 0.2499606747801532
		iv := option.CalculateIV(4.0, 0)
		if iv != e {
			t.Errorf("got %0.2f, want %0.2f", iv, e)
		}
	})
}
