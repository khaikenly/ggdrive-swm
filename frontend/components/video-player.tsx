'use client';

import { getVideoPreviewUrl } from '@/lib/api';

interface VideoPlayerProps {
  fileId: string | null;
}

export function VideoPlayer({ fileId }: VideoPlayerProps) {
  if (!fileId) {
    return (
      <div className="aspect-video w-full max-h-[calc(100vh-5.5rem)] bg-[#000] rounded shadow-[0_4px_24px_rgba(0,0,0,0.5)] flex items-center justify-center border border-[var(--border)]">
        <p className="text-[var(--text-muted)]">Chọn một bài học để xem</p>
      </div>
    );
  }

  const previewUrl = getVideoPreviewUrl(fileId);

  return (
    <div className="aspect-video w-full max-h-[calc(100vh-5.5rem)] bg-black rounded overflow-hidden shadow-[0_4px_24px_rgba(0,0,0,0.5)] border border-[var(--border)]">
      <iframe
        key={fileId}
        src={previewUrl}
        title="Video player"
        className="w-full h-full"
        allow="autoplay"
        allowFullScreen
      />
    </div>
  );
}
