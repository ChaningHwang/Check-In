package main

import (
    "encoding/json"
	"io"
    "fmt"
    "net/http"
	"log"
	"database/sql"
	"strconv"

	"github.com/golang/protobuf/proto"
    "github.com/gorilla/mux"
    pb "github.com/ChaningHwang/checkin/pkg/proto" // Import the generated Go code from Protocol Buffers
	_ "github.com/go-sql-driver/mysql"
)

var messages = []*pb.Message{
	{Id: "1", Text: "Hello"},
	{Id: "2", Text: "World"},
}

type GetFamilyData struct {
	MemberId int;
  	FamilyId int;
	FirstName string;
	LastName string;
	SchoolID string;
	SchoolName string;
}

type ListEventData struct {
	EventID  int;
	EventName  string;
}

type CheckinRequest struct {
	MemberIDs  []int;
	EventID   string;
}

var db *sql.DB

func GetMessages(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/x-protobuf")

	// Create a MessageList message
	messageList := &pb.MessageList{
		ContextMessage: "This is the context message.",
		Messages:       messages,
	}

	// Encode MessageList as Protocol Buffers and write to response
	data, err := proto.Marshal(messageList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding message list: %v", err)
		return
	}

	if _, err := w.Write(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}
}

func PostMessages(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	var message *pb.Message

	err = json.Unmarshal(body, &message)
	if err!=nil{
	panic(err)
	}
	fmt.Printf("%+v\n", message)

	insertVote := fmt.Sprintf("INSERT INTO messages(id, text) VALUES('%s', '%s')", message.Id, message.Text)

	fmt.Printf("insertVote here: %+v\n", insertVote)

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	result, err := db.Exec(insertVote)
	fmt.Println("here is the result: ", result)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}
}

func connectToMysql() (*sql.DB, error) {
	// Define MySQL connection parameters
	dsn := "root:root1234@tcp(localhost:3306)/testDB"

	// Open a connection to the MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to MySQL database: %v", err)
		return nil, err
	}
	// defer db.Close()

	// Test the connection by pinging the database
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging MySQL database: %v", err)
		return nil, err
	}

	fmt.Println("Connected to MySQL database!")
	return db, nil
}

func CreateFamily(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	var family *pb.Family

	err = json.Unmarshal(body, &family)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%+v\n", family)

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	familyId, err := getLargestFamilyId(db)
	fmt.Println("here is the familyId", familyId)

	for _, member := range family.Members {
		insertQuery := fmt.Sprintf("INSERT INTO members(familyId, firstName, lastName, schoolID, schoolName) VALUES(%d, '%s', '%s', '%s', '%s')", 
		familyId, member.FirstName, member.LastName, member.SchoolID, member.SchoolName)
		fmt.Printf("insertVote here: %+v\n", insertQuery)

		result, err := db.Exec(insertQuery)
		fmt.Println("here is the result for insert: ", result)
		if err!=nil {
			fmt.Fprintf(w, "Error writing response: %v", err)
			continue
		}
	}

	defer db.Close()
}

func GetFamily(w http.ResponseWriter, r *http.Request) {
	familyID := r.Header.Get("familyID")
	if familyID == "" {
		return
	}

	num, err := strconv.Atoi(familyID)
	if err!=nil {
		fmt.Println("Error converting string to int: %v", err.Error())
		return
	}

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM members WHERE FamilyID=%d", 
	num)
	fmt.Printf("select query here: %+v\n", query)

	rows, err := db.Query(query)
	fmt.Println("here is the result for select query: ", rows)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}

	defer rows.Close()

	// An album slice to hold data from returned rows.
    var familyData []GetFamilyData

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var fd GetFamilyData
        if err := rows.Scan(&fd.MemberId, &fd.FamilyId, &fd.FirstName,
            &fd.LastName, &fd.SchoolID, &fd.SchoolName); err != nil {
            return
        }
		fmt.Printf("MemberID: %d, FamilyId: %s, FirstName: %d, LastName: %d, SchoolId: %d, SchooldName: %d\n", 
			fd.MemberId, fd.FamilyId, fd.FirstName, fd.LastName, fd.SchoolID, fd.SchoolName)
        familyData = append(familyData, fd)
    }

	fmt.Println("here is family data: ", familyData)

	json.NewEncoder(w).Encode(familyData)
}

func GetMembersByLastName(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Add("Access-Control-Allow-Origin", "*")
	// w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// w.Header().Add("Access-Control-Allow-Headers", "*")
	// w.Header().Add("Access-Control-Allow-Credentials", "true")

	// if r.Method == "OPTIONS" {
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Println("options request here")
	// 	return
	// }
	HandleCors(w, r)

	lastName := r.Header.Get("LastName")

	fmt.Println("here is", lastName)
	if lastName == "" {
		return
	}

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	defer db.Close()

	query := fmt.Sprintf("SELECT MemberID,FamilyID,FirstName,LastName,SchoolID,SchoolName FROM members WHERE LastName='%s'", 
	lastName)
	fmt.Printf("select query here: %+v\n", query)

	rows, err := db.Query(query)
	fmt.Println("here is the result for select query: ", rows)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}


	// An album slice to hold data from returned rows.
    var familyData []GetFamilyData

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var fd GetFamilyData
        if err := rows.Scan(&fd.MemberId, &fd.FamilyId, &fd.FirstName,
            &fd.LastName, &fd.SchoolID, &fd.SchoolName); err != nil {
            return
        }
		fmt.Printf("MemberID: %d, FamilyId: %s, FirstName: %d, LastName: %d, SchoolId: %d, SchooldName: %d\n", 
			fd.MemberId, fd.FamilyId, fd.FirstName, fd.LastName, fd.SchoolID, fd.SchoolName)
        familyData = append(familyData, fd)
    }

	defer rows.Close()
	fmt.Println("here is family data: ", familyData)

	json.NewEncoder(w).Encode(familyData)
}

func GetMembersBySchoolID(w http.ResponseWriter, r *http.Request) {
	HandleCors(w, r)
	schoolID := r.Header.Get("SchoolID")
	if schoolID == "" {
		return
	}

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	defer db.Close()

	query := fmt.Sprintf("SELECT MemberID,FamilyID,FirstName,LastName,SchoolID,SchoolName FROM members WHERE SchoolID='%s'", 
	schoolID)
	fmt.Printf("select query here: %+v\n", query)

	rows, err := db.Query(query)
	fmt.Println("here is the result for select query: ", rows)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}

	defer rows.Close()

	// An album slice to hold data from returned rows.
    var familyData []GetFamilyData

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var fd GetFamilyData
        if err := rows.Scan(&fd.MemberId, &fd.FamilyId, &fd.FirstName,
            &fd.LastName, &fd.SchoolID, &fd.SchoolName); err != nil {
            return
        }
		fmt.Printf("MemberID: %d, FamilyId: %s, FirstName: %d, LastName: %d, SchoolId: %d, SchooldName: %d\n", 
			fd.MemberId, fd.FamilyId, fd.FirstName, fd.LastName, fd.SchoolID, fd.SchoolName)
        familyData = append(familyData, fd)
    }

	fmt.Println("here is family data: ", familyData)

	json.NewEncoder(w).Encode(familyData)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	
}

func PostEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: insert event data to events table
	body, err := io.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	var event *pb.Event

	err = json.Unmarshal(body, &event)
	if err!=nil{
		panic(err)
		return
	}
	fmt.Printf("%+v\n", event)

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
		return
	}

	defer db.Close()

	eventId, err := getLargestEventId(db)
	fmt.Println("here is the eventId", eventId)

	insertQuery := fmt.Sprintf("INSERT INTO events(eventId, EventName, Description, Time) VALUES(%d, '%s', '%s', '%s')", 
	eventId, event.EventName, event.Description, event.Time)
	fmt.Printf("insert query here: %+v\n", insertQuery)

	result, err := db.Exec(insertQuery)
	fmt.Println("here is the result for insert: ", result)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}

	// TODO: Add new column in members table, bool type default false
	addStatement := fmt.Sprintf("ALTER TABLE members ADD %s BOOLEAN DEFAULT FALSE", "event_"+strconv.Itoa(eventId))
	fmt.Printf("insert query here: %+v\n", addStatement)

	addResult, err := db.Exec(addStatement)
	fmt.Println("here is the result for alter table: ", addResult)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Get event data from events table
	HandleCors(w, r)
	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	defer db.Close()

	query := fmt.Sprintf("SELECT eventID, eventName FROM events")
	fmt.Printf("select query here: %+v\n", query)

	rows, err := db.Query(query)
	fmt.Println("here is the result for select query: ", rows)
	if err!=nil {
		fmt.Fprintf(w, "Error writing response: %v", err)
		return
	}

	defer rows.Close()

	// An album slice to hold data from returned rows.
    var eventsData []ListEventData

    // Loop through rows, using Scan to assign column data to struct fields.
    for rows.Next() {
        var fd ListEventData
        if err := rows.Scan(&fd.EventID, &fd.EventName); err != nil {
            return
        }
		fmt.Printf("EventID: %d, EventName: %s\n", 
			fd.EventID, fd.EventName)
		eventsData = append(eventsData, fd)
    }

	fmt.Println("here is eventsData: ", eventsData)

	json.NewEncoder(w).Encode(eventsData)
}

func CheckinEventByFamily(w http.ResponseWriter, r *http.Request) {
	HandleCors(w, r)
	body, err := io.ReadAll(r.Body)
	if err!=nil {
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	var checkInRequest CheckinRequest

	err = json.Unmarshal(body, &checkInRequest)
	if err!=nil{
		panic(err)
	}
	fmt.Printf("checkIn request: %+v\n", checkInRequest)

	db, err := connectToMysql()
	if err!=nil {
		fmt.Println("err here", err)
	}

	// familyId, err := getLargestFamilyId(db)
	// fmt.Println("here is the familyId", familyId)

	for _, memberID := range checkInRequest.MemberIDs {
		updateQuery := fmt.Sprintf("UPDATE members SET %s=true where MemberID=%d", 
		checkInRequest.EventID, memberID)
		fmt.Printf("insertVote here: %+v\n", updateQuery)

		result, err := db.Exec(updateQuery)
		fmt.Println("here is the result for update: ", result)
		if err!=nil {
			fmt.Fprintf(w, "Error writing response: %v", err)
			continue
		}
	}

	defer db.Close()
}

func getLargestFamilyId(db *sql.DB) (int, error) {
	query := "SELECT MAX(familyId) FROM members"

	rows, err := db.Query(query)
	if err!=nil {
		fmt.Println("Error executing query: %v", err.Error())
		return -1, err
	}

	fmt.Println("hre is t rows: ", rows)

	// if !rows.Next() {
	// 	return 0, nil
	// }

	result := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		result = append(result, id)
	}
	fmt.Println("here is the result: ", result)

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(result)==0 {
		return 0, nil
	}

	num, err := strconv.Atoi(result[0])
	if err!=nil {
		fmt.Println("Error converting string to int: %v", err.Error())
		return -1, err
	}

	return num+1, nil
}

func getLargestEventId(db *sql.DB) (int, error) {
	query := "SELECT MAX(eventId) FROM events"

	rows, err := db.Query(query)
	if err!=nil {
		fmt.Println("Error executing query: %v", err.Error())
		return -1, err
	}

	fmt.Println("hre is t rows: ", rows)

	// if !rows.Next() {
	// 	return 0, nil
	// }

	result := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		result = append(result, id)
	}
	fmt.Println("here is the result: ", result)

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if len(result)==0 {
		return 0, nil
	}

	num, err := strconv.Atoi(result[0])
	if err!=nil {
		fmt.Println("Error converting string to int: %v", err.Error())
		return -1, err
	}

	return num+1, nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HandleCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func main() {
	// err := connectToMysql()
	// if err!=nil {
	// 	fmt.Println("err here", err)
	// }

    router := mux.NewRouter()

    // Define your API routes
	router.HandleFunc("/messages", PostMessages).Methods("POST")
    router.HandleFunc("/messages", GetMessages).Methods("GET")
	router.HandleFunc("/families", CreateFamily).Methods("POST")
	router.HandleFunc("/families", GetFamily).Methods("GET")

	// router.HandleFunc("/members", UpdateFamily).Methods("PUT")

	router.HandleFunc("/GetMembersByLastName", GetMembersByLastName).Methods("GET")
	router.HandleFunc("/GetMembersByLastName", HandleCors).Methods("OPTIONS")
	router.HandleFunc("/GetMembersBySchoolID", GetMembersBySchoolID).Methods("GET")
	router.HandleFunc("/GetMembersBySchoolID", HandleCors).Methods("OPTIONS")

	router.HandleFunc("/events", PostEvents).Methods("POST")
	router.HandleFunc("/events", ListEvents).Methods("GET")

	router.HandleFunc("/checkin", CheckinEventByFamily).Methods("POST")
	router.HandleFunc("/checkin", HandleCors).Methods("OPTIONS")

	// router.Handle("/", corsMiddleware(http.DefaultServeMux))

    // Start the HTTP server
    fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", router)
}
