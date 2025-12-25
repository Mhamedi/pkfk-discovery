"use client";

export default function AIGuardrailsPage() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">AI Guardrails & Prompts</h1>

      <div className="bg-card border rounded-lg p-6 max-w-3xl">
        <div className="space-y-6">
          <div>
            <h2 className="text-xl font-semibold mb-4">Allowed Inputs</h2>
            <div className="space-y-2">
              <label className="flex items-center gap-2">
                <input type="checkbox" checked readOnly className="rounded" />
                <span>Schema identifiers</span>
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked readOnly className="rounded" />
                <span>Probe outputs</span>
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked readOnly className="rounded" />
                <span>SQL templates</span>
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked readOnly className="rounded" />
                <span>Validation errors</span>
              </label>
              <label className="flex items-center gap-2">
                <input type="checkbox" checked readOnly className="rounded" />
                <span>EXPLAIN output</span>
              </label>
            </div>
          </div>

          <div>
            <h2 className="text-xl font-semibold mb-4">Output Format</h2>
            <p className="text-muted-foreground">
              AI outputs must be unified diff patches affecting only adapter bundle files.
            </p>
          </div>

          <div>
            <h2 className="text-xl font-semibold mb-4">Prompt Templates</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-1">System Prompt</label>
                <textarea
                  readOnly
                  className="w-full px-3 py-2 border rounded-md font-mono text-sm"
                  rows={4}
                  value="You are an expert SQL optimizer. Analyze the provided SQL template and suggest improvements. Output only unified diff format."
                />
              </div>
              <div>
                <label className="block text-sm font-medium mb-1">User Prompt Template</label>
                <textarea
                  readOnly
                  className="w-full px-3 py-2 border rounded-md font-mono text-sm"
                  rows={6}
                  value="Optimize the following SQL template for better performance:\n\n{SQL_TEMPLATE}\n\nValidation errors: {VALIDATION_ERRORS}\n\nEXPLAIN output: {EXPLAIN_OUTPUT}"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

