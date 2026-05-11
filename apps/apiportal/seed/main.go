package seed

import "github.com/golangci/golangci-lint/pkg/config"




func main(){
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}


	

}