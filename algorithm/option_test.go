package algorithm

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	nameBaseline   = string(SearchAlgorithmBaseline)
	nameUnweighted = string(SearchAlgorithmUnweighted)
	nameUnknown    = "freshness"
)

func TestSearchAlgorithmNames(t *testing.T) {
	Convey("Given the registered search algorithms", t, func() {
		Convey("When SearchAlgorithmNames is called", func() {
			names := SearchAlgorithmNames()

			Convey("Then it should return one name per registered algorithm, in canonical order", func() {
				So(names, ShouldResemble, []string{"baseline", "unweighted"})
				So(names, ShouldHaveLength, len(AllSearchAlgorithms()))
			})
		})
	})
}

func TestParseSearchAlgorithm(t *testing.T) {
	Convey("Given a known algorithm name", t, func() {
		Convey("When ParseSearchAlgorithm is called with surrounding whitespace and mixed case", func() {
			algo, err := ParseSearchAlgorithm("  Unweighted ")

			Convey("Then it should resolve to the unweighted algorithm", func() {
				So(err, ShouldBeNil)
				So(algo, ShouldEqual, SearchAlgorithmUnweighted)
			})
		})
	})

	Convey("Given an unknown algorithm name", t, func() {
		Convey("When ParseSearchAlgorithm is called", func() {
			algo, err := ParseSearchAlgorithm(nameUnknown)

			Convey("Then it should error rather than falling back to baseline", func() {
				So(err, ShouldNotBeNil)
				So(algo, ShouldBeEmpty)
			})

			Convey("Then the error should list the valid algorithms", func() {
				So(err.Error(), ShouldEqual, `unknown algorithm "freshness": valid algorithms are baseline, unweighted`)
			})
		})
	})
}

func TestParseSearchAlgorithms(t *testing.T) {
	Convey("Given no algorithm names", t, func() {
		Convey("When ParseSearchAlgorithms is called with an empty list", func() {
			algorithms, err := ParseSearchAlgorithms(nil)

			Convey("Then it should resolve to every registered algorithm", func() {
				So(err, ShouldBeNil)
				So(algorithms, ShouldResemble, AllSearchAlgorithms())
			})
		})

		Convey("When ParseSearchAlgorithms is called with only blank names", func() {
			algorithms, err := ParseSearchAlgorithms([]string{"", "  "})

			Convey("Then it should resolve to every registered algorithm", func() {
				So(err, ShouldBeNil)
				So(algorithms, ShouldResemble, AllSearchAlgorithms())
			})
		})
	})

	Convey("Given a single algorithm name", t, func() {
		Convey("When ParseSearchAlgorithms is called", func() {
			algorithms, err := ParseSearchAlgorithms([]string{nameUnweighted})

			Convey("Then it should resolve to that algorithm alone", func() {
				So(err, ShouldBeNil)
				So(algorithms, ShouldResemble, []SearchAlgorithm{SearchAlgorithmUnweighted})
			})
		})
	})

	Convey("Given several algorithm names out of canonical order", t, func() {
		Convey("When ParseSearchAlgorithms is called", func() {
			algorithms, err := ParseSearchAlgorithms([]string{nameUnweighted, nameBaseline})

			Convey("Then it should resolve them in canonical order", func() {
				So(err, ShouldBeNil)
				So(algorithms, ShouldResemble, []SearchAlgorithm{SearchAlgorithmBaseline, SearchAlgorithmUnweighted})
			})
		})
	})

	Convey("Given a repeated algorithm name", t, func() {
		Convey("When ParseSearchAlgorithms is called", func() {
			algorithms, err := ParseSearchAlgorithms([]string{nameBaseline, nameBaseline})

			Convey("Then the duplicate should be discarded", func() {
				So(err, ShouldBeNil)
				So(algorithms, ShouldResemble, []SearchAlgorithm{SearchAlgorithmBaseline})
			})
		})
	})

	Convey("Given a list holding one unknown algorithm name", t, func() {
		Convey("When ParseSearchAlgorithms is called", func() {
			algorithms, err := ParseSearchAlgorithms([]string{nameBaseline, nameUnknown})

			Convey("Then the whole list should be rejected", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, `unknown algorithm "freshness"`)
				So(algorithms, ShouldBeNil)
			})
		})
	})
}
