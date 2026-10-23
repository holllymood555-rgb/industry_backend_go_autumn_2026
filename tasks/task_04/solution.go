package main

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	s := Stats{}

	if len(nums) < 2 {
		return Stats{}
	}

	firsPair := nums[1] - nums[0]
	s.Count = 1
	s.Sum = firsPair
	s.Min = firsPair
	s.Max = firsPair

	for i := 2; i < len(nums); i++ {
		s.Count ++
		Pair := nums[i] - nums[i-1]
		s.Sum += Pair
		if s.Min > Pair {
			s.Min = Pair
		}
		if s.Max < Pair {
			s.Max = Pair
		}
	}
	return s
}
