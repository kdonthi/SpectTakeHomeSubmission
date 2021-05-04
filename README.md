# Passenger Pickup Algorithm

## Introduction
The goal of this program is to take in passenger times, in minutes, and output a list of lists with the minimum number of drivers and each driver's pickup schedule. 

## Specs
Some of the special things about this program include: 
1. Each driver's driving time and rest time can be set as variables (```balancedDrivers(driversLists, driving_time, rest_time)```)
2. A redundant driver is added if the union of all the rest times doesn't cover all the times from the beginning to the end
3. Allows the passengers to be picked up in a 5 minute window from their given time (done by finding minimum possible time in the 5 minute window)
4. The drivers' passenger lists are balanced so that no one driver has too many passengers (done by finding maximum time distance from last passenger of each driver, not using the 5 minute window)

There are two functions to pay attention to.

1. The first function is ```unbalancedDrivers```, which gives the list of drivers. This should give you an idea of what the lists look like before they are balanced. I decided not to include the potential redundant driver in the "minimum" amount of drivers.  
2. The second function is ```balancedDrivers```, which has all of the bonus points as the first function in addition to balancing the passengers amount the different driver lists, including the redundant one. <b> This is the function that delivers the balanced result. </b>

I have created some randomized test cases and feel free to ask if I can clarify anything (you can run the file with ```go run KaushikDSpectTakeHome.go```).
