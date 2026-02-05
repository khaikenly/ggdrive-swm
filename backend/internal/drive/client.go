package drive

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	googleapi "google.golang.org/api/drive/v3"
)

const (
	MimeTypeFolder = "application/vnd.google-apps.folder"
	VideoMimePrefix = "video/"
)

type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size,omitempty"`
}

type Client struct {
	svc *googleapi.Service
}

func NewClient(svc *googleapi.Service) *Client {
	return &Client{svc: svc}
}

func (c *Client) ListSharedFolders(ctx context.Context) ([]File, error) {
	query := "sharedWithMe=true and mimeType='" + MimeTypeFolder + "' and trashed=false"
	return c.listFiles(ctx, query)
}

func (c *Client) ListFolderChildren(ctx context.Context, folderID string) ([]File, error) {
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	return c.listFiles(ctx, query)
}

func (c *Client) ListVideosInFolder(ctx context.Context, folderID string) ([]File, error) {
	query := fmt.Sprintf("'%s' in parents and mimeType contains '%s' and trashed=false",
		folderID, VideoMimePrefix)
	return c.listFiles(ctx, query)
}

func (c *Client) ListSubfoldersAndVideos(ctx context.Context, folderID string) (folders []File, videos []File, err error) {
	all, err := c.ListFolderChildren(ctx, folderID)
	if err != nil {
		return nil, nil, err
	}

	folders = make([]File, 0, len(all))
	videos = make([]File, 0, len(all))

	for _, f := range all {
		if f.MimeType == MimeTypeFolder {
			folders = append(folders, f)
		} else if strings.HasPrefix(f.MimeType, VideoMimePrefix) {
			videos = append(videos, f)
		}
	}

	return folders, videos, nil
}

// FolderTree chứa thư mục và video trong cây thư mục (bao gồm thư mục con).
type FolderTree struct {
	Folder   File
	Videos   []File
	Children []FolderTree
}

// ListFolderTree đệ quy lấy toàn bộ video và thư mục con từ một thư mục.
func (c *Client) ListFolderTree(ctx context.Context, folderID string) (*FolderTree, error) {
	folder, err := c.GetFile(ctx, folderID)
	if err != nil {
		return nil, err
	}

	subfolders, videos, err := c.ListSubfoldersAndVideos(ctx, folderID)
	if err != nil {
		return nil, err
	}

	children := make([]FolderTree, 0, len(subfolders))
	for _, sf := range subfolders {
		child, err := c.ListFolderTree(ctx, sf.ID)
		if err != nil {
			return nil, err
		}
		children = append(children, *child)
	}

	return &FolderTree{
		Folder:   *folder,
		Videos:   videos,
		Children: children,
	}, nil
}

func (c *Client) DownloadFile(ctx context.Context, fileID string) (*http.Response, error) {
	return c.svc.Files.Get(fileID).
		Context(ctx).
		SupportsAllDrives(true).
		AcknowledgeAbuse(true).
		Download()
}

func (c *Client) GetFile(ctx context.Context, fileID string) (*File, error) {
	f, err := c.svc.Files.Get(fileID).
		Context(ctx).
		Fields("id,name,mimeType,size").
		Do()
	if err != nil {
		return nil, fmt.Errorf("get file %s: %w", fileID, err)
	}

	return &File{
		ID:       f.Id,
		Name:     f.Name,
		MimeType: f.MimeType,
		Size:     f.Size,
	}, nil
}

func (c *Client) listFiles(ctx context.Context, query string) ([]File, error) {
	var result []File
	pageToken := ""

	for {
		call := c.svc.Files.List().
			Context(ctx).
			Q(query).
			PageSize(100).
			Fields("nextPageToken, files(id, name, mimeType, size)").
			OrderBy("name")

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		resp, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("list files: %w", err)
		}

		for _, f := range resp.Files {
			result = append(result, File{
				ID:       f.Id,
				Name:     f.Name,
				MimeType: f.MimeType,
				Size:     f.Size,
			})
		}

		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return result, nil
}
