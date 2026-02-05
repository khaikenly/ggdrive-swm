package drive

import (
	"context"
	"fmt"
)

type Lesson struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FileID   string `json:"fileId"`
	Duration int    `json:"duration,omitempty"`
}

type Section struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	FolderID string   `json:"folderId"`
	Lessons  []Lesson `json:"lessons"`
}

type Course struct {
	Sections []Section `json:"sections"`
}

func (c *Client) BuildCourse(ctx context.Context, folderIDs []string) (Course, error) {
	sections := make([]Section, 0)

	for _, folderID := range folderIDs {
		tree, err := c.ListFolderTree(ctx, folderID)
		if err != nil {
			return Course{}, err
		}
		flattenTreeToSections(tree, &sections)
	}

	for i := range sections {
		for j := range sections[i].Lessons {
			sections[i].Lessons[j].ID = formatLessonID(i, j)
		}
	}

	return Course{Sections: sections}, nil
}

func flattenTreeToSections(tree *FolderTree, sections *[]Section) {
	lessons := make([]Lesson, 0, len(tree.Videos))
	for _, v := range tree.Videos {
		lessons = append(lessons, Lesson{Title: v.Name, FileID: v.ID})
	}

	if len(lessons) > 0 {
		*sections = append(*sections, Section{
			ID:       tree.Folder.ID,
			Title:    tree.Folder.Name,
			FolderID: tree.Folder.ID,
			Lessons:  lessons,
		})
	}

	for _, child := range tree.Children {
		flattenTreeToSections(&child, sections)
	}
}

func formatLessonID(sectionIdx, lessonIdx int) string {
	return fmt.Sprintf("%d-%d", sectionIdx, lessonIdx)
}
