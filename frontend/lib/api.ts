const API_BASE = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:8080';

const fetchOpts: RequestInit = {
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
};

export async function checkAuth(): Promise<{ authenticated: boolean }> {
  const res = await fetch(`${API_BASE}/api/auth/me`, { ...fetchOpts, method: 'GET' });
  if (!res.ok) return { authenticated: false };
  return res.json();
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/api/auth/logout`, { ...fetchOpts, method: 'POST' });
}

export function getLoginUrl(): string {
  return `${API_BASE}/api/auth/login`;
}

export interface DriveFile {
  id: string;
  name: string;
  mimeType: string;
  size?: number;
}

export async function listFolders(): Promise<DriveFile[]> {
  const res = await fetch(`${API_BASE}/api/folders`, fetchOpts);
  if (!res.ok) throw new Error('Failed to list folders');
  const data = await res.json();
  return data.folders ?? [];
}

export async function listFolderChildren(folderId: string): Promise<{
  folders: DriveFile[];
  videos: DriveFile[];
}> {
  const res = await fetch(`${API_BASE}/api/folders/${folderId}/children`, fetchOpts);
  if (!res.ok) throw new Error('Failed to list folder children');
  return res.json();
}

export interface Lesson {
  id: string;
  title: string;
  fileId: string;
  duration?: number;
}

export interface Section {
  id: string;
  title: string;
  folderId: string;
  lessons: Lesson[];
}

export interface Course {
  sections: Section[];
}

export async function buildCourse(folderIds: string[]): Promise<Course> {
  const res = await fetch(`${API_BASE}/api/courses/build`, {
    ...fetchOpts,
    method: 'POST',
    body: JSON.stringify({ folderIds }),
  });
  if (!res.ok) throw new Error('Failed to build course');
  return res.json();
}

export function getVideoStreamUrl(fileId: string): string {
  return `${API_BASE}/api/videos/${fileId}/stream`;
}

/** Google Drive preview URL - stream trong trình xem của Drive, không cần quyền Download */
export function getVideoPreviewUrl(fileId: string): string {
  return `https://drive.google.com/file/d/${fileId}/preview`;
}
