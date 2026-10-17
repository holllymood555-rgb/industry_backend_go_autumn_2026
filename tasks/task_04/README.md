package main

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	result := Stats{}

	if len(nums) < 2 {
		return result
	}

	diff := nums[1] - nums[0]
	result.Count, result.Sum, result.Min, result.Max = len(nums)-1, diff, diff, diff

	for i := 2; i < len(nums); i++ {
		diff = nums[i] - nums[i-1]

		result.Sum += diff

		if result.Min > diff {
			result.Min = diff
		}

		if result.Max < diff {
			result.Max = diff
		}
	}
	return result
}
