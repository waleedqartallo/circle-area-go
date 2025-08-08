package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	Radius float64
	Area   float64
}

var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head><title>Circle Area Calculator</title></head>
<body>
    <h2>Calculate Circle Area</h2>
    <form method="POST" action="/calculate">
        Radius: <input type="text" name="radius">
        <input type="submit" value="Calculate">
    </form>

    {{if .Result}}
        <p>Area = {{.Result}}</p>
    {{end}}

    <h3>Saved Results</h3>
    <form method="POST" action="/clear" style="margin-bottom:10px;">
    <input type="submit" value="Clear All">
</form>

    <table border="1">
        <tr><th>Radius</th><th>Area</th></tr>
        {{range .Records}}
        <tr><td>{{.Radius}}</td><td>{{.Area}}</td></tr>
        {{end}}
    </table>
</body>
</html>
`))

var db *sql.DB

func main() {
	var err error

	// ⚠️ عدل الباسورد لو عندك باسورد للـ MySQL
	db, err = sql.Open("mysql", "root:1234567@tcp(127.0.0.1:3306)/db1")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/clear", handleClear)
	http.HandleFunc("/", showForm)
	http.HandleFunc("/calculate", handleCalc)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func showForm(w http.ResponseWriter, r *http.Request) {
	records := getAllRecords()
	tmpl.Execute(w, map[string]interface{}{
		"Records": records,
	})
}

func handleCalc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var radius float64
	_, err := fmt.Sscanf(r.FormValue("radius"), "%f", &radius)
	if err != nil {
		http.Error(w, "Please enter a valid radius.", http.StatusBadRequest)
		return
	}

	area := math.Pi * radius * radius
	saveRecord(radius, area)

	records := getAllRecords()
	tmpl.Execute(w, map[string]interface{}{
		"Result":  fmt.Sprintf("%.2f", area),
		"Records": records,
	})
}

func saveRecord(radius, area float64) {
	_, err := db.Exec("INSERT INTO table1 (radius, area) VALUES (?, ?)", radius, area)
	if err != nil {
		log.Println("❌ Error saving:", err)
	} else {
		log.Println("✅ Saved to DB: radius =", radius, "area =", area)
	}
}

func getAllRecords() []Record {
	rows, err := db.Query("SELECT radius, area FROM table1")
	if err != nil {
		log.Println("Error fetching:", err)
		return nil
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Radius, &r.Area); err != nil {
			continue
		}
		records = append(records, r)
	}
	return records
}
func handleClear(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("DELETE FROM table1")
	if err != nil {
		log.Println("❌ Error clearing table:", err)
	} else {
		log.Println("✅ All records deleted.")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
