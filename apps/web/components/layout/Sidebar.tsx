"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Package,
  Database,
  Scan,
  Bot,
  Users,
  FileText,
  Settings,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { cn } from "@/lib/utils";

const navigation = [
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  {
    name: "Adapters",
    icon: Package,
    children: [
      { name: "Registry", href: "/adapters" },
      { name: "Studio", href: "/studio/new" },
      { name: "Validation Runs", href: "/validation-runs" },
      { name: "Releases", href: "/adapters/releases" },
    ],
  },
  { name: "Connections", href: "/connections", icon: Database },
  {
    name: "Engine",
    icon: Scan,
    children: [
      { name: "Scans", href: "/engine/scans" },
      { name: "Results", href: "/engine/results" },
      { name: "Policies", href: "/engine/policies" },
    ],
  },
  {
    name: "AI",
    icon: Bot,
    children: [
      { name: "Providers", href: "/ai/providers" },
      { name: "Guardrails & Prompts", href: "/ai/guardrails" },
      { name: "Audit", href: "/ai/audit" },
    ],
  },
  {
    name: "Administration",
    icon: Users,
    children: [
      { name: "Users & Roles", href: "/admin/users" },
      { name: "Audit Logs", href: "/admin/audit" },
      { name: "System Logs", href: "/admin/system-logs" },
    ],
  },
  { name: "Settings", href: "/settings", icon: Settings },
];

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const [expandedItems, setExpandedItems] = useState<string[]>([]);
  const pathname = usePathname();

  const toggleExpanded = (name: string) => {
    setExpandedItems((prev) =>
      prev.includes(name) ? prev.filter((n) => n !== name) : [...prev, name]
    );
  };

  return (
    <div
      className={cn(
        "bg-card border-r transition-all duration-300 flex flex-col",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="p-4 border-b flex items-center justify-between">
        {!collapsed && <h1 className="text-xl font-bold">PK-FK Discovery</h1>}
        <button
          onClick={() => setCollapsed(!collapsed)}
          className="p-1 hover:bg-accent rounded"
        >
          {collapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
        </button>
      </div>

      <nav className="flex-1 overflow-y-auto p-4 space-y-1">
        {navigation.map((item) => {
          const isActive = pathname === item.href || (item.children && item.children.some(child => pathname === child.href));
          const isExpanded = expandedItems.includes(item.name);

          if (item.children) {
            return (
              <div key={item.name}>
                <button
                  onClick={() => toggleExpanded(item.name)}
                  className={cn(
                    "w-full flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                    isActive
                      ? "bg-primary text-primary-foreground"
                      : "hover:bg-accent hover:text-accent-foreground"
                  )}
                >
                  <item.icon size={20} />
                  {!collapsed && <span className="flex-1 text-left">{item.name}</span>}
                </button>
                {!collapsed && isExpanded && (
                  <div className="ml-8 mt-1 space-y-1">
                    {item.children.map((child) => (
                      <Link
                        key={child.href}
                        href={child.href}
                        className={cn(
                          "block px-3 py-2 rounded-md text-sm transition-colors",
                          pathname === child.href
                            ? "bg-primary text-primary-foreground"
                            : "hover:bg-accent hover:text-accent-foreground"
                        )}
                      >
                        {child.name}
                      </Link>
                    ))}
                  </div>
                )}
              </div>
            );
          }

          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                isActive
                  ? "bg-primary text-primary-foreground"
                  : "hover:bg-accent hover:text-accent-foreground"
              )}
            >
              <item.icon size={20} />
              {!collapsed && <span>{item.name}</span>}
            </Link>
          );
        })}
      </nav>
    </div>
  );
}

