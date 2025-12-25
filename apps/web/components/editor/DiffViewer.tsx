"use client";

import { useEffect, useRef } from "react";
import Editor from "@monaco-editor/react";
import * as monaco from "monaco-editor";

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
  const editorRef = useRef<monaco.editor.IStandaloneDiffEditor | null>(null);

  useEffect(() => {
    if (editorRef.current) {
      const originalModel = monaco.editor.createModel(original, "sql");
      const modifiedModel = monaco.editor.createModel(modified, "sql");
      editorRef.current.setModel({
        original: originalModel,
        modified: modifiedModel,
      });
    }

    return () => {
      if (editorRef.current) {
        const model = editorRef.current.getModel();
        if (model) {
          model.original?.dispose();
          model.modified?.dispose();
        }
      }
    };
  }, [original, modified]);

  return (
    <div className="border rounded-lg overflow-hidden">
      <Editor
        height={height}
        language="sql"
        theme="vs-dark"
        options={{
          readOnly: true,
          minimap: { enabled: true },
          renderSideBySide: true,
          automaticLayout: true,
        }}
        onMount={(editor) => {
          editorRef.current = editor as monaco.editor.IStandaloneDiffEditor;
        }}
      />
    </div>
  );
}

