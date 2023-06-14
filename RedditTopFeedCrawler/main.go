package main

import (
	route "app/RedditTopFeedCrawler/route"
	"app/RedditTopFeedCrawler/server"
	"log"
	"net/http"
)

func main() {

	routehandler, err := server.InitRouteHandler()

	if err != nil {
		log.Fatal(err)
		return
	}

	router := route.NewRouter(routehandler)

	log.Fatal(http.ListenAndServe(":8000", router))

	// Either we can hit or listen the server like i did or we can create option to hit the particular api

	// var option int
	// fmt.Println("Select an option:")
	// fmt.Println("1. Crawl Post")
	// fmt.Println("2. Search Profile")
	// fmt.Println("3. Insights")
	// fmt.Print("Enter your choice: ")
	// fmt.Scanln(&option)

	// switch option {
	// case 1:
	// 	log.Fatal(http.ListenAndServe(":8000", router))
	// case 2:
	// 	log.Fatal(http.ListenAndServe(":8000", router))
	// case 3:
	// 	log.Fatal(http.ListenAndServe(":8000", router))
	// default:
	// 	fmt.Println("Invalid option.")
	// 	os.Exit(0)

	// }

}
