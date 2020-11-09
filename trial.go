package main
import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Article struct {
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Headline string       `json:"Headline"`
	Content     string    `json:"content"`
	Link         string    `json:"link"`
}

type Output struct {
	Status       string    `json:"status"`
	Total 			 int       	`json:"total"`
	Articles     []Article `json:"articles"`
}

type Search struct {
	search  string
	next   int
	total int
	results    Output
}



func (s *Search) current_page() int {
	if s.next == 1 {
		return s.next
	}
	return s.next - 1
}

func (s *Search) prev_page() int {
	return s.current_page() - 1
}

func (s *Search) is_last() bool {
	return s.next >= s.total
}







//Search Requests
func searchHandler(w http.ResponseWriter, r *http.Request) {

	u, error := url.Parse(r.URL.String())
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	params := u.Query()
	searchKey := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is: ", searchKey)
	search := &Search{}
	search.SearchKey = searchKey

	next, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Unexpected server error", http.StatusInternalServerError)
		return
	}

	search.next = next
	pageSize := 20

	endpoint := fmt.Sprintf("https://localhost/v2/all_articles?q=%s&pageSize=%d&page=%d&apiKey=%s&language=en", url.QueryEscape(search.SearchKey), pageSize, search.NextPage, apiKey)

	resp, err := http.Get(endpoint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

var temp = template.Must(template.ParseFiles("index.html"))
var apiKey string

func indexHandler(w http.ResponseWriter, r *http.Request) {
	temp.Execute(w, nil)
}

func main() {
	//FIXME: Fetch these inputs from some input file
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	apiKey = os.Getenv("API_KEY")

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	fmt.Println("Starting Server...")
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}
