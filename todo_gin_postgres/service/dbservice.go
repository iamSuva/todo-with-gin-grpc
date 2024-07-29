package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"todowithgin/database"
	"todowithgin/models"
	"todowithgin/utils"

	"github.com/lib/pq"
)

var (
	ErrorNoRows         = errors.New("no data found")
	ErrorInternalserver = errors.New("internal server error")
	ErrorUniqueTitle    = errors.New("title already exists")
	ErrorUpdateTask     = errors.New("failed to update task")
	ErrorDeleteTask     = errors.New("failed to delete task")
	ErrorExecutingQuery = errors.New("error executing query")
	ErrorReadingRows    = errors.New("error in reading row")
	ErrorUniqueUsername=errors.New("username already exists")
    ErrorUnauthorized=errors.New("user not authorized")
	ErrorInvalidUserId=errors.New("invalid user id")
	ErrorPasswordNotMatch=errors.New("password not match")

)

type TaskService interface {
	GetTasks() ([]models.Task, error)
	GetTask(int) (models.Task, error)
	CreateTask( models.Task) error
	UpdateTask( models.Task, int,int) error
	DeleteTask(int,int) error
	SignUpUser(models.User) (int,error)
	LoginUser(models.User) (models.User,error)
}

type DBService struct {
}

func NewDBservice() *DBService {
	return &DBService{}
}

const Unique_code = "23505"
func UniqueUserError(err error)error{
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == Unique_code{
			return ErrorUniqueUsername
		}
		return ErrorInternalserver
	}
	return nil
}
func UniqueErrorHandler(err error) error {
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == Unique_code {
			return ErrorUniqueTitle
		}
		return ErrorInternalserver
	}
	return nil
}
func (d *DBService) GetTasks() ([]models.Task, error) {
	var tasks []models.Task
	sqlQuery := "SELECT id, title, description, isCompleted, createdAt_UTC, updatedAt_UTC,user_Id FROM mytodos"
	rows, err := database.DB.Query(sqlQuery)
	if err != nil {
		return nil, ErrorExecutingQuery
	}
	defer rows.Close()
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt_UTC, &task.UpdatedAt_UTC,&task.UserId)
		if err != nil {
			return nil, ErrorReadingRows
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrorInternalserver
	}

	return tasks, nil
}
func (d *DBService) GetTask(id int) (models.Task, error) {
	var task models.Task
	sqlQuery := "SELECT id,title,description,isCompleted,createdAt_UTC,updatedAt_UTC,user_id from mytodos WHERE id=$1"
	err := database.DB.QueryRow(sqlQuery, id).Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.CreatedAt_UTC, &task.UpdatedAt_UTC,&task.UserId)
	fmt.Println("get err: ",err)
	if err != nil {
		if errors.Is(err,sql.ErrNoRows){
			return task, ErrorNoRows
		}else{
			return task, ErrorInternalserver
		}
	}
	return task, nil
}
func (d *DBService) CreateTask(task models.Task) error {
	title := task.Title
	description := task.Description
	completed := task.IsCompleted
	createdAt_UTC := task.CreatedAt_UTC.Format(time.RFC3339)
	updatedAt_UTC := task.CreatedAt_UTC.Format(time.RFC3339)
	userId:=task.UserId
	fmt.Printf("type of id %d %T",userId,userId);
	sqlStatement := "INSERT INTO mytodos(title,description,isCompleted,createdAt_UTC,updatedAt_UTC,user_Id) VALUES($1,$2,$3,$4,$5,$6) RETURNING id"
	var id int
	err := database.DB.QueryRow(sqlStatement, title, description, completed, createdAt_UTC, updatedAt_UTC,userId).Scan(&id)
	uniqueErr := UniqueErrorHandler(err)
	fmt.Println("create err:", uniqueErr)

	if uniqueErr != nil {
		return uniqueErr
	}
	return nil
}
func (d *DBService) UpdateTask(task models.Task, id int,userid int) error {
	title:=task.Title
	description:=task.Description
	isCompleted:=task.IsCompleted
	sqlQuery := "UPDATE mytodos SET title=$1,description=$2 ,isCompleted=$3,updatedAt_UTC=$4 WHERE id=$5 AND user_id=$6"
	currentTime := time.Now().UTC().Format(time.RFC3339)
	
	result, err := database.DB.Exec(sqlQuery, title, description,isCompleted, currentTime, id,userid)
	
	uniqueErr := UniqueErrorHandler(err)
	
	fmt.Println("update err: ", uniqueErr)
	
	if uniqueErr != nil {
		return uniqueErr
	}
	rowsUpdated, err := result.RowsAffected()
	fmt.Println(rowsUpdated)
	if err != nil {
		return ErrorUpdateTask 
	}
	if rowsUpdated==0{
		return ErrorNoRows
	}
	return nil
}
func (d *DBService) DeleteTask(id int,userid int) error {
	sqlQuery := "DELETE FROM mytodos WHERE id=$1 AND user_id=$2"
	row, err := database.DB.Exec(sqlQuery, id,userid)
	if err != nil {
		return ErrorInternalserver
	}
	rowsaffected, err := row.RowsAffected()
	fmt.Println("del err",err)
	fmt.Println(rowsaffected)
	if err != nil {
		return ErrorDeleteTask
	}
	if rowsaffected==0{
		return ErrorNoRows
	}
	return nil
}

func(d *DBService) SignUpUser(user models.User) (int,error){
     username:=user.Username
	 password:=user.Password
	sqlQuery:="INSERT INTO todousers(username,password) values ($1,$2) RETURNING userid";
	var userid  int;
    
	err:=database.DB.QueryRow(sqlQuery,username,password).Scan(&userid)
    uniqueError:=UniqueUserError(err)
	if uniqueError!=nil{
		return 0,uniqueError
	}
    

	return userid,nil

}
func(d *DBService) LoginUser(user models.User) (models.User,error){
	username:=user.Username
	password:=user.Password
	sqlQuery:="SELECT  userid, username ,password FROM todoUsers WHERE username=$1";
	row:=database.DB.QueryRow(sqlQuery,username)
	var logUser models.User;
    err:=row.Scan(&logUser.UserId,&logUser.Username,&logUser.Password);

	fmt.Println("err: ",err,logUser)
	if err!=nil{
		if errors.Is(err,sql.ErrNoRows){
			return logUser,ErrorNoRows
		}
		return logUser,ErrorInternalserver
	}
	isPasswordMatch:=utils.CheckHashedPassword(password,logUser.Password)
	if !isPasswordMatch{
        return logUser,ErrorPasswordNotMatch
	}
	return logUser,nil
}
