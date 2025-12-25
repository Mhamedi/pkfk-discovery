"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { api, APIError } from "@/lib/api";
import { MonacoEditor } from "@/components/editor/MonacoEditor";

interface Draft {
  id: string;
  name: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export default function StudioDraftPage() {
  const params = useParams();
  const draftId = params.draftId as string;
  const [draft, setDraft] = useState<Draft | null>(null);
  const [activeTab, setActiveTab] = useState("overview");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [sqlTemplate, setSqlTemplate] = useState("-- SQL template\nSELECT * FROM information_schema.tables;");

  useEffect(() => {
    if (draftId) {
      loadDraft();
    }
  }, [draftId]);

  const loadDraft = async () => {
    try {
      const data = await api.get<Draft>(`/api/v1/studio/drafts/${draftId}`);
      setDraft(data);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load draft");
      }
    } finally {
      setLoading(false);
    }
  };

  const handleValidate = async () => {
    try {
      const result = await api.post<{ passed: boolean }>(`/api/v1/studio/drafts/${draftId}/validate`, {});
      alert(`Validation ${result.passed ? "passed" : "failed"}`);
    } catch (err) {
      alert("Validation failed");
    }
  };

  const handlePublish = async () => {
    if (!confirm("Are you sure you want to publish this adapter?")) {
      return;
    }

    try {
      const adapter = await api.post<{ id: string }>(`/api/v1/studio/drafts/${draftId}/publish`, {
        maturity_level: "L2",
      });
      alert("Adapter published successfully!");
      window.location.href = `/adapters/${adapter.id}`;
    } catch (err) {
      alert("Failed to publish adapter");
    }
  };

  if (loading) {
    return (
      <div className="p-8">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-1/4"></div>
          <div className="h-64 bg-gray-200 rounded"></div>
        </div>
      </div>
    );
  }

  if (error || !draft) {
    return (
      <div className="p-8">
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          {error || "Draft not found"}
        </div>
      </div>
    );
  }

  const tabs = [
    { id: "overview", label: "Overview" },
    { id: "probe", label: "Probe" },
    { id: "templates", label: "Templates" },
    { id: "validate", label: "Validate" },
    { id: "optimize", label: "Optimize" },
    { id: "publish", label: "Publish" },
    { id: "audit", label: "Audit" },
  ];

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-2">{draft.name}</h1>
      <p className="text-muted-foreground mb-6">
        Status: <span className="px-2 py-1 bg-muted rounded text-xs">{draft.status}</span>
      </p>

      <div className="border-b mb-6">
        <nav className="flex gap-4">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`px-4 py-2 border-b-2 transition-colors ${
                activeTab === tab.id
                  ? "border-primary text-primary"
                  : "border-transparent hover:border-muted-foreground"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      <div className="bg-card border rounded-lg p-6">
        {activeTab === "overview" && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">Name</label>
              <input
                type="text"
                value={draft.name}
                onChange={(e) => setDraft({ ...draft, name: e.target.value })}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
          </div>
        )}

        {activeTab === "probe" && (
          <div className="space-y-4">
            <p className="text-muted-foreground mb-4">
              Probe a connection to discover database characteristics and generate templates.
            </p>
            <button
              onClick={() => {
                api.post(`/api/v1/studio/drafts/${draftId}/probe`, {}).then(() => {
                  alert("Probe job queued");
                });
              }}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
            >
              Run Probe
            </button>
          </div>
        )}

        {activeTab === "templates" && (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">SQL Template</label>
              <MonacoEditor
                value={sqlTemplate}
                onChange={(value) => setSqlTemplate(value || "")}
                language="sql"
                height="500px"
              />
            </div>
          </div>
        )}

        {activeTab === "validate" && (
          <div className="space-y-4">
            <p className="text-muted-foreground mb-4">
              Run validation tests to verify adapter functionality.
            </p>
            <button
              onClick={handleValidate}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
            >
              Run Validation
            </button>
          </div>
        )}

        {activeTab === "optimize" && (
          <div className="space-y-4">
            <p className="text-muted-foreground mb-4">
              Use AI to optimize SQL templates for better performance.
            </p>
            <button
              onClick={() => {
                api.post(`/api/v1/studio/drafts/${draftId}/optimize`, {}).then(() => {
                  alert("Optimization not yet implemented");
                });
              }}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
            >
              Optimize with AI
            </button>
          </div>
        )}

        {activeTab === "publish" && (
          <div className="space-y-4">
            <p className="text-muted-foreground mb-4">
              Package and publish this adapter to the registry.
            </p>
            <button
              onClick={handlePublish}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
            >
              Publish Adapter
            </button>
          </div>
        )}

        {activeTab === "audit" && (
          <div className="space-y-4">
            <p className="text-muted-foreground">Audit trail for this draft will be displayed here.</p>
          </div>
        )}
      </div>
    </div>
  );
}

