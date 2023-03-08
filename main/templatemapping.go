package main

import (
	"mailprojesi/datamanagement"
	"net/http"
)

type templateDataStruct struct {
	User  *datamanagement.User
	Flash string
}

func (config *Config) addData(req *http.Request) *templateDataStruct {
	var data *templateDataStruct
	if config.session.Exists(req.Context(), "user") {
		data = config.createStruct(config.session.Get(req.Context(), "user").(datamanagement.User), req)
	} else {
		data = config.createStruct(datamanagement.User{}, req)
	}
	return data
}
func (config *Config) createStruct(user datamanagement.User, req *http.Request) *templateDataStruct {
	dataStruct := &templateDataStruct{
		User:  &user,
		Flash: config.session.PopString(req.Context(), "flash"),
	}

	return dataStruct
}
