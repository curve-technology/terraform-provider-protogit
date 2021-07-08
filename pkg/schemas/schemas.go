package schemas

import (
	"fmt"

	"github.com/curve-technology/terraform-provider-protogit/pkg/store"
)

type Section string

const (
	SectionKey   Section = "key"
	SectionValue Section = "value"
)

// Subject returns the subject given a topic and a section (key or value).
func Subject(topic string, section Section) string {
	return fmt.Sprintf("%s-%s", topic, section)
}

// BuildRecords takes a list of entries and a storer and compiles a list of records.
func BuildRecords(storer store.Storer, entries Entries) (records Records, err error) {

	dagBuilder := NewDAGsBuilder()

	// Get fileDescriptor for each message
	for _, entry := range entries {
		subject := Subject(entry.Topic, entry.Section)

		fileDescriptor, err := storer.GetFileDescriptor(entry.Filepath)
		if err != nil {
			return records, err
		}

		err = dagBuilder.AddDag(fileDescriptor, subject)
		if err != nil {
			return records, err
		}
	}

	records, err = dagBuilder.Records()
	if err != nil {
		return records, err
	}

	return records, nil
}
