"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { api, APIError } from "@/lib/api";

interface Adapter {
  id: string;
  name: string;
  vendor: string;
  db_family: string;
  version: string;
  maturity_level: string;
  bundle_path: string;
  signature: string;
  created_at: string;
  updated_at: string;
}

export default function AdapterDetailsPage() {
  const params = useParams();
  const adapterId = params.adapterId as string;
  const [adapter, setAdapter] = useState<Adapter | null>(null);
  const [activeTab, setActiveTab] = useState("overview");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (adapterId) {
      loadAdapter();
    }
  }, [adapterId]);

  const loadAdapter = async () => {
    try {
      const data = await api.get<Adapter>(`/api/v1/adapters/${adapterId}`);
      setAdapter(data);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load adapter");
      }
    } finally {
      setLoading(false);
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

  if (error || !adapter) {
    return (
      <div className="p-8">
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          {error || "Adapter not found"}
        </div>
      </div>
    );
  }

  const tabs = [
    { id: "overview", label: "Overview" },
    { id: "templates", label: "Templates" },
    { id: "capabilities", label: "Capabilities" },
    { id: "tests", label: "Tests" },
    { id: "releases", label: "Releases" },
    { id: "security", label: "Security/Audit" },
  ];

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-2">{adapter.name}</h1>
      <p className="text-muted-foreground mb-6">
        {adapter.vendor} • {adapter.db_family} • v{adapter.version}
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
              <h3 className="font-semibold mb-2">Metadata</h3>
              <dl className="grid grid-cols-2 gap-4">
                <div>
                  <dt className="text-sm text-muted-foreground">Name</dt>
                  <dd>{adapter.name}</dd>
                </div>
                <div>
                  <dt className="text-sm text-muted-foreground">Vendor</dt>
                  <dd>{adapter.vendor}</dd>
                </div>
                <div>
                  <dt className="text-sm text-muted-foreground">DB Family</dt>
                  <dd>{adapter.db_family}</dd>
                </div>
                <div>
                  <dt className="text-sm text-muted-foreground">Version</dt>
                  <dd>{adapter.version}</dd>
                </div>
                <div>
                  <dt className="text-sm text-muted-foreground">Maturity Level</dt>
                  <dd>
                    <span className="px-2 py-1 bg-muted rounded text-xs">
                      {adapter.maturity_level}
                    </span>
                  </dd>
                </div>
                <div>
                  <dt className="text-sm text-muted-foreground">Created</dt>
                  <dd>{new Date(adapter.created_at).toLocaleDateString()}</dd>
                </div>
              </dl>
            </div>
          </div>
        )}

        {activeTab === "templates" && (
          <div>
            <p className="text-muted-foreground">SQL templates will be displayed here</p>
          </div>
        )}

        {activeTab === "capabilities" && (
          <div>
            <p className="text-muted-foreground">Capabilities will be displayed here</p>
          </div>
        )}

        {activeTab === "tests" && (
          <div>
            <p className="text-muted-foreground">Test definitions and results will be displayed here</p>
          </div>
        )}

        {activeTab === "releases" && (
          <div>
            <p className="text-muted-foreground">Version history will be displayed here</p>
          </div>
        )}

        {activeTab === "security" && (
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold mb-2">Signature</h3>
              <p className="text-sm text-muted-foreground font-mono break-all">
                {adapter.signature}
              </p>
            </div>
            <div>
              <h3 className="font-semibold mb-2">Bundle Path</h3>
              <p className="text-sm text-muted-foreground">{adapter.bundle_path}</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

