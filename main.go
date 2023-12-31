package main

import (
	"fmt"

	"github.com/mingyuanc/GovTech-Technical/routes"
	"github.com/mingyuanc/GovTech-Technical/utils"
)

func main() {
	fmt.Println("Starting Server...")
	fmt.Println("Connecting to database...")
	db := utils.Connect()
	fmt.Println("Connected to database...")
	fmt.Println("Starting backend...")
	routes.Run(db)
}
