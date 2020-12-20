package main

import(
	"fmt"
	s3 "AlertBot/pkg/s3"
	
)

func main(){
	res, err := s3.HandleRequest("testUser123")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}