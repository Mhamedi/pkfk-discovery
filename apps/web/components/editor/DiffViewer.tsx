"use client";

import { DiffEditor } from "@monaco-editor/react";

interface DiffViewerProps {
  original: string;
  modified: string;
  height?: string;
}

export function DiffViewer({
  original,
  modified,
  height = "400px",
}: DiffViewerProps) {
  return (
    <div className="border rounded-lg overflow-hidden">
      <DiffEditor
        height={height}
        language="sql"
        theme="vs-dark"
        original={original}
        modified={modified}
        options={{
          readOnly: true,
          minimap: { enabled: true },
          automaticLayout: true,
        }}
      />
    </div>
  );
}

