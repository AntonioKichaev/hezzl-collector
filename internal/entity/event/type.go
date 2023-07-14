package event

import (
	"fmt"
	"time"
)

type ClickhouseEvent struct {
	Id          int       `json:"Id"`
	CampaignId  int       `json:"CampaignId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description,omitempty"`
	Priority    int       `json:"Priority,omitempty"`
	Removed     bool      `json:"Removed,omitempty"`
	EventTime   time.Time `json:"EventTime"`
}

func (cl *ClickhouseEvent) String() string {
	return fmt.Sprintf(
		"ID:%d, CampaignId:%d, Name:%s, Description:%s, Priority:%d, Removed:%v, EventTime:%s",
		cl.Id,
		cl.CampaignId,
		cl.Name,
		cl.Description,
		cl.Priority,
		cl.Removed,
		cl.EventTime,
	)

}
