package main

import "fmt"

func diff(originalFile string, updatedFile string, og_len int, ud_len int) int{

	var solution = make([][]int,og_len)

	for iter := range solution {
		solution[iter] = make([]int, ud_len)
	}
	
	for iter:=0;iter < og_len;iter++ {
		solution[iter][0] = 0
	}

	for iter:=0;iter < ud_len;iter++ {
		solution[0][iter] = 0
	}

	for i:=1;i<og_len;i++ {
		for j:=1;j<ud_len;j++ {
			if originalFile[i] == updatedFile[j] {
				solution[i][j] = 1 + solution[i-1][j-1]
			} else {
				solution[i][j] = max(solution[i-1][j], solution[i][j-1]) 
			}
		}
	}
	//fmt.Println(solution)
	return solution[og_len-1][ud_len-1]
}

func main() {
	str1 := "Ankit is an amazing guy"
	str2 := "Ankit is a great guy"

	fmt.Println(diff(str1,str2,len(str1),len(str2)))
}
