'use client';

import { useState } from 'react';
import type { Course, Section, Lesson } from '@/lib/api';

interface LessonListProps {
  course: Course;
  currentFileId: string | null;
  onSelectLesson: (fileId: string) => void;
}

export function LessonList({ course, currentFileId, onSelectLesson }: LessonListProps) {
  return (
    <div className="py-2">
      {course.sections.map((section, sectionIdx) => (
        <SectionBlock
          key={section.id}
          section={section}
          sectionIdx={sectionIdx}
          currentFileId={currentFileId}
          onSelectLesson={onSelectLesson}
        />
      ))}
    </div>
  );
}

function SectionBlock({
  section,
  sectionIdx,
  currentFileId,
  onSelectLesson,
}: {
  section: Section;
  sectionIdx: number;
  currentFileId: string | null;
  onSelectLesson: (fileId: string) => void;
}) {
  const [open, setOpen] = useState(true);

  const currentIdx = section.lessons.findIndex((l) => l.fileId === currentFileId);
  const progressText = currentIdx >= 0
    ? `${currentIdx + 1}/${section.lessons.length}`
    : `${section.lessons.length} bài`;

  return (
    <div className="border-b border-[var(--border)] last:border-0">
      <button
        onClick={() => setOpen((o) => !o)}
        className="w-full px-4 py-3 flex items-center justify-between text-left hover:bg-white/5 transition-colors gap-2"
      >
        <div className="flex items-center gap-2 min-w-0 flex-1">
          <svg
            className={`w-4 h-4 shrink-0 text-[var(--text-muted)] transition-transform ${open ? 'rotate-180' : ''}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
          <span className="text-sm font-medium truncate">
            Phần {sectionIdx + 1}: {section.title}
          </span>
        </div>
        <span className="text-xs text-[var(--text-muted)] shrink-0">
          {progressText}
        </span>
      </button>
      {open && (
        <ul className="pb-2">
          {section.lessons.map((lesson, lessonIdx) => (
            <LessonItem
              key={lesson.id}
              lesson={lesson}
              isActive={currentFileId === lesson.fileId}
              onSelect={() => onSelectLesson(lesson.fileId)}
            />
          ))}
        </ul>
      )}
    </div>
  );
}

function LessonItem({
  lesson,
  isActive,
  onSelect,
}: {
  lesson: Lesson;
  isActive: boolean;
  onSelect: () => void;
}) {
  return (
    <li>
      <button
        onClick={onSelect}
        className={`w-full px-4 pl-10 py-2.5 flex items-center gap-3 text-left text-sm transition-colors ${
          isActive ? 'bg-[var(--accent)]/20 text-[var(--accent)]' : 'hover:bg-white/5 text-[var(--text-muted)] hover:text-white'
        }`}
      >
        <span className="w-5 shrink-0 flex items-center justify-center">
          {isActive ? (
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clipRule="evenodd" />
            </svg>
          ) : (
            <span className="w-3 h-3 rounded-full border border-current opacity-60" />
          )}
        </span>
        <span className="truncate flex-1">{lesson.title}</span>
      </button>
    </li>
  );
}
