//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

type dateAndCount struct {
	date  int
	count int
}

func ComputeLoad(guests []Guest) []Load {
	var dateAndCountSlice []dateAndCount
	var loadSlice []Load
	for _, g := range guests {
		dateAndCountSlice = append(dateAndCountSlice, dateAndCount{g.CheckInDate, 1}, dateAndCount{g.CheckOutDate, -1})
	}

	sort.Slice(dateAndCountSlice, func(i, j int) bool {
		return dateAndCountSlice[i].date < dateAndCountSlice[j].date
	})

	sumCount := 0
	for i := 0; i < len(dateAndCountSlice)-1; i++ {
		if dateAndCountSlice[i].date == dateAndCountSlice[i+1].date {
			sumCount += dateAndCountSlice[i].count
			if i == len(dateAndCountSlice)-2 {
				loadSlice = append(loadSlice, Load{dateAndCountSlice[i].date, sumCount + dateAndCountSlice[i+1].count})
			}
		} else {
			sumCount += dateAndCountSlice[i].count
			if len(loadSlice) > 0 && loadSlice[len(loadSlice)-1].GuestCount != sumCount {
				loadSlice = append(loadSlice, Load{dateAndCountSlice[i].date, sumCount})
			}
			if len(loadSlice) == 0 {
				loadSlice = append(loadSlice, Load{dateAndCountSlice[i].date, sumCount})
			}
			if i == len(dateAndCountSlice)-2 {
				loadSlice = append(loadSlice, Load{dateAndCountSlice[i+1].date, sumCount + dateAndCountSlice[i+1].count})
			}
		}

	}

	return loadSlice
}
