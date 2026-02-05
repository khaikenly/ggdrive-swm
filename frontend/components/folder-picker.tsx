'use client';

import { useState, useCallback } from 'react';
import type { DriveFile } from '@/lib/api';
import { listFolders, listFolderChildren } from '@/lib/api';

interface FolderPickerProps {
  selectedIds: Set<string>;
  onToggle: (id: string, name: string) => void;
}

export function FolderPicker({ selectedIds, onToggle }: FolderPickerProps) {
  const [folders, setFolders] = useState<DriveFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const [childrenCache, setChildrenCache] = useState<
    Record<string, { folders: DriveFile[]; videos: DriveFile[] }>
  >({});

  const loadFolders = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const list = await listFolders();
      setFolders(list);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadChildren = useCallback(async (folderId: string) => {
    if (childrenCache[folderId]) {
      setExpanded((prev) => new Set(prev).add(folderId));
      return;
    }
    try {
      const data = await listFolderChildren(folderId);
      setChildrenCache((prev) => ({ ...prev, [folderId]: data }));
      setExpanded((prev) => new Set(prev).add(folderId));
    } catch {
      setError('Failed to load folder contents');
    }
  }, [childrenCache]);

  const toggleExpand = (folderId: string) => {
    if (expanded.has(folderId)) {
      setExpanded((prev) => {
        const next = new Set(prev);
        next.delete(folderId);
        return next;
      });
    } else {
      loadChildren(folderId);
    }
  };

  return (
    <div className="space-y-3">
      <button
        onClick={loadFolders}
        disabled={loading}
        className="w-full px-3 py-2 bg-white/5 hover:bg-white/10 rounded text-sm font-medium disabled:opacity-50 transition-colors"
      >
        {loading ? 'Đang tải...' : 'Tải danh sách thư mục'}
      </button>
      {error && (
        <p className="text-red-400 text-xs">{error}</p>
      )}
      {folders.length === 0 && !loading && (
        <p className="text-[var(--text-muted)] text-xs">
          Nhấn nút trên để xem thư mục Shared with me
        </p>
      )}
      <ul className="space-y-0.5 max-h-48 overflow-y-auto sidebar-scroll">
        {folders.map((f) => (
          <FolderItem
            key={f.id}
            file={f}
            isExpanded={expanded.has(f.id)}
            onToggleExpand={() => toggleExpand(f.id)}
            onSelect={() => onToggle(f.id, f.name)}
            isSelected={selectedIds.has(f.id)}
            children={childrenCache[f.id]}
            loadChildren={() => loadChildren(f.id)}
          />
        ))}
      </ul>
    </div>
  );
}

function FolderItem({
  file,
  isExpanded,
  onToggleExpand,
  onSelect,
  isSelected,
  children,
  loadChildren,
}: {
  file: DriveFile;
  isExpanded: boolean;
  onToggleExpand: () => void;
  onSelect: () => void;
  isSelected: boolean;
  children?: { folders: DriveFile[]; videos: DriveFile[] };
  loadChildren: () => void;
}) {
  const hasChildren = (children?.folders?.length ?? 0) + (children?.videos?.length ?? 0) > 0;

  return (
    <li>
      <div className="flex items-center gap-1.5 py-1 group">
        <button
          onClick={onToggleExpand}
          className="w-5 h-5 flex items-center justify-center text-[var(--text-muted)] hover:text-white shrink-0"
          aria-label={isExpanded ? 'Thu gọn' : 'Mở rộng'}
        >
          {hasChildren ? (isExpanded ? '−' : '+') : ' '}
        </button>
        <label className="flex-1 flex items-center gap-2 cursor-pointer min-w-0">
          <input
            type="checkbox"
            checked={isSelected}
            onChange={onSelect}
            className="rounded border-white/20 shrink-0"
          />
          <span className="text-sm truncate">{file.name}</span>
        </label>
      </div>
      {isExpanded && children && (
        <ul className="ml-5 pl-2 border-l border-[var(--border)] space-y-0.5 py-1">
          {children.folders?.map((sub) => (
            <li key={sub.id} className="text-xs text-[var(--text-muted)]">
              📁 {sub.name}
            </li>
          ))}
          {children.videos?.map((v) => (
            <li key={v.id} className="text-xs text-[var(--text-muted)]">
              🎬 {v.name}
            </li>
          ))}
        </ul>
      )}
    </li>
  );
}
