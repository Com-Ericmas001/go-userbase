package userbase

import (
	"time"

	"github.com/leonelquinteros/gorand"
)

//TokenSuccessResponse is a response
type TokenSuccessResponse struct {
	Success bool
	Token   Token
}

//ConnectUserResponse is a response
type ConnectUserResponse struct {
	TokenResponse TokenSuccessResponse
	IDUser        int64
}

func invalidTokenSuccessResponse() TokenSuccessResponse {
	return TokenSuccessResponse{
		Success: false,
		Token: Token{
			ID:         "",
			ValidUntil: time.Now()}}
}
func invalidConnectUserResponse() ConnectUserResponse {
	return ConnectUserResponse{
		IDUser:        0,
		TokenResponse: invalidTokenSuccessResponse()}
}

func (context DbContext) newTokenSuccessResponse(id int64) ConnectUserResponse {

	uuid, err := gorand.UUID()
	checkErr(err)

	token := Token{ID: uuid, ValidUntil: time.Now().Add(time.Minute * time.Duration(10))}

	stmt, err := context.Db.Prepare("INSERT INTO UserTokens(IdUser, Token, Expiration) VALUES(?, ?, ?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id, token.ID, token.ValidUntil)
	checkErr(err)

	return ConnectUserResponse{
		IDUser: id,
		TokenResponse: TokenSuccessResponse{
			Success: true,
			Token:   token}}
}
