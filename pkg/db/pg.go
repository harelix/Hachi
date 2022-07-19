package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "hachi"
	dbname   = "postgres"
)

func BuildConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

type PostgresProvider struct {
	dbAvailable      bool
	connectionString string
	db               sql.DB
}

var instance *PostgresProvider
var once sync.Once

func GetInstance() *PostgresProvider {
	once.Do(func() {
		instance = &PostgresProvider{}
		instance.dbAvailable = false

	})
	return instance
}

func (client *PostgresProvider) IsAvailable() bool {
	return client.dbAvailable
}

func (client *PostgresProvider) testDBAvailability() {
	//if !client.dbAvailable {
	//	panic(errors.New("Postgres DB is not available under the supplied connection string:" + client.connectionString))
	//}
}

func (client *PostgresProvider) Init(connectionString string) {
	client.connectionString = connectionString
	db, err := sql.Open("postgres", client.connectionString)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(20)
	if err != nil {
		log.Error("Error: Postgres -> open postgres on Init " + err.Error() + " connection string: " + client.connectionString)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		log.Error("Error: Postgres -> open postgres on Init " + err.Error() + " connection string: " + client.connectionString)
	}
	instance.dbAvailable = true
	client.db = *db
}

func (client *PostgresProvider) GetConnection() *sql.DB {
	return &client.db
}

func (client *PostgresProvider) GetConnectionFromString(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Error("Error: Postgres -> open postgres on Init " + err.Error() + " connection string: " + client.connectionString)
		return nil
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		log.Error("Error: Postgres -> open postgres on Init " + err.Error() + " connection string: " + client.connectionString)
		return nil
	}
	instance.dbAvailable = true
	return db
}

func (client *PostgresProvider) ExecuteQuery(query string) string {
	client.testDBAvailability()
	var raw_message string
	err := client.GetConnection().QueryRow(query).Scan(&raw_message)
	if err == sql.ErrNoRows {
		log.Error("Error: Execute Query " + err.Error())
	}
	return raw_message
}

func (client *PostgresProvider) ExecuteServerQuery(connectionString string, query string) string {
	client.testDBAvailability()
	var raw_message string
	err := client.GetConnectionFromString(connectionString).QueryRow(query).Scan(&raw_message)
	if err == sql.ErrNoRows {
		log.Error("Error: Execute Query " + err.Error())
	}
	return raw_message
}

func (client *PostgresProvider) RegisterAgent(agentID string, dedicatedChannel string) (bool, error) {
	fmt.Println(agentID)
	fmt.Println(dedicatedChannel)
	return true, nil
}

/*
func (client *PostgresProvider) WriteRawIncomingMessage(channel string, capsule runners.MessageCapsule) {

	message := capsule.Message
	sqlStatement := `INSERT INTO public."snitch_messages_capture"(message_capsule, incoming_raw_message,  is_valid_message, is_known_message, is_valid_json_message, time)
						VALUES ('%s','%s',%t,%t,%t,now())`

	stringCapsule, _ := json.Marshal(capsule)

	query := fmt.Sprintf(sqlStatement, string(stringCapsule), message, capsule.IsValidMessage, capsule.IsKnownMessage, capsule.IsValidJsonFormatMessage)

	conn := client.GetConnection()
	_, err := conn.Exec(query)
	if err != nil {
		common.Slog().Write("Error: Write Raw Incoming Message" + err.Error())
	}

}

//func (client *PostgresProvider) Write(sinkConf sinks.SinkConfiguration, capsule runners.MessageCapsule, messageJSON string) {
func (client *PostgresProvider) Write(table string, capsule runners.MessageCapsule, messageJSON string) {

	if table == "snitch_messages_flow" {
		client.WriteToMessageFlow(capsule)
	} else {
		client.writeMessageToExistingTable(table, messageJSON)
	}

}

func (client *PostgresProvider) WriteMessageToEspyTable(table string, capsule runners.MessageCapsule) {
	for k, row := range capsule.FlatMessages {
		//str := fmt.Sprintf("%v", capsule.Tags)

		s := make([]string, len(capsule.Tags))
		for i, v := range capsule.Tags {
			s[i] = escapeChars(fmt.Sprint(v))
		}

		str := "ARRAY [ '" + strings.Join(s, "','") + "']"

		insertQuery := fmt.Sprintf(`INSERT INTO public.%s(message, strcuture, tags) VALUES ('%s', '%s', %s);`, table, row, k, str)
		client.ExecuteGenericJsonQuery(insertQuery)
	}
}

func (client *PostgresProvider) writeMessageToExistingTable(table string, message string) {
	insertQuery := fmt.Sprintf(`insert into %s select * from json_populate_record(null::%s, '%s')`, table, table, escapeChars(message))
	client.ExecuteGenericJsonQuery(insertQuery)
}

func (client *PostgresProvider) ExecuteGenericJsonQuery(query string) {
	conn := client.GetConnection()
	_, err := conn.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	return
}
*/
