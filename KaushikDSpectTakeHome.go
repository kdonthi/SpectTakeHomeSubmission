package main

import (
	"fmt"
	"math/rand"
	"time"
)

func max (a int, b int) (int) {
	if a > b {
		return a
	} else {
		return b
	}
}

func mergeSort (list []int) (sortedArray []int) {
	var list_len int = len(list)
	if list_len <= 1 {
		return list
	} else {
		var lista []int = mergeSort(list[:list_len/2])
		var listb []int = mergeSort(list[list_len/2:])
		var newlist []int = []int{}
		for len(lista) > 0 || len(listb) > 0 {
			if (len(lista) > 0 && len(listb) > 0) {
				if lista[0] < listb[0] {
					newlist = append(newlist, lista[0])
					lista = lista[1:]
				} else {
					newlist = append(newlist, listb[0])
					listb = listb[1:]
				}
			} else if (len(lista) > 0) {
				newlist = append(newlist, lista...) //unpacks list a
				break
			} else {
				newlist = append(newlist, listb...)
				break
			}
		}
		return newlist
	}
}

func updatedPassengers(passengers []int) (updatedpassengers []int) {
	var newpassengers []int = []int{}
	for _, passenger := range(passengers) {
		if passenger != -1 {
			newpassengers = append(newpassengers, passenger)
		}
	}
	return (newpassengers)
}

func needs_extra_driver(drivers [][]int, min_time int, max_time int, ride_time int) (bool) {
	times := make(map[int]bool) //has bool value to say whether time is covered or not
	for i := min_time; i <= max_time; i++ {
		times[i] = false
	}

	for _, driverSchedule := range(drivers) {
		if len(driverSchedule) == 0 { //if any fully free drivers, then we are all taken care of
			return false
		}
		for passenger := 0; passenger < len(driverSchedule); passenger++ {
			if passenger == 0 {
				for recovery := 0; recovery < driverSchedule[0]; recovery++ {
					if !times[recovery] {
						times[recovery] = true
					}
				}
			}
			if passenger < len(driverSchedule) - 1 {
				for recovery := driverSchedule[passenger] + ride_time; recovery < driverSchedule[passenger + 1]; recovery++ {
					if !times[recovery] {
						times[recovery] = true
					}
				}
			} else {
				for recovery := driverSchedule[passenger] + ride_time; recovery <= max_time; recovery++ { //if at the end, go to the end
					if !times[recovery] {
						times[recovery] = true
					}
				}
			}
		}
	}

	for _, value := range(times) {
		if !value {
			return true
		}
	}
	return false
}

func unbalancedDrivers(passengers []int, ride_time int, recovery int) (int, [][]int) {
	var passengers_copy []int = make([]int, len(passengers))
	var num_passengers int = len(passengers)
	var drivers [][]int = [][]int{}
	var driver []int
	var wait_time int = ride_time + recovery
	var min_valid_time int
	var append_val int
	var min_time int
	var max_time int

	copy(passengers_copy, passengers) //passing passengers by value into passengers_copy
	if num_passengers == 0 { //if no passengers, then no need to have backup drivers either
		return 0, drivers
	}
	passengers_copy = mergeSort(passengers_copy) //sorting in case not sorted already
	min_time = passengers_copy[0]
	max_time = passengers_copy[len(passengers_copy)-1]
	for len(passengers_copy) != 0 {
		driver = []int{}
		for index, passenger := range(passengers_copy) {
			if index == 0 { //shifting first index as far back as possible
				append_val = passenger
				if append_val <= 5 {
					append_val = 0
				} else {
					append_val -= 5
				}
				driver = append(driver, append_val)
				passengers_copy[index] = -1 //preparing for removal
			} else {
				min_valid_time = max(driver[len(driver) - 1] + wait_time, passenger - 5) //pickup time in window has to satisfy two conditions - 1) be wait_time away from previous passenger, and 2) be inside the 5 min window of each the passenger time
				if (min_valid_time <= passenger + 5) { //only adding if under the upper bound
					driver = append(driver, min_valid_time)
					passengers_copy[index] = -1 //preparing for removal by updatedPassengers
				}
			}
		}
		passengers_copy = updatedPassengers(passengers_copy)
		drivers = append(drivers, driver)
	}
	if needs_extra_driver(drivers, min_time, max_time + (ride_time - 1), ride_time) {
		drivers = append(drivers, []int{}) //adding backup driver
		return len(drivers) - 1, drivers //not including redundant driver in number of drivers
	} else { //no redundant driver included
		return len(drivers), drivers
	}
}

func balancedDrivers(passengers []int, ride_time int, recovery int) (int, [][]int) {
	var num_passengers int = len(passengers)
	var num_drivers, unbalanced_drivers = unbalancedDrivers(passengers, ride_time, recovery)
	var balanced_drivers [][]int = [][]int{}
	var append_val int
	var new_passenger_times []int
	var extra_driver bool = false
	//using times potentially shifted by 5 min from unbalancedDrivers and sorting them
	for _, driver := range(unbalanced_drivers) {
		for _, passenger_time := range(driver) {
			new_passenger_times = append(new_passenger_times, passenger_time)
		}
	}
	new_passenger_times = mergeSort(new_passenger_times)

	//returning blank array if no passengers and therefore no drivers
	if num_drivers == 0 {
		return num_drivers, unbalanced_drivers
	}
	if len(unbalanced_drivers) != num_drivers {
		extra_driver = true
		num_drivers++
	}
	//put the first driver_num passengers in different lists
	for passenger_index := 0; passenger_index < num_drivers; passenger_index++ {
		if passenger_index < num_passengers {
			append_val = new_passenger_times[passenger_index]
			balanced_drivers = append(balanced_drivers, []int{append_val}) //putting the first num_driver passengers into different lists and reducing as much as possible
		} else {
			balanced_drivers = append(balanced_drivers, []int{})
		}
	}
	//getting rid of already added passengers or returning driver list if we have added all the passengers already
	if num_passengers > num_drivers {
		new_passenger_times = new_passenger_times[num_drivers:]
	} else {
		return len(balanced_drivers) - 1, balanced_drivers
	}

	var max_direct_time_diff int //eventually becomes max difference between given current passenger's pickup time and each driver's latest pickup time - this encourages drivers to take passengers further away in time from their previous passenger and therefore balances the system
	var max_driver_index int //eventually becomes index of driver with max_direct_time_diff
	var max_passenger_time int //eventually becomes passenger's time with max_direct_time_diff
	var direct_time_diff int //current difference between given current passenger pickup time and current driver's latest pick time
	for _, passenger := range(new_passenger_times) {
		max_passenger_time = -1
		max_driver_index = -1
		max_direct_time_diff = -1
		for driver_index, driver := range(balanced_drivers) {
			direct_time_diff = passenger - driver[len(driver) - 1]
			if direct_time_diff > max_direct_time_diff {
				max_direct_time_diff = direct_time_diff
				max_passenger_time = passenger
				max_driver_index = driver_index
			}
		}
		balanced_drivers[max_driver_index] = append(balanced_drivers[max_driver_index], max_passenger_time)
	}
	if extra_driver {
		return len(balanced_drivers) - 1, balanced_drivers
	} else {
		return len(balanced_drivers), balanced_drivers
	}
}

func main() {
	var array1 []int = []int{0,20,40,100}
	fmt.Println("Array: ", array1)
	print("Unbalanced: ")
	fmt.Println(unbalancedDrivers(array1, 15, 5))
	print("Balanced: ")
	fmt.Println(balancedDrivers(array1, 15, 5))

	var array2 []int = []int{0,15,30,45}
	fmt.Println()
	fmt.Println("Array: ", array2)
	print("Unbalanced: ")
	fmt.Println(unbalancedDrivers(array2, 15, 5))
	print("Balanced: ")
	fmt.Println(balancedDrivers(array2, 15, 5))

	var array3 []int = []int{0,0,10}
	fmt.Println()
	fmt.Println("Array: ", array3)
	print("Unbalanced: ")
	fmt.Println(unbalancedDrivers(array3,15,5))
	print("Balanced: ")
	fmt.Println(balancedDrivers(array3, 15, 5))

	fmt.Println("\nExtra Cases: ")
	var array4 []int = []int{}
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		arr_len := rand.Int() % 10
		for j := 0; j < arr_len; j++ {
			array4 = append(array4, rand.Int() % 30)
		}
		fmt.Println("\nArray: ", array4)
		print("Unbalanced: ")
		fmt.Println(unbalancedDrivers(array4, 20, 5))
		//fmt.Println("Array again: ", array4)
		print("Balanced: ")
		fmt.Println(balancedDrivers(array4, 20, 5))
		array4 = nil
	}

	var array5 []int = []int{0,15,30,45}
	fmt.Println("\nArray: ", array5)
	print("Unbalanced: ")
	fmt.Println(unbalancedDrivers(array5, 15, 30))
	print("Balanced: ")
	fmt.Println(balancedDrivers(array5, 15, 30))
}
