package main

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	s := Stats{}
	n := len(nums) 
	s.Min = nums[1] - nums[0]

	if n < 2 {
    return Stats{}
	}
		
	for i:= 1; i <= n; i++ {
		diff := nums[i] - nums[i-1]
		s.Count ++
		s.Sum += diff
		
		if s.Min > diff {
			s.Min = diff
		}
		if s.Max < diff {
			s.Max = diff
		}
	}
	return s
}
