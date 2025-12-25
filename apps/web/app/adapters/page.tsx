"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Plus } from "lucide-react";
import { api, APIError } from "@/lib/api";

interface Adapter {
  id: string;
  name: string;
  vendor: string;
  db_family: string;
  version: string;
  maturity_level: string;
  created_at: string;
}

export default function AdaptersPage() {
  const [adapters, setAdapters] = useState<Adapter[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAdapters();
  }, []);

  const loadAdapters = async () => {
    try {
      setLoading(true);
      const data = await api.get<Adapter[]>("/api/v1/adapters");
      setAdapters(data);
      setError(null);
    } catch (err) {
      if (err instanceof APIError) {
        setError(err.message);
      } else {
        setError("Failed to load adapters");
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
        <h1 className="text-3xl font-bold">Adapter Registry</h1>
        <Link
          href="/studio/new"
          className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          <Plus size={20} />
          New Adapter
        </Link>
      </div>

      {error && (
        <div className="mb-4 p-4 bg-destructive/10 text-destructive rounded-md">
          {error}
        </div>
      )}

      {adapters.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-muted-foreground mb-4">No adapters found</p>
          <Link
            href="/studio/new"
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
          >
            Create your first adapter
          </Link>
        </div>
      ) : (
        <div className="bg-card border rounded-lg overflow-hidden">
          <table className="w-full">
            <thead className="bg-muted">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold">Name</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Vendor</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">DB Family</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Version</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Maturity</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Created</th>
              </tr>
            </thead>
            <tbody>
              {adapters.map((adapter) => (
                <tr key={adapter.id} className="border-t hover:bg-muted/50">
                  <td className="px-6 py-4">
                    <Link
                      href={`/adapters/${adapter.id}`}
                      className="text-primary hover:underline"
                    >
                      {adapter.name}
                    </Link>
                  </td>
                  <td className="px-6 py-4">{adapter.vendor}</td>
                  <td className="px-6 py-4">{adapter.db_family}</td>
                  <td className="px-6 py-4">{adapter.version}</td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 bg-muted rounded text-xs">
                      {adapter.maturity_level}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-muted-foreground">
                    {new Date(adapter.created_at).toLocaleDateString()}
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

