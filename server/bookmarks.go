package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

const (
	// StoreBookmarksKey is the key used to store bookmarks in the plugin KV store
	StoreBookmarksKey = "bookmarks"
)

// storeBookmarks stores all the users bookmarks
func (p *Plugin) storeBookmarks(userID string, bmarks *Bookmarks) error {
	jsonBookmarks, jsonErr := json.Marshal(bmarks)
	if jsonErr != nil {
		return jsonErr
	}

	key := getBookmarksKey(userID)
	appErr := p.MattermostPlugin.API.KVSet(key, jsonBookmarks)
	if appErr != nil {
		return errors.New(appErr.Error())
	}

	return nil
}

// getBookmark returns a bookmark with the specified bookmarkID
func (p *Plugin) getBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, err
	}

	_, ok := bmarks.exists(bmarkID)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	if bmarks == nil {
		return nil, nil
	}

	for _, bmark := range bmarks.ByID {
		if bmark.PostID == bmarkID {
			return bmark, nil
		}
	}

	return nil, nil
}

// addBookmark stores the bookmark in a map,
func (p *Plugin) addBookmark(userID string, bmark *Bookmark) (*Bookmarks, error) {

	// get all bookmarks for user
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// no marks, initialize the store first
	if bmarks == nil {
		bmarks = NewBookmarks()
	}

	// user doesn't have any bookmarks add first bookmark and return
	if len(bmarks.ByID) == 0 {
		bmarks.add(bmark)
		if err = p.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark already exists, update ModifiedAt and save
	_, ok := bmarks.exists(bmark.PostID)
	if ok {
		bmarks.updateTimes(bmark.PostID)
		bmarks.updateLabels(bmark)

		if err = p.storeBookmarks(userID, bmarks); err != nil {
			return nil, errors.New(err.Error())
		}
		return bmarks, nil
	}

	// bookmark doesn't exist. Add it
	bmarks.add(bmark)
	if err = p.storeBookmarks(userID, bmarks); err != nil {
		return nil, errors.New(err.Error())
	}
	return bmarks, nil
}

// getBookmarks returns a users bookmarks.  If the user has no bookmarks,
// return nil bookmarks
func (p *Plugin) getBookmarks(userID string) (*Bookmarks, error) {

	// if a user not not have bookmarks, bb will be nil
	bb, appErr := p.API.KVGet(getBookmarksKey(userID))
	if appErr != nil {
		return nil, appErr
	}

	if bb == nil {
		return nil, nil
	}

	// return initialized bookmarks
	bmarks := NewBookmarks()
	jsonErr := json.Unmarshal(bb, &bmarks)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return bmarks, nil
}

// ByPostCreateAt returns an array of bookmarks sorted by post.CreateAt times
func (p *Plugin) ByPostCreateAt(bmarks *Bookmarks) ([]*Bookmark, error) {
	// build temp map
	tempMap := make(map[int64]string)
	for _, bmark := range bmarks.ByID {
		post, appErr := p.API.GetPost(bmark.PostID)
		if appErr != nil {
			return nil, appErr
		}
		tempMap[post.CreateAt] = bmark.PostID
	}

	// sort post.CreateAt (keys)
	keys := make([]int, 0, len(tempMap))
	for k := range tempMap {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	// reconstruct the bookmarks in a sorted array
	var bookmarks []*Bookmark
	for _, k := range keys {
		bmark := bmarks.ByID[tempMap[int64(k)]]
		bookmarks = append(bookmarks, bmark)
	}

	return bookmarks, nil
}

// getBookmarksWithLabel return a Bookmarks with bookmarks that contains given
// label specified by labelName
func (p *Plugin) getBookmarksWithLabel(userID, labelName string) (*Bookmarks, error) {
	bmarks, err := p.getBookmarks(userID)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	if bmarks == nil {
		return nil, nil
	}

	bmarksWithLabel := NewBookmarks()

	labelIDs, err := p.getLabelIDsFromNames(userID, []string{labelName})
	labelID := labelIDs[0]

	for _, bmark := range bmarks.ByID {
		if bmark.hasLabels(bmark) {
			for _, id := range bmark.LabelIDs {
				if labelID == id {
					bmarksWithLabel.add(bmark)
				}
			}
		}
	}

	return bmarksWithLabel, nil
}

// deleteBookmark deletes a bookmark from the store
func (p *Plugin) deleteBookmark(userID, bmarkID string) (*Bookmark, error) {
	bmarks, err := p.getBookmarks(userID)
	var bmark *Bookmark
	if err != nil {
		return bmark, errors.New(err.Error())
	}

	if bmarks == nil {
		return bmark, errors.New(fmt.Sprintf("User doesn't have any bookmarks"))
	}

	_, ok := bmarks.exists(bmarkID)
	if !ok {
		return bmark, errors.New(fmt.Sprintf("Bookmark `%v` does not exist", bmarkID))
	}

	bmark = bmarks.get(bmarkID)

	bmarks.delete(bmarkID)
	p.storeBookmarks(userID, bmarks)

	return bmark, nil
}

// getBookmarkLabelIDs returns an array of label UUIDs for a given bookmark
func (p *Plugin) getBookmarkLabelIDs(userID string, bmarkID string) ([]string, error) {
	bmark, err := p.getBookmark(userID, bmarkID)
	if err != nil {
		return nil, err
	}

	return bmark.LabelIDs, nil
}

func (p *Plugin) getBookmarkLabelNames(userID string, bmark *Bookmark) ([]string, error) {
	var labelNames []string
	for _, id := range bmark.LabelIDs {
		name, err := p.getLabelNameByID(userID, id)
		if err != nil {
			return nil, err
		}
		labelNames = append(labelNames, name)
	}
	return labelNames, nil
}

func getBookmarksKey(userID string) string {
	return fmt.Sprintf("%s_%s", StoreBookmarksKey, userID)
}
