package main

import (
	"bytes"
	"encoding/base32"
	"encoding/json"
	"fmt"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

// StoreLabelsKey is the key used to store labels in the plugin KV store
const StoreLabelsKey = "labels"

// storeLabels stores all the users labels
func (l *Labels) storeLabels() error {
	bb, jsonErr := json.Marshal(l)
	if jsonErr != nil {
		return jsonErr
	}

	key := getLabelsKey(l.userID)
	appErr := l.api.KVSet(key, bb)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getNameFromID returns the Name of a Label
func (l *Labels) getNameFromID(id string) (string, error) {
	label, err := l.get(id)
	if err != nil {
		return "", err
	}

	return label.Name, nil
}

// getLabels returns a users labels
func (l *Labels) getLabels() (*Labels, error) {
	// if a user does not have labels, bb will be nil
	bb, appErr := l.api.KVGet(getLabelsKey(l.userID))
	if appErr != nil {
		return nil, appErr
	}

	if bb == nil {
		return l, nil
	}

	jsonErr := json.Unmarshal(bb, l)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return l, nil
}

// getLabelsForUser returns a users labels
func (l *Labels) getLabelsForUser() (*Labels, error) {
	// if a user does not have labels, bb will be nil
	bb, appErr := l.api.KVGet(getLabelsKey(l.userID))
	if appErr != nil {
		return nil, appErr
	}
	if bb == nil {
		return nil, nil
	}

	var labels *Labels
	jsonErr := json.Unmarshal(bb, &labels)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return labels, nil
}

// getLabelByName returns a label with the provided label name
func (l *Labels) getLabelByName(labelName string) *Label {
	for _, label := range l.ByID {
		if label.Name == labelName {
			return label
		}
	}
	return nil
}

// getIDFromName returns a label name with the corresponding label ID
func (l *Labels) getIDFromName(labelName string) (string, error) {
	if l == nil {
		return "", errors.New("user does not have any labels")
	}

	// return the labelId if found
	for id, label := range l.ByID {
		if label.Name == labelName {
			return id, nil
		}
	}
	return "", errors.New(fmt.Sprintf("Label: `%s` does not exist", labelName))
}

// addLabel stores a label into the users label store
func (l *Labels) addLabel(labelName string) (*Labels, error) {
	// check if name already exists
	label := l.getLabelByName(labelName)

	// User already has label with this labelName
	if label != nil {
		return nil, errors.New(fmt.Sprintf("Label with name `%s` already exists", label.Name))
	}

	labelID := NewID()
	label = &Label{
		Name: labelName,
	}
	l.add(labelID, label)

	if err := l.storeLabels(); err != nil {
		return nil, err
	}

	return l, nil
}

// deleteByID deletes a label from the store
func (l *Labels) deleteByID(labelID string) error {
	l.delete(labelID)
	if err := l.storeLabels(); err != nil {
		return err
	}
	return nil
}

func getLabelsKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreLabelsKey, userID)
}

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h769")

// NewID is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// with the padding stripped off.
func NewID() string {
	var b bytes.Buffer
	encoder := base32.NewEncoder(encoding, &b)
	_, _ = encoder.Write(uuid.NewRandom())
	encoder.Close()
	b.Truncate(26) // removes the '==' padding
	return b.String()
}
