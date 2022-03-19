#!/usr/bin/env bash


data() {
	echo "A  B  C"
	sleep 1
	echo "1 1 12345"
	sleep 1 
	echo "111 1 12"
	sleep 1
	echo "2 111 11"
	sleep 1
	echo "33333 211123 12313"
	sleep 1
	echo "11321321 211123 12313"
	sleep 1
	echo "1 6348343483499 12313"
	sleep 1
	echo "878732476234786213486 1 1"
}

echo "1st test case: column"
echo "---------------------"

data | go run main.go column -t


data2() {
cat <<EOF
1
99
3
4
EOF
sleep 1
cat <<EOF
45
3
77
33
EOF
sleep 1
cat <<EOF
56
78
44
1111
9
EOF
}

echo
echo "2nd test case: sort"
echo "-------------------"

data2 | go run main.go sort -n


echo
echo "3nd test case: wc"
echo "-------------------"
data | go run main.go wc
