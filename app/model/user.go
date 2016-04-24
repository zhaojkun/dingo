package model

import (
	"database/sql"
	"github.com/dinever/dingo/app/utils"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id       int64
	Name     string
	Slug     string
	Avatar   string
	Email    string
	Image    string // NULL
	Cover    string // NULL
	Bio      string // NULL
	Website  string // NULL
	Location string // NULL
	Role     int    //1 = Administrator, 2 = Editor, 3 = Author, 4 = Owner
}

func (u *User) Save(hashedPassword string, createdBy int64) error {
	_, err := InsertUser(u.Name, u.Slug, hashedPassword, u.Email, u.Image, u.Cover, time.Now(), createdBy)
	if err != nil {
		return err
	}
	//	err = InsertRoleUser(u.Role, userId)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (u *User) UpdateUser(updatedById int64) error {
	err := UpdateUser(u.Id, u.Name, u.Slug, u.Email, u.Image, u.Cover, u.Bio, u.Website, u.Location, time.Now(), updatedById)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ChangePassword(password string) error {
	hashedPassword, err := EncryptPassword(password)
	if err != nil {
		return err
	}
	WriteDB, err := db.Begin()
	if err != nil {
		WriteDB.Rollback()
		return err
	}
	_, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, u.Id)
	if err != nil {
		WriteDB.Rollback()
		return err
	}
	return WriteDB.Commit()
}

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) CheckPassword(password string) bool {
	hashedPassword, err := GetHashedPasswordForUser(string(u.Email))
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false
	}
	return true
}

func scanUser(user *User, row *sql.Row) error {
	var (
		nullImage    sql.NullString
		nullCover    sql.NullString
		nullBio      sql.NullString
		nullWebsite  sql.NullString
		nullLocation sql.NullString
	)
	err := row.Scan(&user.Id, &user.Name, &user.Slug, &user.Email, &nullImage, &nullCover, &nullBio, &nullWebsite, &nullLocation)
	user.Avatar = utils.Gravatar(user.Email, "150")
	user.Image = nullImage.String
	user.Cover = nullCover.String
	user.Bio = nullBio.String
	user.Website = nullWebsite.String
	user.Location = nullLocation.String
	return err
}

func GetHashedPasswordForUser(email string) ([]byte, error) {
	var hashedPassword []byte
	row := db.QueryRow(stmtGetHashedPasswordByEmail, email)
	err := row.Scan(&hashedPassword)
	if err != nil {
		return []byte{}, err
	}
	return hashedPassword, nil
}

func GetUserById(id int64) (*User, error) {
	user := new(User)
	// Get user
	row := db.QueryRow(stmtGetUserById, id)
	err := scanUser(user, row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserBySlug(slug string) (*User, error) {
	user := new(User)
	// Get user
	row := db.QueryRow(stmtGetUserBySlug, slug)
	err := scanUser(user, row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByName(name string) (*User, error) {
	user := new(User)
	// Get user
	row := db.QueryRow(stmtGetUserByName, name)
	err := scanUser(user, row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := new(User)
	// Get user
	row := db.QueryRow(stmtGetUserByEmail, email)
	err := scanUser(user, row)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func InsertUser(name string, slug string, password string, email string, image string, cover string, created_at time.Time, created_by int64) (int64, error) {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	result, err := writeDB.Exec(stmtInsertUser, nil, uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen), name, slug, password, email, image, cover, created_at, created_by, created_at, created_by)
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	userId, err := result.LastInsertId()
	if err != nil {
		writeDB.Rollback()
		return 0, err
	}
	return userId, writeDB.Commit()
}

func InsertRoleUser(role_id int, user_id int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtInsertRoleUser, nil, role_id, user_id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UpdateUser(id int64, name string, slug string, email string, image string, cover string, bio string, website string, location string, updated_at time.Time, updated_by int64) error {
	writeDB, err := db.Begin()
	if err != nil {
		writeDB.Rollback()
		return err
	}
	_, err = writeDB.Exec(stmtUpdateUser, name, slug, email, image, cover, bio, website, location, updated_at, updated_by, id)
	if err != nil {
		writeDB.Rollback()
		return err
	}
	return writeDB.Commit()
}

func UserChangeEmail(email string) bool {
	var count int64
	row := db.QueryRow(stmtGetUsersCountByEmail, email)
	err := row.Scan(&count)
	if count > 0 || err != nil {
		return false
	}
	return true
}

func GetNumberOfUsers() (int64, error) {
	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	err := row.Scan(&count)
	return count, err
}

func CreateNewUser(email, name, password string) error {
	user := new(User)
	user.Email = email
	user.Name = name
	hashedPassword, err := EncryptPassword(password)
	if err != nil {
		return err
	}
	return user.Save(hashedPassword, 0)
}
