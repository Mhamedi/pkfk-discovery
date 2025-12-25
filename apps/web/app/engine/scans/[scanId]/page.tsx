"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { api, APIError } from "@/lib/api";

interface Scan {
  id: string;
  connection_id: string;
  adapter_id: string;
  status: string;
  policy_json: Record<string, unknown>;
  results_json?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

export default function ScanDetailsPage() {
  const params = useParams();
  const scanId = params.scanId as string;
  const [scan, setScan] = useState<Scan | null>(null);
  const [results, setResults] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (scanId) {
      loadScan();
      loadResults();
    }
  }, [scanId]);

  const loadScan = async () => {
    try {
      const data = await api.get<Scan>(`/api/v1/scans/${scanId}`);
      setScan(data);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load scan");
      }
    } finally {
      setLoading(false);
    }
  };

  const loadResults = async () => {
    try {
      const data = await api.get(`/api/v1/scans/${scanId}/results`);
      setResults(data);
    } catch (err) {
      // Results may not be available yet
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

  if (error || !scan) {
    return (
      <div className="p-8">
        <div className="p-4 bg-destructive/10 text-destructive rounded-md">
          {error || "Scan not found"}
        </div>
      </div>
    );
  }

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Scan Details</h1>

      <div className="space-y-6">
        <div className="bg-card border rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4">Status</h2>
          <div className="space-y-2">
            <div>
              <span className="font-medium">Status: </span>
              <span className="px-2 py-1 bg-muted rounded text-xs">{scan.status}</span>
            </div>
            <div>
              <span className="font-medium">Created: </span>
              <span className="text-muted-foreground">
                {new Date(scan.created_at).toLocaleString()}
              </span>
            </div>
            <div>
              <span className="font-medium">Updated: </span>
              <span className="text-muted-foreground">
                {new Date(scan.updated_at).toLocaleString()}
              </span>
            </div>
          </div>
        </div>

        {results && (
          <div className="bg-card border rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">Results</h2>
            <pre className="bg-muted p-4 rounded overflow-auto">
              {JSON.stringify(results, null, 2)}
            </pre>
          </div>
        )}

        {scan.status === "pending" || scan.status === "running" ? (
          <div className="bg-card border rounded-lg p-6">
            <p className="text-muted-foreground">Scan is in progress. Results will appear here when complete.</p>
          </div>
        ) : scan.status === "completed" && !results ? (
          <div className="bg-card border rounded-lg p-6">
            <p className="text-muted-foreground">No results available.</p>
          </div>
        ) : null}
      </div>
    </div>
  );
}

