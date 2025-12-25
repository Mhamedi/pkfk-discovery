"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";

export default function NewStudioPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    name: "",
    vendor: "",
    db_family: "postgresql",
    description: "",
  });
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const draft = await api.post<{ id: string }>("/api/v1/studio/drafts", {
        name: formData.name,
        status: "draft",
      });
      router.push(`/studio/${draft.id}`);
    } catch (err) {
      alert("Failed to create draft");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Create New Adapter</h1>

      <div className="bg-card border rounded-lg p-6 max-w-2xl">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Adapter Name</label>
            <input
              type="text"
              required
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
              placeholder="e.g., PostgreSQL Adapter"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Vendor</label>
            <input
              type="text"
              value={formData.vendor}
              onChange={(e) => setFormData({ ...formData, vendor: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
              placeholder="e.g., pkfk-discovery"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Database Family</label>
            <select
              value={formData.db_family}
              onChange={(e) => setFormData({ ...formData, db_family: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="postgresql">PostgreSQL</option>
              <option value="mysql">MySQL</option>
              <option value="sqlserver">SQL Server</option>
              <option value="oracle">Oracle</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Description</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
              rows={4}
              placeholder="Brief description of the adapter"
            />
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <button
              type="button"
              onClick={() => router.back()}
              className="px-4 py-2 border rounded-md"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md disabled:opacity-50"
            >
              {submitting ? "Creating..." : "Create Draft"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

