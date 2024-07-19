package filter

import (
	"fmt"

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

func (uf *UserFilter) String() string {
	return fmt.Sprintf(
		"FirstName:%v, LastName:%v, Nickname:%v, Country:%v, Email:%v, Offset:%v, Limit:%v",
		stringValue(uf.FirstName), stringValue(uf.LastName), stringValue(uf.Nickname),
		stringValue(uf.Country), stringValue(uf.Email), int64Value(uf.Offset), int64Value(uf.Limit),
	)
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

func stringValue(s *string) string {
	if s == nil {
		return "empty"
	}
	return *s
}

func int64Value(i *int64) string {
	if i == nil {
		return "empty"
	}
	return fmt.Sprintf("%d", *i)
}

type filterBuilder struct {
	filter *UserFilter
}

func NewFilterBuilder() *filterBuilder {
	return &filterBuilder{filter: &UserFilter{}}
}

func (f *filterBuilder) ByFirstName(firstName *string) *filterBuilder {
	f.filter.FirstName = firstName
	return f
}

func (f *filterBuilder) ByLastName(lastName *string) *filterBuilder {
	f.filter.LastName = lastName
	return f
}

func (f *filterBuilder) ByNickname(nickname *string) *filterBuilder {
	f.filter.Nickname = nickname
	return f
}

func (f *filterBuilder) ByCountry(country *string) *filterBuilder {
	f.filter.Country = country
	return f
}

func (f *filterBuilder) ByEmail(email *string) *filterBuilder {
	f.filter.Email = email
	return f
}

func (f *filterBuilder) WithLimit(limit *int64) *filterBuilder {
	f.filter.Limit = limit
	return f
}

func (f *filterBuilder) WithOffset(offset *int64) *filterBuilder {
	f.filter.Offset = offset
	return f
}

func (f *filterBuilder) Build() *UserFilter {
	return f.filter
}
