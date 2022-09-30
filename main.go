package main

import (
	"fmt"

	"github.com/go-air/gini"
	"github.com/go-air/gini/z"
)

func main() {
	bdc_names := []string{"Alex","Anik","Ankita","Austin","Per","Catherine","Jordan","Per","Tim","Tyler"}
	num_weeks := 11
	repeat_min := len(bdc_names) / 4 // minimum number of weeks between repeating assignments
	pair_min := len(bdc_names) / 3 // minimum number of weeks between repeating a pairing
	base := len(bdc_names) + num_weeks
	//rules:
	// Each week must have 2 BDC
	// The same person cannot be both BDC
	// If a person is selected on one week, don't prefer them for adjacent weeks
	// If a person is paired with someone on a week, try to pair with someone else on the next week
	g := gini.New()
	var lit = func(p1, p2, week int) z.Lit {
		// This is a way to enumerate the choices available for each week
		// week - number of the week we are on
		// p1   - position of first selection in the bdc_names list
		// p2   - position of the second bdc_names on the list
		n := week + base * (p1 + base * p2)
		return z.Var(n + 1).Pos()
	}


	var foranywith = func(p1, w int)[] z.Lit {
		// all combinations with a particular person on a given week
		lits := []z.Lit{}
		for i := 0; i < len(bdc_names); i++ {
			if i == p1 {
				continue
			}
			if i < p1 {
				lits = append(lits, lit(i, p1, w))
				continue
			}
			lits = append(lits, lit(p1, i, w))
		}
		return lits
	}

	var openClause bool

	var printClause = func (l z.Lit) {
		if l.String() == "0" {
			fmt.Println()
			openClause = false
			return
		}
		if openClause {
			fmt.Print(" OR ")
		}
		openClause = true
		v := l.Dimacs()
		if l.Sign() == -1 {
			fmt.Print("NOT ")
			v *= -1
		}
		v -= 1
		w := v % base
		p2 := v / base / base
		p1 := (v - w - p2*base*base) /base
		fmt.Printf("w%d(%s,%s)", w, bdc_names[p1], bdc_names[p2])
	}

	var add = func(l z.Lit) {
		g.Add(l)
//		printClause(l)
	}

	for w := 0; w < num_weeks; w++ {
		var nonEmptyClause bool
		for p1 := range bdc_names {
			// Exclude possibility of selecting same person twice
			for p2 := p1 + 1; p2 < len(bdc_names); p2++ {
				nonEmptyClause = true
				add(lit(p1, p2, w))
			}
		}
		if nonEmptyClause{
			// terminate with 0 to end a set of OR'd literals
			add(0)
		}
	}

	// The same pair should not repeat
	for p1 := range bdc_names {
		for p2 := p1 + 1; p2 < len(bdc_names); p2++ {
			for w := 0; w < num_weeks; w++ {
				a := lit(p1, p2, w)
				for i := 1; i <= pair_min; i++ {
					if w-i < 0 {
						continue
					}
					b := lit(p1, p2, w-i)
					add(a.Not())
					add(b.Not())
					add(0)
				}
			}
		}
	}

	// The same person should not serve consecutive weeks
	for p1 := range bdc_names {
		for w := 0; w < num_weeks; w++ {
			setA := foranywith(p1, w)
			for i := 1; i <= repeat_min; i++ {
				if w-i < 0 {
					continue
				}
				setB := foranywith(p1, w-i)
				for _, a := range setA {
					for _, b := range setB {
						// ~A OR ~B = ~(A AND B)
						add(a.Not())
						add(b.Not())
						add(0)
					}
				}
			}
		}
	}

// TODO: initialize with known pairs here

	if g.Solve() != 1 {
		reasons := make([]z.Lit, 0)
		result := g.Why(reasons)
		fmt.Println("Failed to schedule BDC:")
		for _, a := range append(result, reasons...) {
			printClause(a)
		}
		return
	}

out:
	for w := 0; w < num_weeks; w++ {
		fmt.Printf("Week %d: ", w)
		for p1 := range bdc_names {
			for p2 := p1 + 1; p2 < len(bdc_names); p2++ {
				if g.Value(lit(p1, p2, w)) {
					fmt.Println(bdc_names[p1], ",", bdc_names[p2])
					continue out
				}
			}
		}
	}
}


