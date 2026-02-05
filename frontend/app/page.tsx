'use client';

import { useState, useEffect } from 'react';
import { checkAuth, logout, getLoginUrl, buildCourse } from '@/lib/api';
import { FolderPicker } from '@/components/folder-picker';
import { LessonList } from '@/components/lesson-list';
import { VideoPlayer } from '@/components/video-player';
import type { Course } from '@/lib/api';

export default function Home() {
  const [auth, setAuth] = useState<boolean | null>(null);
  const [course, setCourse] = useState<Course | null>(null);
  const [currentFileId, setCurrentFileId] = useState<string | null>(null);
  const [selectedFolders, setSelectedFolders] = useState<Map<string, string>>(new Map());
  const [building, setBuilding] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [showFolderPicker, setShowFolderPicker] = useState(true);

  useEffect(() => {
    checkAuth().then(({ authenticated }) => setAuth(authenticated));
  }, []);

  const handleToggleFolder = (id: string, name: string) => {
    setSelectedFolders((prev) => {
      const next = new Map(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.set(id, name);
      }
      return next;
    });
  };

  const handleBuildCourse = async () => {
    const ids = Array.from(selectedFolders.keys());
    if (ids.length === 0) return;
    setBuilding(true);
    try {
      const c = await buildCourse(ids);
      setCourse(c);
      const firstLesson = c.sections[0]?.lessons[0];
      setCurrentFileId(firstLesson?.fileId ?? null);
    } catch {
      setCourse(null);
    } finally {
      setBuilding(false);
    }
  };

  if (auth === null) {
    return (
      <main className="min-h-screen flex items-center justify-center bg-[var(--bg)]">
        <div className="animate-pulse text-[var(--text-muted)]">Đang tải...</div>
      </main>
    );
  }

  if (!auth) {
    return (
      <main className="min-h-screen flex flex-col items-center justify-center gap-6 p-8 bg-[var(--bg)]">
        <h1 className="text-2xl font-bold">Drive Course Viewer</h1>
        <p className="text-[var(--text-muted)] text-center max-w-md">
          Đăng nhập với Google để xem video trong Drive (Share with me) dưới dạng khóa học.
        </p>
        <a
          href={getLoginUrl()}
          className="px-6 py-3 bg-[var(--accent)] hover:bg-[var(--accent-hover)] rounded font-medium transition-colors"
        >
          Đăng nhập với Google
        </a>
      </main>
    );
  }

  return (
    <main className="min-h-screen flex flex-col bg-[var(--bg)]">
      <header className="h-14 px-4 flex items-center justify-between border-b border-[var(--border)] bg-[var(--bg)] shrink-0">
        <div className="flex items-center gap-4">
          <h1 className="text-base font-bold truncate max-w-md">
            {course ? 'Drive Course Viewer' : 'Drive Course Viewer'}
          </h1>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setSidebarOpen((o) => !o)}
            className="px-3 py-2 text-sm text-[var(--text-muted)] hover:text-white hover:bg-white/5 rounded transition-colors"
          >
            {sidebarOpen ? 'Ẩn nội dung' : 'Nội dung khóa học'}
          </button>
          <button
            onClick={async () => {
              await logout();
              setAuth(false);
            }}
            className="px-3 py-2 text-sm text-[var(--text-muted)] hover:text-white hover:bg-white/5 rounded transition-colors"
          >
            Đăng xuất
          </button>
        </div>
      </header>

      <div className="flex-1 flex min-h-0 relative">
        <div
          className={`absolute inset-0 flex items-center justify-center p-3 md:p-4 transition-[padding] duration-300 ${
            sidebarOpen ? 'pr-80 md:pr-[360px]' : ''
          }`}
        >
          <div
            className={`w-full max-h-full transition-all duration-300 ${
              sidebarOpen ? 'max-w-6xl' : 'max-w-7xl'
            }`}
          >
            <VideoPlayer fileId={currentFileId} />
          </div>
        </div>

        <aside
          className={`fixed top-14 right-0 bottom-0 z-10 bg-[var(--bg-sidebar)] border-l border-[var(--border)] flex flex-col shadow-[-4px_0_24px_rgba(0,0,0,0.3)] transition-transform duration-300 ease-in-out ${
            sidebarOpen ? 'translate-x-0 w-80 md:w-[360px]' : 'translate-x-full w-80 md:w-[360px]'
          }`}
        >
          <div className="w-full h-full flex flex-col overflow-hidden">
            <div className="p-4 border-b border-[var(--border)] shrink-0 space-y-4">
              <div className="flex items-center justify-between">
                <h2 className="text-sm font-bold">Nội dung khóa học</h2>
                <button
                  onClick={() => setSidebarOpen(false)}
                  className="p-1.5 text-[var(--text-muted)] hover:text-white hover:bg-white/5 rounded transition-colors"
                  aria-label="Đóng"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <div className="space-y-3">
                <button
                  onClick={() => setShowFolderPicker((o) => !o)}
                  className="w-full flex items-center justify-between text-left py-2 px-3 rounded hover:bg-white/5 transition-colors"
                >
                  <span className="text-sm font-medium">Chọn thư mục</span>
                  <svg
                    className={`w-4 h-4 text-[var(--text-muted)] transition-transform ${showFolderPicker ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                {showFolderPicker && (
                  <div className="pl-1">
                    <FolderPicker
                      selectedIds={new Set(selectedFolders.keys())}
                      onToggle={handleToggleFolder}
                    />
                    <button
                      onClick={handleBuildCourse}
                      disabled={selectedFolders.size === 0 || building}
                      className="mt-3 w-full py-2.5 bg-[var(--accent)] hover:bg-[var(--accent-hover)] disabled:opacity-50 disabled:cursor-not-allowed rounded text-sm font-medium transition-colors"
                    >
                      {building ? 'Đang tạo...' : `Tạo khóa học (${selectedFolders.size} thư mục)`}
                    </button>
                  </div>
                )}
              </div>
            </div>

            {course && (
              <nav className="flex-1 min-h-0 overflow-y-auto sidebar-scroll">
                <LessonList
                  course={course}
                  currentFileId={currentFileId}
                  onSelectLesson={setCurrentFileId}
                />
              </nav>
            )}
          </div>
        </aside>
      </div>
    </main>
  );
}
