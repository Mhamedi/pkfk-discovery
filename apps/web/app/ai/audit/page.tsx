"use client";

import { useState, useEffect } from "react";
import { api, APIError } from "@/lib/api";

interface AIInteraction {
  id: string;
  user_id: string;
  provider: string;
  model: string;
  prompt_hash: string;
  response_hash: string;
  adapter_draft_id?: string;
  created_at: string;
}

export default function AIAuditPage() {
  const [interactions, setInteractions] = useState<AIInteraction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadInteractions();
  }, []);

  const loadInteractions = async () => {
    try {
      setLoading(true);
      // TODO: Implement AI interactions API endpoint
      // const data = await api.get<AIInteraction[]>("/api/v1/ai/audit");
      setInteractions([]);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load AI interactions");
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
      <h1 className="text-3xl font-bold mb-6">AI Audit</h1>

      {error && (
        <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-md">
          {error}
        </div>
      )}

      {interactions.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground">No AI interactions found</p>
        </div>
      ) : (
        <div className="bg-card border rounded-lg overflow-hidden">
          <table className="w-full">
            <thead className="bg-muted">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold">Timestamp</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">User</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Provider</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Model</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Prompt Hash</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Response Hash</th>
              </tr>
            </thead>
            <tbody>
              {interactions.map((interaction) => (
                <tr key={interaction.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4 text-sm">
                    {new Date(interaction.created_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    {interaction.user_id.substring(0, 8)}...
                  </td>
                  <td className="px-6 py-4">{interaction.provider}</td>
                  <td className="px-6 py-4">{interaction.model}</td>
                  <td className="px-6 py-4 text-sm font-mono text-muted-foreground">
                    {interaction.prompt_hash.substring(0, 16)}...
                  </td>
                  <td className="px-6 py-4 text-sm font-mono text-muted-foreground">
                    {interaction.response_hash.substring(0, 16)}...
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

