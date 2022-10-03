# bdc-gini

Toy problem using [gini](https://github.com/go-air/gini) to solve a rotation scheduling problem.

## Problem

Given a set of people, a number of weeks to schedule for and an optional initial condition,
identify a weekly rotating schedule of two people with the following conditions:
- each week must have two different people on the rotation (only one person or no people scheduled for a given week is not an acceptable solution)
- a person cannot be on the rotation two weeks in a row
- a person cannot be on successive rotations with the same partner
