"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Plus, Play } from "lucide-react";
import { api, APIError } from "@/lib/api";

interface Scan {
  id: string;
  connection_id: string;
  adapter_id: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export default function ScansPage() {
  const [scans, setScans] = useState<Scan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);

  useEffect(() => {
    loadScans();
  }, []);

  const loadScans = async () => {
    try {
      setLoading(true);
      const data = await api.get<Scan[]>("/api/v1/scans");
      setScans(data);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load scans");
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

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Scans</h1>
        <button
          onClick={() => setShowCreateModal(true)}
          className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          <Plus size={20} />
          New Scan
        </button>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-md">
          {error}
        </div>
      )}

      {scans.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">No scans found</p>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
          >
            Create your first scan
          </button>
        </div>
      ) : (
        <div className="bg-card border rounded-lg overflow-hidden">
          <table className="w-full">
            <thead className="bg-muted">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold">ID</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Status</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Created</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Updated</th>
                <th className="px-6 py-3 text-right text-sm font-semibold">Actions</th>
              </tr>
            </thead>
            <tbody>
              {scans.map((scan) => (
                <tr key={scan.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4">
                    <Link
                      href={`/engine/scans/${scan.id}`}
                      className="text-primary hover:underline"
                    >
                      {scan.id.substring(0, 8)}...
                    </Link>
                  </td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 bg-muted rounded text-xs">
                      {scan.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-muted-foreground">
                    {new Date(scan.created_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 text-sm text-muted-foreground">
                    {new Date(scan.updated_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex justify-end">
                      <Link
                        href={`/engine/scans/${scan.id}`}
                        className="p-2 hover:bg-accent rounded"
                        title="View Details"
                      >
                        <Play size={16} />
                      </Link>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {showCreateModal && (
        <ScanModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={() => {
            setShowCreateModal(false);
            loadScans();
          }}
        />
      )}
    </div>
  );
}

function ScanModal({
  onClose,
  onSuccess,
}: {
  onClose: () => void;
  onSuccess: () => void;
}) {
  const [formData, setFormData] = useState({
    connection_id: "",
    adapter_id: "",
    sample_mode: true,
    timeout: 300,
    max_rows: 10000,
  });
  const [submitting, setSubmitting] = useState(false);
  const [connections, setConnections] = useState<any[]>([]);
  const [adapters, setAdapters] = useState<any[]>([]);

  useEffect(() => {
    // Load connections and adapters
    api.get("/api/v1/connections").then(setConnections).catch(() => {});
    api.get("/api/v1/adapters").then(setAdapters).catch(() => {});
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const policy = {
        sample_mode: formData.sample_mode,
        deep_mode: !formData.sample_mode,
        timeout: formData.timeout,
        max_rows: formData.max_rows,
        concurrency: 5,
      };
      await api.post("/api/v1/scans", {
        connection_id: formData.connection_id,
        adapter_id: formData.adapter_id,
        policy_json: policy,
      });
      onSuccess();
    } catch (err) {
      alert("Failed to create scan");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-card rounded-lg p-6 w-full max-w-md">
        <h2 className="text-2xl font-bold mb-4">Create Scan</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Connection</label>
            <select
              required
              value={formData.connection_id}
              onChange={(e) => setFormData({ ...formData, connection_id: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="">Select connection</option>
              {connections.map((conn) => (
                <option key={conn.id} value={conn.id}>
                  {conn.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Adapter</label>
            <select
              required
              value={formData.adapter_id}
              onChange={(e) => setFormData({ ...formData, adapter_id: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="">Select adapter</option>
              {adapters.map((adapter) => (
                <option key={adapter.id} value={adapter.id}>
                  {adapter.name} ({adapter.version})
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.sample_mode}
                onChange={(e) => setFormData({ ...formData, sample_mode: e.target.checked })}
              />
              <span>Sample Mode</span>
            </label>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Timeout (seconds)</label>
            <input
              type="number"
              value={formData.timeout}
              onChange={(e) => setFormData({ ...formData, timeout: parseInt(e.target.value) })}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Max Rows</label>
            <input
              type="number"
              value={formData.max_rows}
              onChange={(e) => setFormData({ ...formData, max_rows: parseInt(e.target.value) })}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
          <div className="flex justify-end gap-2 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border rounded-md"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md disabled:opacity-50"
            >
              {submitting ? "Creating..." : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

