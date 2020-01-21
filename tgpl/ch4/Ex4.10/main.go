// Issues prints a table of Github issues matching the search terms
package main

import (
	"exercises-the_go_programming_language/ch4/github"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	category := make(map[string][]string)
	now := time.Now()
	for _, item := range result.Items {
		fmt.Printf("#%-5d %#v\n", item.Number, item.CreatedAt.Format("20060101"))
		s := fmt.Sprintf("#%-5d %#v", item.Number, item.CreatedAt.Format("20060101"))
		if item.CreatedAt.Add(time.Hour * 24 * 30).After(now) {
			category["lessThan1Month"] = append(category["lessThan1Month"], s)
		} else if item.CreatedAt.Add(time.Hour * 24 * 30 * 365).After(now) {
			category["lessThan1Year"] = append(category["lessThan1Year"], s)
		} else {
			category["moreThan1Year"] = append(category["moreThan1Year"], s)
		}
	}

	for k, v := range category {
		fmt.Printf("\n%s\t%d\n", k, len(v))
		for _, issue := range v {
			fmt.Println(issue)
		}
	}
}

// go run main.go repo:golang/go is:open json decoder
