package sql_handler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ii-varsha-ii/CloudLab-RDMARoCE-profile/container/data"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	TableName = "rdmaredis"
	Driver = "postgres"
)

var (
	sqlDB *sql.DB
	err   error
	dbName string
)

func InitializeSQLClient(sqlDBUser, sqlDBPassword, sqlDBHost, sqlDBPort, sqlDBName string) error {
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", sqlDBHost, sqlDBPort, sqlDBUser, sqlDBPassword, sqlDBName)
	log.Infof("InitializeSQLClient: SQL connect string: %s", connectStr)
	sqlDB, err = sql.Open(Driver, connectStr)
	if err != nil {
		log.Errorf("Exception while initializing SQL client: %v", err)
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Errorf("Exception while testing SQL client: %v", err)
		return err
	}
	dbName = sqlDBName
	return nil
}

func Close() error {
	return sqlDB.Close()
}

// CreateTable creates the necessary table in the SQL database
func CreateTable() error {
	createTableStmt := `CREATE TABLE IF NOT EXISTS $1 (ID  SERIAL PRIMARY KEY, name TEXT, Message TEXT, SourceType TEXT, MessageSizeInKB INT, WriteTime TIMESTAMP, ReadTime TIMESTAMP, DiffInMs INT)`
	if _, err := sqlDB.Exec(createTableStmt, TableName); err != nil {
		log.Errorf("CreateTable: Exception while creating table %s on SQL DB: %v", TableName, err)
		return err
	}
	return nil
}

// RecordToDatabase records the MessageID, WrittenTime, and ReadTime to the SQL database
func RecordToDatabase(message *data.Message) error {
	insertMsgStmt := `insert into $1 ("Message", "SourceType", "MessageSizeInKB", "WriteTime", "ReadTime", "DiffInMs") values($2, $3, $4, $5, $6, $7)`
	if _, err := sqlDB.Exec(insertMsgStmt, TableName, message.Message, message.GetSourceTypeAsStr(), message.MessageSizeInKB, message.WriteTime, message.ReadTime, message.DiffInMs); err != nil {
		log.Errorf("RecordToDatabase: Exception while writing \"%s\" to SQL DB: %v", message.String(), err)
		return err
	}

	log.Errorf("RecordToDatabase: Successfully wrote message \"%s\" to SQL DB.", message.String())
	return nil
}

// ReadAllData retrieves all data from the SQL database table
func ReadAllData() ([]data.Message, error) {
	result := make([]data.Message, 0)
	selectMsgsStmt := `select * from $1`
	rows, err := sqlDB.Query(selectMsgsStmt, TableName)
	if err != nil {
		log.Errorf("ReadAllData: Exception while reading all data from table %s on SQL DB: %v", TableName, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id              int
			message         string
			sourceType      string
			messageSizeInKB int
			writeTime       time.Time
			readTime        time.Time
			diffInMs        int64
		)

		if err := rows.Scan(&id, &message, &sourceType, &messageSizeInKB, &writeTime, &readTime, &diffInMs); err != nil {
			log.Errorf("ReadAllData: Exception while reading a data from table %s on SQL DB: %v", TableName, err)
			return nil, err
		}
		messageStruct := data.Message{
			ID:              id,
			Message:         message,
			SourceType:      data.String_To_Source_Type[sourceType],
			MessageSizeInKB: messageSizeInKB,
			WriteTime:       writeTime,
			ReadTime:        readTime,
			DiffInMs:        diffInMs,
		}
		result = append(result, messageStruct)
	}
	log.Errorf("ReadAllData: Successfully read all data from table %s on SQL DB.", TableName)
	return result, nil
}
