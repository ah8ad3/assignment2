package main

import (
	"errors"
	"time"
)

type (
	UserID   int
	UniqueID int
)

// Input is the first flow of accepting users data.
type Input struct {
	// UniqueID per data that system accepts.
	UniqueID UniqueID `json:"unique_id"`
	// UserID is the id of each individual user.
	UserID UserID `json:"user_id"`
}

// NewInput to create new input per new request.
func NewInput(uniqueID UniqueID, userID UserID) Input {
	return Input{UniqueID: uniqueID, UserID: userID}
}

// Accept to import the new data.
// Should return error if user's quota is exceeded.
func (i Input) Accept() error {
	return Storage.Put(i.UniqueID)
}

func (i Input) checkRate(quota Quota) error {
	err := quota.checkMinute()
	if err != nil {
		return err
	}
	err = quota.checkMonthly()
	if err != nil {
		return err
	}
	return nil
}

// Quota to specifying and limit the usage of each user.
type Quota struct {
	// UserID is the id of each individual user.
	UserID         UserID
	monthlyLimiter *RateLimiter
	minuteLimiter  *RateLimiter
}

// NewQuota creates quota for a userID based on month and minute use.
func NewQuota(userID UserID, monthlyLimit int, minuteLimit int) Quota {
	monthlyLimiter := NewRateLimiter(monthlyLimit, time.Hour*24*30)
	minuteLimiter := NewRateLimiter(minuteLimit, time.Minute)
	return Quota{
		UserID:         userID,
		monthlyLimiter: monthlyLimiter,
		minuteLimiter:  minuteLimiter}
}

func (q *Quota) checkMonthly() error {
	if !q.monthlyLimiter.Allow() {
		return errors.New("monthly limit exceeded")
	}
	return nil
}

func (q *Quota) checkMinute() error {
	if !q.minuteLimiter.Allow() {
		return errors.New("minute limit exceeded")
	}
	return nil
}
