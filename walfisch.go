// Implements a simple phase delay from different n-Antenna elements
package pathloss

import (
	"fmt"
	"log"
	"math"

	"github.com/wiless/cellular/deployment"
	"github.com/wiless/vlib"
)

type WalfischIke struct {
	wsettings ModelSetting
}

func (w *WalfischIke) Set(ModelSetting) {

}
func (w WalfischIke) Get() ModelSetting {
	return ModelSetting{}
}
func (w WalfischIke) LossInDbNodes(txnode, rxnode deployment.Node, freqGHz float64) (plDb float64, valid bool) {

	return 0, true
}
func (w WalfischIke) LossInDb3D(src, dest vlib.Location3D, freqGHz float64) (plDb float64, valid bool) {
	//fmt.Printf("Hello World")
	FreqMHz := freqGHz * 1.0e3
	lamda := 3.0e8 / (FreqMHz * 1.0e6)
	distance := src.DistanceFrom(dest) / 1.0e3
	var result float64
	var theta float64
	var Lbf, Lmsd, Lrts, Lori, Lbsh float64
	var ka, kf, kd float64
	hb := src.Z
	hm := dest.Z
	theta = 90.0
	hroof := 30.0
	w1 := 15.0
	b := 30.0
	l := 1000.0
	if freqGHz < 100 && distance < 1 {
		Lbf = 32.4 + 20*math.Log10(distance) + 20*math.Log10(FreqMHz)
		if theta < 35 {
			Lori = -10 + 0.354*theta
		} else if theta < 55 {
			Lori = 2.5 + 0.075*(theta-35)
		} else {
			Lori = 4.0 - 0.114*(theta-55)
		}
		Lrts = -8.2 - 10*math.Log10(w1) + 10*math.Log10(FreqMHz) + 20*math.Log10(hroof-hm) + Lori
		diff_in_TxRxh := hb - hm
		d := distance / (diff_in_TxRxh)
		ds := lamda * math.Pow(d, 2)
		//fmt.Printf("The value of Ds:%f", ds)
		if l > ds {
			if hb > hroof {
				Lbsh = -18 * math.Log10(1+diff_in_TxRxh)
				if FreqMHz >= 2000 {
					ka = 71.4
					kf = -8
				} else {
					ka = 54
				}
				kd = 18
				kf = -4.0 + 1.5*(FreqMHz/925-1)
			} else {
				Lbsh = 0
				kd = 18 - 15*(diff_in_TxRxh/hroof)
				if distance > 0.5 {
					ka = 54 - 0.8*diff_in_TxRxh
				} else {
					ka = 54 - 1.6*diff_in_TxRxh
				}
			}
			Lmsd = Lbsh + ka + kd*math.Log10(distance) + kf*math.Log10(FreqMHz) - 9*math.Log10(b)
		} else {
			Q := math.Atan((hb - hm) / b)
			rho1 := math.Pow(hb-hm, 2)
			rho2 := math.Pow(b, 2)
			rho := math.Sqrt(rho1 + rho2)
			var Qm float64
			if hb > hroof {
				f := (hb - hm) / d
				g := math.Sqrt(b / lamda)
				Qm = 2.35 * (math.Pow(f*g, 0.9))
			} else if hb == hroof {
				Qm = b / d
			} else {
				f1 := b / (2 * math.Pi * d)
				f2 := math.Sqrt(lamda / rho)
				f3 := 1/Q - 1/(2*math.Pi+Q)
				Qm = f1 * f2 * f3
			}
			Lmsd = -20 * math.Log(Qm)

		}
		if Lmsd+Lrts > 0 {
			result = Lbf + Lmsd + Lrts
		} else {
			result = Lbf
		}
		fmt.Printf("The pathloss value in Db is %f\n", result)

	} else if freqGHz < 2 && distance > 1 {

		if FreqMHz >= 150 && FreqMHz < 1500 && distance > 0.05 {
			var Ch float64
			// Ch = 0.8 + (1.1*math.Log10(FreqMHz)-0.7)*dest.Z - 1.56*math.Log10(FreqMHz)
			if FreqMHz >= 150.0 && FreqMHz <= 200.0 {
				Ch = 8.29*math.Pow(math.Log10(1.54*dest.Z), 2) - 1.1
			} else if FreqMHz > 200.0 && FreqMHz <= 1500.0 {
				Ch = 3.2*math.Pow(math.Log10(11.75*dest.Z), 2) - 4.97
			}
			result = 69.55 + 26.16*math.Log10(FreqMHz) - 13.82*math.Log10(src.Z) - Ch + (44.9-6.55*math.Log10(src.Z))*math.Log10(distance)

		} else if FreqMHz >= 1500 && FreqMHz < 2000 && distance > 0.05 {
			a := (1.1*math.Log10(FreqMHz)-0.7)*dest.Z - (1.56*math.Log10(FreqMHz) - 0.8)
			result = 46.3 + 33.9*math.Log10(FreqMHz) - 13.82*math.Log10(src.Z) - a + (44.9-6.55*math.Log10(src.Z))*math.Log10(distance) + 3

		} else if FreqMHz >= 150 && FreqMHz < 2000 && distance <= 0.05 {
			result = 20*math.Log10(distance) + 20*math.Log10(FreqMHz) + 32.45

		}
	} else {
		fmt.Printf("distance= %f\n", distance)
		log.Panic("Path loss model does not valid for given frequency")
		return 0, false
	}
	return result, true

}
