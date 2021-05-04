/*datamaps is a simple tool to extract from and send data to spreadsheets.
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yulqen/datamaps-go/datamaps"
)

func main() {
	opts := datamaps.ParseOptions()
	if opts.Command == "help" {
		os.Stdout.WriteString(datamaps.Usage)
		os.Exit(0)
	}
	// TODO - removed this to handle "setup" bug below.
	// Check that removing this has no consequences.
	// dbpc := datamaps.NewDBPathChecker(os.UserConfigDir)
	// if !dbpc.Check() {
	// 	datamaps.SetUp()
	// }
	switch opts.Command {
	case "checkdb":
		dbpc := datamaps.NewDBPathChecker(os.UserConfigDir)
		if !dbpc.Check() {
			log.Println("No database file exists. Please run datamaps setup")
		}
	case "import":
		if err := datamaps.ImportToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "datamap":
		if err := datamaps.DatamapToDB(opts); err != nil {
			log.Fatal(err)
		}
	case "setup":
		// BUG This gets called twice if the !dbpc.Check()
		// call above reveals that the config dir is present
		_, err := datamaps.SetUp()
		if err != nil {
			log.Fatal(err)
		}
	case "createmaster":
		if err := datamaps.CreateMaster(opts); err != nil {
			log.Fatal(err)
		}
	case "server":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			fmt.Fprintf(w, "Welcome to datamaps!")
			// or you could write it thus
			// w.Write([]byte("Hello from datamaps"))
		})
		log.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
