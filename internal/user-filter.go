package internal

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Filter interface {
	ToBSON() bson.M
}

type UserFilter struct {
	FirstName *string
	LastName  *string
	Nickname  *string
	Country   *string
	Email     *string
	Offset    *int64
	Limit     *int64
}

func (u *UserFilter) ToBSON() bson.M {
	query := bson.M{}

	if u.FirstName != nil {
		query["first_name"] = *u.FirstName
	}

	if u.LastName != nil {
		query["last_name"] = u.LastName
	}

	if u.Nickname != nil {
		query["nickname"] = *u.Nickname
	}

	if u.Country != nil {
		query["country"] = *u.Country
	}

	if u.Email != nil {
		query["email"] = *u.Email
	}

	return query
}

type FilterBuilder struct {
	filter *UserFilter
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{filter: &UserFilter{}}
}

func (f *FilterBuilder) ByFirstName(firstName *string) *FilterBuilder {
	f.filter.FirstName = firstName
	return f
}

func (f *FilterBuilder) ByLastName(lastName *string) *FilterBuilder {
	f.filter.LastName = lastName
	return f
}

func (f *FilterBuilder) ByNickname(nickname *string) *FilterBuilder {
	f.filter.Nickname = nickname
	return f
}

func (f *FilterBuilder) ByCountry(country *string) *FilterBuilder {
	f.filter.Country = country
	return f
}

func (f *FilterBuilder) ByEmail(email *string) *FilterBuilder {
	f.filter.Email = email
	return f
}

func (f *FilterBuilder) Build() *UserFilter {
	return f.filter
}
